package repositories

import (
	"errors"
	"time"

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
	FindAllWithPagination(page, limit int, filters map[string]string, sort string) ([]models.Manga, int, error)
	FindChaptersByMangaID(mangaID uint) ([]models.Chapter, error)
	FindFavoritesWithUpdates(userID uint, since time.Time) ([]models.Manga, error)
	FindBySlug(slug string) (*models.Manga, error)
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

func (r *mangaRepository) FindUserFavorites(userID uint, page, limit int) ([]models.Manga, int, error) {
	var favorites []models.Manga
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.Manga{}).
		Joins("JOIN user_favorites ON user_favorites.manga_id = mangas.id").
		Where("user_favorites.user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Tags").Preload("Chapters").Preload("Website").
		Joins("JOIN user_favorites ON user_favorites.manga_id = mangas.id").
		Where("user_favorites.user_id = ?", userID).
		Offset(offset).Limit(limit).
		Find(&favorites).Error

	return favorites, int(total), err
}

func (r *mangaRepository) FindAllWithPagination(page, limit int, filters map[string]string, sort string) ([]models.Manga, int, error) {
	var mangas []models.Manga
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&models.Manga{}).Preload("Tags").Preload("Chapters").Preload("Website")

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" ILIKE ?", "%"+value+"%")
	}

	// Apply sorting
	if sort != "" {
		query = query.Order(sort)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Find(&mangas).Error
	return mangas, int(total), err
}

func (r *mangaRepository) FindChaptersByMangaID(mangaID uint) ([]models.Chapter, error) {
	var chapters []models.Chapter
	err := r.db.Where("manga_id = ?", mangaID).Order("number DESC").Find(&chapters).Error // DESC for latest first
	return chapters, err
}

func (r *mangaRepository) FindFavoritesWithUpdates(userID uint, since time.Time) ([]models.Manga, error) {
	var mangas []models.Manga
	err := r.db.Table("mangas").
		Joins("JOIN user_favorites ON user_favorites.manga_id = mangas.id").
		Where("user_favorites.user_id = ? AND mangas.update_time > ?", userID, since).
		Preload("Tags").Preload("Chapters").Preload("Website").
		Order("mangas.update_time DESC").
		Find(&mangas).Error
	return mangas, err
}

func (r *mangaRepository) FindBySlug(slug string) (*models.Manga, error) {
	var manga models.Manga
	err := r.db.Where("slug = ?", slug).First(&manga).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &manga, nil
}
