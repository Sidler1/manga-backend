package handlers

import (
	"net/http"
	_ "strconv"

	"github.com/gin-gonic/gin"
	"github.com/sidler1/manga-backend/internal/repositories"
	"github.com/sidler1/manga-backend/internal/services"
)

func AddWebsite(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			URL  string `json:"url"`
			Name string `json:"name"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Assume admin check middleware
		err := s.AddWebsite(req.URL, req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "website added"})
	}
}

func GetWebsites(repo repositories.WebsiteRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		websites, err := repo.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, websites)
	}
}
