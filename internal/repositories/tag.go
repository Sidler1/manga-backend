package repositories

import (
	"github.com/sidler1/manga-backend/internal/models"

	"gorm.io/gorm"
)

type TagRepository interface {
	FindOrCreate(name string) (*models.Tag, error)
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
