package scraper

import "errors"

// Update represents a recent manga chapter update from the site.
type Update struct {
	MangaTitle    string
	MangaSlug     string // Slug for the manga (e.g., "solo-leveling-manhwa")
	ChapterNumber string
	UpdateDate    string
}

// Manga represents detailed metadata for a manga.
type Manga struct {
	Title       string
	Description string
	Author      string
	Status      string // e.g., "Ongoing", "Completed"
	CoverURL    string
	Tags        []string
}

// Chapter represents a single chapter in a manga's list.
type Chapter struct {
	Number string
	Title  string // Optional, as some sites may not provide chapter titles
	Date   string
	URL    string // Link to the original chapter on the site
}

// Scraper defines the interface for site-specific manga scrapers.
type Scraper interface {
	GetLatestUpdates() ([]Update, error)           // Fetch recent site-wide updates for hourly update checks
	GetMangaDetails(slug string) (Manga, error)    // Fetch details for a specific manga
	GetChapterList(slug string) ([]Chapter, error) // Fetch full chapter list for bookmarking and update detection
	GetBaseUrl() string                            // Get the base URL for the scraped site
}

// ErrScrapeFailed is a generic error for scraping issues.
var ErrScrapeFailed = errors.New("scrape failed due to site access or parsing error")
