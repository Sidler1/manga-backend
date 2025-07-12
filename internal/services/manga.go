package services

import (
	"errors"
	"time"

	"github.com/sidler1/manga-backend/internal/models"
	"github.com/sidler1/manga-backend/internal/repositories"
)

type MangaService interface {
	GetAll(page, limit int, search map[string]string, tags string) ([]models.Manga, int, error)
	GetByID(id uint) (*models.Manga, error)
	SearchByTags(tags string) ([]models.Manga, error)
	FavoriteManga(userID uint, mangaID uint) error
	GetUserFavorites(userID uint) ([]models.Manga, error)
	SetBookmark(userID uint, mangaID uint, chapter uint) error
	GetBookmark(userID uint, mangaID uint) (uint, error)
	AddWebsite(url string, name string) error
	GetMangaChapters(mangaID uint) ([]models.Chapter, error)
	UnfavoriteManga(userID uint, mangaID uint) error
	GetFavoriteUpdates(userID uint, since time.Time) ([]models.Manga, error)
	GetWebsites() ([]models.Website, error)
	SearchMangas(query string) ([]models.Manga, error)
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

func NewMangaService(mangaRepo repositories.MangaRepository, userRepo repositories.UserRepository, bookmarkRepo repositories.BookmarkRepository, chapterRepo repositories.ChapterRepository, tagRepo repositories.TagRepository, scraperService ScraperService, notificationService NotificationService) *mangaService {
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

func (s *mangaService) GetAll(page, limit int, search string, tags []string) ([]models.Manga, int, error) {
	return s.mangaRepo.FindAllWithPagination(page, limit, search, tags)
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
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return user.Favorites, nil
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

func (s *mangaService) GetMangaChapters(mangaID uint) ([]models.Chapter, error) {
	return s.mangaRepo.FindChaptersByMangaID(mangaID)
}

func (s *mangaService) UnfavoriteManga(userID uint, mangaID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	for i, fav := range user.Favorites {
		if fav.ID == mangaID {
			user.Favorites = append(user.Favorites[:i], user.Favorites[i+1:]...)
			break
		}
	}
	return s.userRepo.Update(user)
}

func (s *mangaService) GetFavoriteUpdates(userID uint, since time.Time) ([]models.Manga, error) {
	return s.mangaRepo.FindFavoritesWithUpdates(userID, since)
}

func (s *mangaService) GetWebsites() ([]models.Website, error) {
	return s.scraperService.GetAllWebsites()
}

func (s *mangaService) SearchMangas(query string) ([]models.Manga, error) {
	return []models.Manga{}, nil
}
