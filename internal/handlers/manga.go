package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sidler1/manga-backend/internal/services"
)

func GetMangas(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mangas, err := s.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, mangas)
	}
}

func GetManga(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		manga, err := s.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "manga not found"})
			return
		}
		c.JSON(http.StatusOK, manga)
	}
}

func SearchMangas(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Tags []string `json:"tags"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mangas, err := s.SearchByTags(req.Tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, mangas)
	}
}

func FavoriteManga(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		mangaID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid manga id"})
			return
		}
		userID := getUserIDFromContext(c) // Assume middleware sets userID in context
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		err = s.FavoriteManga(userID, uint(mangaID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "favorited"})
	}
}
