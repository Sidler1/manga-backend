package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sidler1/manga-backend/internal/repositories"
)

func AdminMiddleware(userRepo repositories.UserRepository, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement JWT token validation using jwtSecret
		// Extract user ID from the token
		// Check if the user is an admin
		// If not an admin, return Forbidden status
		// Otherwise, call c.Next()
	}
}
