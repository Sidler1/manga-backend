package services

import (
	"log"
	"time"

	"github.com/sidler1/manga-backend/internal/models"
	"github.com/sidler1/manga-backend/internal/repositories"
)

type MangaUpdate struct {
	MangaID     uint
	NewChapter  string
	ChapterNum  uint
	Title       string
	ReleaseDate time.Time
	URL         string
}

type ScraperService interface {
	CheckForUpdates() error
	ScrapeWebsite(website *models.Website) ([]MangaUpdate, error)
	AddWebsite(website *models.Website) error
	GetAllWebsites() ([]models.Website, error)
}

type scraperService struct {
	websiteRepo         repositories.WebsiteRepository
	mangaRepo           repositories.MangaRepository
	chapterRepo         repositories.ChapterRepository
	tagRepo             repositories.TagRepository
	notificationService NotificationService
}

func NewScraperService(websiteRepo repositories.WebsiteRepository, mangaRepo repositories.MangaRepository, chapterRepo repositories.ChapterRepository, tagRepo repositories.TagRepository, notificationService NotificationService) *scraperService {
	return &scraperService{
		websiteRepo:         websiteRepo,
		mangaRepo:           mangaRepo,
		chapterRepo:         chapterRepo,
		tagRepo:             tagRepo,
		notificationService: notificationService,
	}
}

func (s *scraperService) CheckForUpdates() error {
	websites, err := s.websiteRepo.FindAll()
	if err != nil {
		return err
	}

	for _, w := range websites {
		if time.Since(w.LastChecked) < time.Hour {
			continue
		}

		updates, err := s.ScrapeWebsite(&w)
		if err != nil {
			log.Printf("Error scraping %s: %v", w.URL, err)
			continue
		}

		for _, update := range updates {
			manga, err := s.mangaRepo.FindByID(update.MangaID)
			if err != nil {
				log.Printf("Manga not found: %d", update.MangaID)
				continue
			}
			if update.NewChapter > manga.LastChapter {
				manga.LastChapter = update.NewChapter
				manga.UpdateTime = time.Now()

				newChapter := &models.Chapter{
					MangaID:     manga.ID,
					Number:      update.ChapterNum,
					Title:       update.Title,
					ReleaseDate: update.ReleaseDate,
					URL:         update.URL,
				}
				if err := s.chapterRepo.Create(newChapter); err != nil {
					log.Printf("Error creating chapter: %v", err)
				}

				chapters, _ := s.chapterRepo.FindByMangaID(manga.ID)
				manga.EstimatedNext = calculateEstimatedNext(chapters)

				if err := s.mangaRepo.Update(manga); err != nil {
					log.Printf("Error updating manga: %v", err)
				}

				// Send notification
				_ = s.notificationService.SendUpdateNotification(manga) // Ignore err for now
			}
		}

		w.LastChecked = time.Now()
		err = s.websiteRepo.Update(&w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *scraperService) ScrapeWebsite(website *models.Website) ([]MangaUpdate, error) {
	// Implement site-specific scraping logic
	// This is a placeholder implementation
	return []MangaUpdate{}, nil
}

func (s *scraperService) AddWebsite(website *models.Website) error {
	return s.websiteRepo.Create(website)
}

func calculateEstimatedNext(chapters []models.Chapter) time.Time {
	// Implement logic to estimate next chapter release
	// This is a placeholder implementation
	return time.Now().Add(24 * time.Hour)
}
