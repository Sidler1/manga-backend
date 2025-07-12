package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sidler1/manga-backend/internal/models"
	"github.com/sidler1/manga-backend/internal/repositories"
	"github.com/sidler1/manga-backend/scraper"
)

var scrapers = map[string]scraper.Scraper{}

func RegisterScrapers() {
	scrapers["https://www.mangaread.org/"] = scraper.NewMangaReadScraper()
}

func GetScraperForWebsite(url string) (scraper.Scraper, bool) {
	s, ok := scrapers[url]
	return s, ok
}

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

func NewScraperService(websiteRepo repositories.WebsiteRepository, mangaRepo repositories.MangaRepository, chapterRepo repositories.ChapterRepository, tagRepo repositories.TagRepository, notificationService NotificationService) ScraperService {
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

				_ = s.notificationService.SendUpdateNotification(manga)
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
	scraper, ok := GetScraperForWebsite(website.URL)
	if !ok {
		return nil, fmt.Errorf("no scraper found for website: %s", website.URL)
	}

	updates, err := scraper.GetLatestUpdates()
	if err != nil {
		return nil, err
	}

	var mangaUpdates []MangaUpdate
	for _, update := range updates {
		manga, err := s.mangaRepo.FindBySlug(update.MangaSlug)
		if err != nil {
			mangaDetails, err := scraper.GetMangaDetails(update.MangaSlug)
			if err != nil {
				log.Printf("Error fetching manga details for %s: %v", update.MangaSlug, err)
				continue
			}
			manga = &models.Manga{
				Title:       mangaDetails.Title,
				Description: mangaDetails.Description,
				Author:      mangaDetails.Author,
			}
			err = s.mangaRepo.Create(manga)
			if err != nil {
				log.Printf("Error creating manga %s: %v", mangaDetails.Title, err)
				continue
			}
			for _, tag := range mangaDetails.Tags {
				err = s.tagRepo.AddTagToManga(manga.ID, tag)
				if err != nil {
					log.Printf("Error adding tag %s to manga %s: %v", tag, manga.Title, err)
				}
			}
		}

		chapterNum, _ := strconv.ParseUint(strings.TrimPrefix(update.ChapterNumber, "Chapter "), 10, 32)
		mangaUpdates = append(mangaUpdates, MangaUpdate{
			MangaID:     manga.ID,
			NewChapter:  update.ChapterNumber,
			ChapterNum:  uint(chapterNum),
			Title:       update.MangaTitle,
			ReleaseDate: time.Now(),
			URL:         website.URL + "manga/" + update.MangaSlug + "/" + update.ChapterNumber,
		})
	}

	return mangaUpdates, nil
}

func (s *scraperService) AddWebsite(website *models.Website) error {
	return s.websiteRepo.Create(website)
}

func (s *scraperService) GetAllWebsites() ([]models.Website, error) {
	return s.websiteRepo.FindAll()
}

func calculateEstimatedNext(chapters []models.Chapter) time.Time {
	// Implement logic to estimate next chapter release
	// This is a placeholder implementation
	return time.Now().Add(24 * time.Hour)
}
