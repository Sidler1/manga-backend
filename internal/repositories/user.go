package repositories

import (
	"github.com/sidler1/manga-backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id uint) (*models.User, error)
	FindFavorites(userID uint) ([]models.Manga, error)
	Update(user *models.User) error
	FindUsersByFavoriteManga(mangaID uint) ([]models.User, error)
	// Add more as needed, e.g., Create for registration
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Favorites").Preload("Bookmarks").First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindFavorites(userID uint) ([]models.Manga, error) {
	var mangas []models.Manga
	err := r.db.Table("mangas").
		Joins("JOIN user_favorites ON user_favorites.manga_id = mangas.id").
		Where("user_favorites.user_id = ?", userID).
		Preload("Tags").Preload("Chapters").Preload("Website").
		Find(&mangas).Error
	return mangas, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}
func (r *userRepository) FindUsersByFavoriteManga(mangaID uint) ([]models.User, error) {
	var users []models.User
	err := r.db.Table("users").
		Joins("JOIN user_favorites ON user_favorites.user_id = users.id").
		Where("user_favorites.manga_id = ?", mangaID).
		Find(&users).Error
	return users, err
}
