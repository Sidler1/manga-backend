package repositories

import (
	"github.com/sidler1/manga-backend/internal/models"
	"gorm.io/gorm"
)

type WebsiteRepository interface {
	Create(website *models.Website) error
	FindByID(id uint) (*models.Website, error)
	FindAll() ([]models.Website, error)
	Update(website *models.Website) error
	Delete(id uint) error
}

type websiteRepository struct {
	db *gorm.DB
}

func NewWebsiteRepository(db *gorm.DB) WebsiteRepository {
	return &websiteRepository{db: db}
}

func (r *websiteRepository) Create(website *models.Website) error {
	return r.db.Create(website).Error
}

func (r *websiteRepository) FindByID(id uint) (*models.Website, error) {
	var website models.Website
	err := r.db.First(&website, id).Error
	return &website, err
}

func (r *websiteRepository) FindAll() ([]models.Website, error) {
	var websites []models.Website
	err := r.db.Find(&websites).Error
	return websites, err
}

func (r *websiteRepository) Update(website *models.Website) error {
	return r.db.Save(website).Error
}

func (r *websiteRepository) Delete(id uint) error {
	return r.db.Delete(&models.Website{}, id).Error
}
