package main

import (
	"log"

	"github.com/sidler1/manga-backend/internal/config"
	"github.com/sidler1/manga-backend/internal/database"
	"github.com/sidler1/manga-backend/internal/handlers"
	"github.com/sidler1/manga-backend/internal/middlewares"
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
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authMiddleware := middlewares.AuthMiddleware(cfg.JWTSecret)
	notificationService := services.NewNotificationService(userRepo, notificationRepo, mangaRepo)
	scraperService := services.NewScraperService(websiteRepo, mangaRepo, chapterRepo, tagRepo, notificationService)
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
	r.Use(middlewares.ErrorHandler())
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.Logger())
	r.Use(middlewares.RateLimiter(10)) // 10 requests per second

	api := r.Group("/api/v1")
	{
		// Auth routes (no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register(authService))
			auth.POST("/login", handlers.Login(authService))
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(authMiddleware)
		{
			mangas := protected.Group("/mangas")
			{
				mangas.GET("/", handlers.GetMangas(mangaService))
				mangas.GET("/:id", handlers.GetManga(mangaService))
				mangas.POST("/search", handlers.SearchMangas(mangaService))
				mangas.POST("/:id/favorite", handlers.FavoriteManga(mangaService))
			}

			bookmarks := protected.Group("/bookmarks")
			{
				bookmarks.POST("/:manga_id", handlers.SetBookmark(mangaService))
				bookmarks.GET("/:manga_id", handlers.GetBookmark(mangaService))
			}

			users := protected.Group("/users")
			{
				users.GET("/me", handlers.GetCurrentUser(userRepo))
				users.GET("/favorites", handlers.GetUserFavorites(mangaService))
				users.GET("/notifications", handlers.GetNotifications(notificationService))
			}

			favorites := protected.Group("/favorites")
			{
				favorites.GET("/", handlers.GetUserFavorites(mangaService))
				favorites.POST("/:manga_id", handlers.FavoriteManga(mangaService))
				favorites.DELETE("/:manga_id", handlers.UnfavoriteManga(mangaService))
				favorites.GET("/updates", handlers.GetFavoriteUpdates(mangaService))
			}
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(authMiddleware, middlewares.AdminMiddleware(userRepo, cfg.JWTSecret))
		{
			admin.POST("/websites", handlers.AddWebsite(mangaService))
		}
	}

	services.RegisterScrapers() // Register scrapers for supported websites

	// Run server
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
