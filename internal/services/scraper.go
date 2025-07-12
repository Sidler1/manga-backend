// internal/services/scraper.go
package services

import (
	"log"
	"time"

	"github.com/sidler1/manga-backend/internal/models"
	"github.com/sidler1/manga-backend/internal/repositories"
	// For scraping: import "github.com/gocolly/colly"
	// Implement site-specific logic; for now, placeholder.
)

type MangaUpdate struct {
	MangaID     uint
	NewChapter  string
	ChapterNum  uint
	Title       string
	ReleaseDate time.Time
	URL         string
	// Add tags if scraped
}

type ScraperService interface {
	CheckForUpdates() error
	ScrapeWebsite(website *models.Website) ([]MangaUpdate, error)
	AddWebsite(website *models.Website) error
}

type scraperService struct {
	websiteRepo repositories.WebsiteRepository
	mangaRepo   repositories.MangaRepository
	chapterRepo repositories.ChapterRepository
	tagRepo     repositories.TagRepository
}

func NewScraperService(websiteRepo repositories.WebsiteRepository, mangaRepo repositories.MangaRepository, chapterRepo repositories.ChapterRepository, tagRepo repositories.TagRepository) ScraperService {
	return &scraperService{
		websiteRepo: websiteRepo,
		mangaRepo:   mangaRepo,
		chapterRepo: chapterRepo,
		tagRepo:     tagRepo,
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
			// Assume LastChapter is string like "Chapter 100", compare accordingly
			if update.NewChapter > manga.LastChapter { // Implement proper comparison, e.g., parse numbers
				manga.LastChapter = update.NewChapter
				manga.UpdateTime = time.Now()

				// Add new chapter
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

				// Update estimated next
				chapters, _ := s.chapterRepo.FindByMangaID(manga.ID)
				manga.EstimatedNext = calculateEstimatedNext(chapters)

				if err := s.mangaRepo.Update(manga); err != nil {
					log.Printf("Error updating manga: %v", err)
				}

				// Tags would be handled during initial scrape or separately

				// Trigger notifications: this service would call notificationService in full impl, but since circular, inject or event-based later
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
	// Placeholder: Implement using colly or similar.
	// For each site, parse updates, match to existing mangas or create new if monitoring all.
	// Return list of updates.
	return []MangaUpdate{}, nil
}

func (s *scraperService) AddWebsite(website *models.Website) error {
	return s.websiteRepo.Create(website)
}

func calculateEstimatedNext(chapters []models.Chapter) time.Time {
	if len(chapters) < 2 {
		return time.Now().Add(7 * 24 * time.Hour) // Default 1 week
	}
	// Simple average interval
	var totalInterval time.Duration
	for i := 1; i < len(chapters); i++ {
		totalInterval += chapters[i].ReleaseDate.Sub(chapters[i-1].ReleaseDate)
	}
	avgInterval := totalInterval / time.Duration(len(chapters)-1)
	return chapters[len(chapters)-1].ReleaseDate.Add(avgInterval)
}
