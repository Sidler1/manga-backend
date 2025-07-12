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
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindFavorites(userID uint) ([]models.Manga, error) {
	var favorites []models.Manga
	err := r.db.Model(&models.User{}).Where("id = ?", userID).Association("Favorites").Find(&favorites)
	return favorites, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) FindUsersByFavoriteManga(mangaID uint) ([]models.User, error) {
	var users []models.User
	err := r.db.Joins("JOIN user_favorites ON user_favorites.user_id = users.id").
		Where("user_favorites.manga_id = ?", mangaID).
		Find(&users).Error
	return users, err
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}
