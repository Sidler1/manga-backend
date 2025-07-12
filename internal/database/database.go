package database

import (
	"github.com/sidler1/manga-backend/internal/models"

	"gorm.io/driver/postgres" // Change to gorm.io/driver/sqlite for local dev if needed
	"gorm.io/gorm"
)

func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Website{},
		&models.Manga{},
		&models.Tag{},
		&models.Chapter{},
		&models.User{},
		&models.Bookmark{},
		&models.Notification{},
	)
}
