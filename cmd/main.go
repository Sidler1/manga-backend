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
	_ "github.com/sidler1/manga-backend/docs" // This is important!
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gorm.io/gorm"
)

// @title           Manga Update API
// @version         0.1.0
// @description     This is a simple API for managing manga updates.
// @termsOfService  https://api.isekai.info/tos

// @contact.name   sidler2
// @contact.url    http://www.isekai.info/
// @contact.email  support@isekai.info

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      https://api.isekai.info
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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
	println("Setting up hourly cron job...")
	c := cron.New()
	_, err = c.AddFunc("@every 10m", func() {
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
			//	@Summary		Register a new user
			//	@Description	Register a new user with the provided information
			//	@Tags			auth
			//	@Accept			json
			//	@Produce		json
			//	@Param			user	body		handlers.RegisterRequest	true	"User Registration Info"
			//	@Success		200		{object}	handlers.AuthResponse
			//	@Failure		400		{object}	handlers.ErrorResponse
			//	@Router			/auth/register [post]
			auth.POST("/register", handlers.Register(authService))
			//	@Summary		Login user
			//	@Description	Authenticate a user and return a JWT token
			//	@Tags			auth
			//	@Accept			json
			//	@Produce		json
			//	@Param			credentials	body		handlers.LoginRequest	true	"Login Credentials"
			//	@Success		200			{object}	handlers.AuthResponse
			//	@Failure		401			{object}	handlers.ErrorResponse
			//	@Router			/auth/login [post]
			auth.POST("/login", handlers.Login(authService))
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(authMiddleware)
		{
			mangas := protected.Group("/mangas")
			{
				//	@Summary		Login user
				//	@Description	Authenticate a user and return a JWT token
				//	@Tags			auth
				//	@Accept			json
				//	@Produce		json
				//	@Param			credentials	body		handlers.LoginRequest	true	"Login Credentials"
				//	@Success		200			{object}	handlers.AuthResponse
				//	@Failure		401			{object}	handlers.ErrorResponse
				//	@Router			/auth/login [post]
				mangas.GET("/", handlers.GetMangas(mangaService))
				//	@Summary		Get a specific manga
				//	@Description	Retrieve details of a specific manga by its ID
				//	@Tags			mangas
				//	@Accept			json
				//	@Produce		json
				//	@Param			id	path		int	true	"Manga ID"
				//	@Success		200	{object}	models.Manga
				//	@Failure		404	{object}	handlers.ErrorResponse
				//	@Failure		500	{object}	handlers.ErrorResponse
				//	@Router			/mangas/{id} [get]
				mangas.GET("/:id", handlers.GetManga(mangaService))
				//	@Summary		Search for mangas
				//	@Description	Search for mangas based on given criteria
				//	@Tags			mangas
				//	@Accept			json
				//	@Produce		json
				//	@Param			query	body		handlers.SearchRequest	true	"Search criteria"
				//	@Success		200		{array}		models.Manga
				//	@Failure		400		{object}	handlers.ErrorResponse
				//	@Failure		500		{object}	handlers.ErrorResponse
				//	@Router			/mangas/search [post]
				mangas.POST("/search", handlers.SearchMangas(mangaService))
				//	@Summary		Favorite a manga
				//	@Description	Add a manga to the user's favorites list
				//	@Tags			mangas
				//	@Accept			json
				//	@Produce		json
				//	@Param			id	path		int	true	"Manga ID"
				//	@Success		200	{object}	handlers.SuccessResponse
				//	@Failure		400	{object}	handlers.ErrorResponse
				//	@Failure		401	{object}	handlers.ErrorResponse
				//	@Failure		404	{object}	handlers.ErrorResponse
				//	@Failure		500	{object}	handlers.ErrorResponse
				//	@Security		ApiKeyAuth
				//	@Router			/mangas/{id}/favorite [post]
				mangas.POST("/:id/favorite", handlers.FavoriteManga(mangaService))
			}

			bookmarks := protected.Group("/bookmarks")
			{
				//	@Summary		Set a bookmark for a manga
				//	@Description	Set or update a bookmark for a specific manga
				//	@Tags			bookmarks
				//	@Accept			json
				//	@Produce		json
				//	@Param			manga_id	path		int							true	"Manga ID"
				//	@Param			bookmark	body		handlers.BookmarkRequest	true	"Bookmark details"
				//	@Success		200			{object}	handlers.SuccessResponse
				//	@Failure		400			{object}	handlers.ErrorResponse
				//	@Failure		401			{object}	handlers.ErrorResponse
				//	@Failure		404			{object}	handlers.ErrorResponse
				//	@Failure		500			{object}	handlers.ErrorResponse
				//	@Security		ApiKeyAuth
				//	@Router			/bookmarks/{manga_id} [post]
				bookmarks.POST("/:manga_id", handlers.SetBookmark(mangaService))
				//	@Summary		Get a bookmark for a manga
				//	@Description	Retrieve the bookmark for a specific manga
				//	@Tags			bookmarks
				//	@Accept			json
				//	@Produce		json
				//	@Param			manga_id	path		int	true	"Manga ID"
				//	@Success		200			{object}	models.Bookmark
				//	@Failure		401			{object}	handlers.ErrorResponse
				//	@Failure		404			{object}	handlers.ErrorResponse
				//	@Failure		500			{object}	handlers.ErrorResponse
				//	@Security		ApiKeyAuth
				//	@Router			/bookmarks/{manga_id} [get]
				bookmarks.GET("/:manga_id", handlers.GetBookmark(mangaService))
			}

			users := protected.Group("/users")
			{
				//	@Summary		Get current user
				//	@Description	Retrieve details of the currently authenticated user
				//	@Tags			users
				//	@Accept			json
				//	@Produce		json
				//	@Success		200	{object}	models.User
				//	@Failure		401	{object}	handlers.ErrorResponse
				//	@Failure		500	{object}	handlers.ErrorResponse
				//	@Security		ApiKeyAuth
				//	@Router			/users/me [get]
				users.GET("/me", handlers.GetCurrentUser(userRepo))
				//	@Summary		Get user's favorite mangas
				//	@Description	Retrieve a list of mangas favorited by the current user
				//	@Tags			users
				//	@Accept			json
				//	@Produce		json
				//	@Success		200	{array}		models.Manga
				//	@Failure		401	{object}	handlers.ErrorResponse
				//	@Failure		500	{object}	handlers.ErrorResponse
				//	@Security		ApiKeyAuth
				//	@Router			/users/favorites [get]
				users.GET("/favorites", handlers.GetUserFavorites(mangaService))
				//	@Summary		Get user's notifications
				//	@Description	Retrieve a list of notifications for the current user
				//	@Tags			users
				//	@Accept			json
				//	@Produce		json
				//	@Success		200	{array}		models.Notification
				//	@Failure		401	{object}	handlers.ErrorResponse
				//	@Failure		500	{object}	handlers.ErrorResponse
				//	@Security		ApiKeyAuth
				//	@Router			/users/notifications [get]
				users.GET("/notifications", handlers.GetNotifications(notificationService))
			}

			favorites := protected.Group("/favorites")
			{
				//	@Summary		Get user's favorite mangas
				//	@Description	Retrieve a list of mangas favorited by the current user
				//	@Tags			favorites
				//	@Accept			json
				//	@Produce		json
				//	@Success		200	{array}		models.Manga
				//	@Failure		401	{object}	handlers.ErrorResponse
				//	@Failure		500	{object}	handlers.ErrorResponse
				//	@Security		ApiKeyAuth
				//	@Router			/favorites [get]
				favorites.GET("/", handlers.GetUserFavorites(mangaService))
				//	@Summary		Add manga to favorites
				//	@Description	Add a specific manga to the user's favorites list
				//	@Tags			favorites
				//	@Accept			json
				//	@Produce		json
				//	@Param			manga_id	path		int	true	"Manga ID"
				//	@Success		200			{object}	handlers.SuccessResponse
				//	@Failure		400			{object}	handlers.ErrorResponse
				//	@Failure		401			{object}	handlers.ErrorResponse
				//	@Failure		404			{object}	handlers.ErrorResponse
				//	@Failure		500			{object}	handlers.ErrorResponse
				//	@Security		ApiKeyAuth
				//	@Router			/favorites/{manga_id} [post]
				favorites.POST("/:manga_id", handlers.FavoriteManga(mangaService))
				//	@Summary		Remove manga from favorites
				//	@Description	Remove a specific manga from the user's favorites list
				//	@Tags			favorites
				//	@Accept			json
				//	@Produce		json
				//	@Param			manga_id	path		int	true	"Manga ID"
				//	@Success		200			{object}	handlers.SuccessResponse
				//	@Failure		400			{object}	handlers.ErrorResponse
				//	@Failure		401			{object}	handlers.ErrorResponse
				//	@Failure		404			{object}	handlers.ErrorResponse
				//	@Failure		500			{object}	handlers.ErrorResponse
				//	@Security		ApiKeyAuth
				//	@Router			/favorites/{manga_id} [delete]
				favorites.DELETE("/:manga_id", handlers.UnfavoriteManga(mangaService))
				//	@Summary		Get updates for favorite mangas
				//	@Description	Retrieve a list of updates for the user's favorite mangas
				//	@Tags			favorites
				//	@Accept			json
				//	@Produce		json
				//	@Success		200	{array}		models.MangaUpdate
				//	@Failure		401	{object}	handlers.ErrorResponse
				//	@Failure		500	{object}	handlers.ErrorResponse
				//	@Security		ApiKeyAuth
				//	@Router			/favorites/updates [get]
				favorites.GET("/updates", handlers.GetFavoriteUpdates(mangaService))
			}
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(authMiddleware, middlewares.AdminMiddleware(userRepo, cfg.JWTSecret))
		{
			//	@Summary		Add a new website
			//	@Description	Add a new website to the system for manga scraping
			//	@Tags			admin
			//	@Accept			json
			//	@Produce		json
			//	@Param			website	body		handlers.WebsiteRequest	true	"Website details"
			//	@Success		201		{object}	handlers.SuccessResponse
			//	@Failure		400		{object}	handlers.ErrorResponse
			//	@Failure		401		{object}	handlers.ErrorResponse
			//	@Failure		403		{object}	handlers.ErrorResponse
			//	@Failure		500		{object}	handlers.ErrorResponse
			//	@Security		ApiKeyAuth
			//	@Router			/admin/websites [post]
			admin.POST("/websites", handlers.AddWebsite(mangaService))
		}
	}

	services.RegisterScrapers() // Register scrapers for supported websites

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
