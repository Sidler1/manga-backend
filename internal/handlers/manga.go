package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sidler1/manga-backend/internal/services"
)

// GetMangas handles the request to get a list of mangas with pagination and filters
func GetMangas(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		if limit < 1 || limit > 100 {
			limit = 20
		}

		filters := c.QueryMap("filters")
		sort := c.DefaultQuery("sort", "")

		mangas, total, err := s.GetAll(page, limit, filters, sort)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"mangas": mangas,
			"total":  total,
			"page":   page,
			"limit":  limit,
		})
	}
}

// GetMangaChapters handles the request to get chapters for a specific manga
func GetMangaChapters(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid manga id"})
			return
		}
		chapters, err := s.GetMangaChapters(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, chapters)
	}
}

// GetManga handles the request to get a single manga by ID
func GetManga(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid manga id"})
			return
		}

		manga, err := s.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "manga not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"manga": manga})
	}
}

// FavoriteManga handles the request to favorite a manga
func FavoriteManga(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		mangaID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid manga id"})
			return
		}

		userID := c.GetUint("userID") // Assuming you set this in your auth middleware
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		err = s.FavoriteManga(userID, uint(mangaID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "manga favorited successfully"})
	}
}

func UnfavoriteManga(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		mangaIDStr := c.Param("manga_id")
		mangaID, err := strconv.ParseUint(mangaIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid manga id"})
			return
		}
		err = s.UnfavoriteManga(userID.(uint), uint(mangaID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "manga unfavorited"})
	}
}

func GetFavoriteUpdates(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		sinceStr := c.DefaultQuery("since", "")
		var since time.Time
		var err error
		if sinceStr != "" {
			since, err = time.Parse(time.RFC3339, sinceStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid since parameter"})
				return
			}
		} else {
			since = time.Now().Add(-24 * time.Hour) // Default to last 24 hours
		}
		updates, err := s.GetFavoriteUpdates(userID.(uint), since)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updates)
	}
}
func SearchMangas(s services.MangaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchQuery struct {
			Query string `json:"query" binding:"required"`
		}
		if err := c.ShouldBindJSON(&searchQuery); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search query"})
			return
		}

		results, err := s.SearchMangas(searchQuery.Query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"results": results})
	}
}
