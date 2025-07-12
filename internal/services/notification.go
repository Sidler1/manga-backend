package services

import (
	"fmt"
	"log"
	"time"

	"github.com/sidler1/manga-backend/internal/models"
	"github.com/sidler1/manga-backend/internal/repositories"
)

type NotificationService interface {
	SendUpdateNotification(manga *models.Manga) error
	GetNotifications(userID uint) ([]models.Notification, error)
}

type notificationService struct {
	userRepo         repositories.UserRepository
	notificationRepo repositories.NotificationRepository
	mangaRepo        repositories.MangaRepository
	// Add email sender or push service
}

func NewNotificationService(userRepo repositories.UserRepository, notificationRepo repositories.NotificationRepository, mangaRepo repositories.MangaRepository) NotificationService {
	return &notificationService{
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		mangaRepo:        mangaRepo,
	}
}

func (s *notificationService) SendUpdateNotification(manga *models.Manga) error {
	// Find users who favorited this manga
	users, err := s.findUsersFavoritedManga(manga.ID)
	if err != nil {
		return err
	}

	for _, user := range users {
		message := fmt.Sprintf("New chapter available for %s: %s", manga.Title, manga.LastChapter)
		notification := &models.Notification{
			UserID:  user.ID,
			MangaID: manga.ID,
			Message: message,
			SentAt:  time.Now(),
		}
		if err := s.notificationRepo.Create(notification); err != nil {
			log.Printf("Error creating notification: %v", err)
			continue
		}
		// Send via email or push; placeholder
	}
	return nil
}

func (s *notificationService) GetNotifications(userID uint) ([]models.Notification, error) {
	return s.notificationRepo.FindByUserID(userID)
}

func (s *notificationService) findUsersFavoritedManga(mangaID uint) ([]models.User, error) {
	return s.userRepo.FindUsersByFavoriteManga(mangaID)
}
