package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sidler1/manga-backend/internal/services"
)

func GetUserFavorites(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserIDFromContext(c)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		favorites, err := s.GetUserFavorites(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, favorites)
	}
}

func GetNotifications(s services.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserIDFromContext(c)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		notifications, err := s.GetNotifications(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, notifications)
	}
}

// Placeholder for auth middleware extraction
func getUserIDFromContext(c *gin.Context) uint {
	// Implement with JWT or similar; return user ID from c.Get("userID")
	return 1 // Dummy
}
