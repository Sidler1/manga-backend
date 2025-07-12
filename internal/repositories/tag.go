package repositories

import (
	"github.com/sidler1/manga-backend/internal/models"

	"gorm.io/gorm"
)

type TagRepository interface {
	FindOrCreate(name string) (*models.Tag, error)
	AddTagToManga(mangaID uint, tag string) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindOrCreate(name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("name = ?", name).FirstOrCreate(&tag).Error
	return &tag, err
}

func (r *tagRepository) AddTagToManga(mangaID uint, tag string) error {
	// Implementation to add a tag to a manga
	// This is just an example and may need to be adjusted based on your actual database schema
	return r.db.Exec("INSERT INTO manga_tags (manga_id, tag) VALUES (?, ?)", mangaID, tag).Error
}
