package repositories

import (
	"errors"
	"github.com/sidler1/manga-backend/internal/models"

	"gorm.io/gorm"
)

type BookmarkRepository interface {
	FindByUserAndManga(userID, mangaID uint) (*models.Bookmark, error)
	Upsert(bookmark *models.Bookmark) error
}

type bookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) BookmarkRepository {
	return &bookmarkRepository{db: db}
}

func (r *bookmarkRepository) FindByUserAndManga(userID, mangaID uint) (*models.Bookmark, error) {
	var bookmark models.Bookmark
	err := r.db.Where("user_id = ? AND manga_id = ?", userID, mangaID).First(&bookmark).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &bookmark, err
}

func (r *bookmarkRepository) Upsert(bookmark *models.Bookmark) error {
	return r.db.Save(bookmark).Error // GORM Save acts as upsert if ID exists
}
