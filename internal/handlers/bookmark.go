package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sidler1/manga-backend/internal/services"
)

func SetBookmark(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mangaIDStr := c.Param("manga_id")
		mangaID, err := strconv.ParseUint(mangaIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid manga id"})
			return
		}
		var req struct {
			Chapter uint `json:"chapter"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID := getUserIDFromContext(c)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		err = s.SetBookmark(userID, uint(mangaID), req.Chapter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "bookmark set"})
	}
}

func GetBookmark(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mangaIDStr := c.Param("manga_id")
		mangaID, err := strconv.ParseUint(mangaIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid manga id"})
			return
		}
		userID := getUserIDFromContext(c)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		chapter, err := s.GetBookmark(userID, uint(mangaID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"chapter": chapter})
	}
}
