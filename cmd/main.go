package main

import (
	"log"

	"github.com/sidler1/manga-backend/internal/config"
	"github.com/sidler1/manga-backend/internal/database"
	"github.com/sidler1/manga-backend/internal/handlers"
	"github.com/sidler1/manga-backend/internal/repositories"
	"github.com/sidler1/manga-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	_ "gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	// Initialize repositories
	mangaRepo := repositories.NewMangaRepository(db)
	userRepo := repositories.NewUserRepository(db)
	bookmarkRepo := repositories.NewBookmarkRepository(db)
	websiteRepo := repositories.NewWebsiteRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	chapterRepo := repositories.NewChapterRepository(db)

	// Initialize services
	scraperService := services.NewScraperService(websiteRepo, mangaRepo, chapterRepo, tagRepo)
	notificationService := services.NewNotificationService(userRepo, notificationRepo, mangaRepo)
	mangaService := services.NewMangaService(mangaRepo, userRepo, bookmarkRepo, chapterRepo, tagRepo, scraperService, notificationService)

	// Set up cron job for hourly updates
	c := cron.New()
	_, err = c.AddFunc("@hourly", func() {
		log.Println("Running hourly manga update check...")
		if err := scraperService.CheckForUpdates(); err != nil {
			log.Printf("Error during update check: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule cron job: %v", err)
	}
	c.Start()
	defer c.Stop()

	// Set up Gin router
	r := gin.Default()

	// API routes
	api := r.Group("/api/v1")
	{
		mangas := api.Group("/mangas")
		{
			mangas.GET("/", handlers.GetMangas(mangaService))
			mangas.GET("/:id", handlers.GetManga(mangaService))
			mangas.POST("/search", handlers.SearchMangas(mangaService))
			mangas.POST("/:id/favorite", handlers.FavoriteManga(mangaService)) // Requires auth middleware for userID
		}

		bookmarks := api.Group("/bookmarks")
		{
			bookmarks.POST("/:manga_id", handlers.SetBookmark(mangaService)) // Requires auth
			bookmarks.GET("/:manga_id", handlers.GetBookmark(mangaService))  // Requires auth
		}

		websites := api.Group("/websites")
		{
			websites.POST("/", handlers.AddWebsite(mangaService)) // Admin endpoint, requires auth/authorization
		}

		users := api.Group("/users")
		{
			users.GET("/favorites", handlers.GetUserFavorites(mangaService))            // Requires auth
			users.GET("/notifications", handlers.GetNotifications(notificationService)) // Requires auth
		}
	}

	// Run server
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
