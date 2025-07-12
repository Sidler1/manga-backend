package services

import (
	"errors"

	"github.com/sidler1/manga-backend/internal/models"
	"github.com/sidler1/manga-backend/internal/repositories"
)

type MangaService interface {
	GetAll() ([]models.Manga, error)
	GetByID(id uint) (*models.Manga, error)
	SearchByTags(tags []string) ([]models.Manga, error)
	FavoriteManga(userID uint, mangaID uint) error
	GetUserFavorites(userID uint) ([]models.Manga, error)
	SetBookmark(userID uint, mangaID uint, chapter uint) error
	GetBookmark(userID uint, mangaID uint) (uint, error)
	AddWebsite(url string, name string) error
	// More as needed
}

type mangaService struct {
	mangaRepo           repositories.MangaRepository
	userRepo            repositories.UserRepository
	bookmarkRepo        repositories.BookmarkRepository
	chapterRepo         repositories.ChapterRepository
	tagRepo             repositories.TagRepository
	scraperService      ScraperService
	notificationService NotificationService
}

func NewMangaService(mangaRepo repositories.MangaRepository, userRepo repositories.UserRepository, bookmarkRepo repositories.BookmarkRepository, chapterRepo repositories.ChapterRepository, tagRepo repositories.TagRepository, scraperService ScraperService, notificationService NotificationService) MangaService {
	return &mangaService{
		mangaRepo:           mangaRepo,
		userRepo:            userRepo,
		bookmarkRepo:        bookmarkRepo,
		chapterRepo:         chapterRepo,
		tagRepo:             tagRepo,
		scraperService:      scraperService,
		notificationService: notificationService,
	}
}

func (s *mangaService) GetAll() ([]models.Manga, error) {
	return s.mangaRepo.FindAll()
}

func (s *mangaService) GetByID(id uint) (*models.Manga, error) {
	return s.mangaRepo.FindByID(id)
}

func (s *mangaService) SearchByTags(tags []string) ([]models.Manga, error) {
	return s.mangaRepo.SearchByTags(tags)
}

func (s *mangaService) FavoriteManga(userID uint, mangaID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	manga, err := s.mangaRepo.FindByID(mangaID)
	if err != nil {
		return err
	}
	// Check if already favorited
	for _, fav := range user.Favorites {
		if fav.ID == mangaID {
			return errors.New("already favorited")
		}
	}
	user.Favorites = append(user.Favorites, *manga)
	return s.userRepo.Update(user)
}

func (s *mangaService) GetUserFavorites(userID uint) ([]models.Manga, error) {
	return s.userRepo.FindFavorites(userID)
}

func (s *mangaService) SetBookmark(userID uint, mangaID uint, chapter uint) error {
	bookmark := &models.Bookmark{
		UserID:  userID,
		MangaID: mangaID,
		Chapter: chapter,
	}
	return s.bookmarkRepo.Upsert(bookmark)
}

func (s *mangaService) GetBookmark(userID uint, mangaID uint) (uint, error) {
	bookmark, err := s.bookmarkRepo.FindByUserAndManga(userID, mangaID)
	if err != nil || bookmark == nil {
		return 0, err
	}
	return bookmark.Chapter, nil
}

func (s *mangaService) AddWebsite(url string, name string) error {
	website := &models.Website{
		URL:  url,
		Name: name,
	}
	return s.scraperService.AddWebsite(website)
}
