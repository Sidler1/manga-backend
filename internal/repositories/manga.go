package repositories

import (
	"github.com/sidler1/manga-backend/internal/models"

	"gorm.io/gorm"
)

type MangaRepository interface {
	FindAll() ([]models.Manga, error)
	FindByID(id uint) (*models.Manga, error)
	SearchByTags(tags []string) ([]models.Manga, error)
	Update(manga *models.Manga) error
	Create(manga *models.Manga) error
	FindByWebsiteID(websiteID uint) ([]models.Manga, error)
}

type mangaRepository struct {
	db *gorm.DB
}

func NewMangaRepository(db *gorm.DB) MangaRepository {
	return &mangaRepository{db: db}
}

func (r *mangaRepository) FindAll() ([]models.Manga, error) {
	var mangas []models.Manga
	err := r.db.Preload("Tags").Preload("Chapters").Preload("Website").Find(&mangas).Error
	return mangas, err
}

func (r *mangaRepository) FindByID(id uint) (*models.Manga, error) {
	var manga models.Manga
	err := r.db.Preload("Tags").Preload("Chapters").Preload("Website").First(&manga, id).Error
	return &manga, err
}

func (r *mangaRepository) SearchByTags(tags []string) ([]models.Manga, error) {
	var mangas []models.Manga
	query := r.db.Joins("JOIN manga_tags ON manga_tags.manga_id = mangas.id").
		Joins("JOIN tags ON tags.id = manga_tags.tag_id").
		Where("tags.name IN ?", tags).
		Group("mangas.id").
		Having("COUNT(DISTINCT tags.name) = ?", len(tags)).
		Preload("Tags").Preload("Chapters").Preload("Website")
	err := query.Find(&mangas).Error
	return mangas, err
}

func (r *mangaRepository) Update(manga *models.Manga) error {
	return r.db.Save(manga).Error
}

func (r *mangaRepository) Create(manga *models.Manga) error {
	return r.db.Create(manga).Error
}

func (r *mangaRepository) FindByWebsiteID(websiteID uint) ([]models.Manga, error) {
	var mangas []models.Manga
	err := r.db.Where("website_id = ?", websiteID).Preload("Tags").Preload("Chapters").Find(&mangas).Error
	return mangas, err
}
