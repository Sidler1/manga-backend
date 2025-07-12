package repositories

import (
	"github.com/sidler1/manga-backend/internal/models"

	"gorm.io/gorm"
)

type ChapterRepository interface {
	Create(chapter *models.Chapter) error
	FindByMangaID(mangaID uint) ([]models.Chapter, error)
}

type chapterRepository struct {
	db *gorm.DB
}

func NewChapterRepository(db *gorm.DB) ChapterRepository {
	return &chapterRepository{db: db}
}

func (r *chapterRepository) Create(chapter *models.Chapter) error {
	return r.db.Create(chapter).Error
}

func (r *chapterRepository) FindByMangaID(mangaID uint) ([]models.Chapter, error) {
	var chapters []models.Chapter
	err := r.db.Where("manga_id = ?", mangaID).Order("number ASC").Find(&chapters).Error
	return chapters, err
}
