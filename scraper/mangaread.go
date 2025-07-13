package scraper

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// MangaReadScraper implements the Scraper interface for https://www.mangaread.org/.
type MangaReadScraper struct {
	baseURL string
}

// NewMangaReadScraper initializes the scraper.
func NewMangaReadScraper() Scraper {
	return &MangaReadScraper{baseURL: "https://www.mangaread.org/"}
}

func (s *MangaReadScraper) GetBaseUrl() string {
	return s.baseURL
}

// GetLatestUpdates fetches the latest chapter updates from the MangaRead homepage.
//
// This function scrapes the homepage of MangaRead.org to extract information about
// the most recently updated manga chapters. It parses the HTML content to collect
// details such as manga title, slug, latest chapter number, and update time.
//
// Returns:
//   - []Update: A slice of Update structs, each containing information about a single
//     manga update, including the manga title, slug, latest chapter number, and update time.
//   - error: An error of type ErrScrapeFailed if the scraping process fails at any point,
//     or nil if the operation is successful.
func (s *MangaReadScraper) GetLatestUpdates() ([]Update, error) {
	resp, err := http.Get(s.baseURL)
	if err != nil {
		return nil, ErrScrapeFailed
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, ErrScrapeFailed
	}

	var updates []Update

	doc.Find(".page-content-listing div").Each(func(i int, selection *goquery.Selection) {
		selection.Find(".col-12").Each(func(j int, subSelection *goquery.Selection) {
			t := subSelection.Find(".post-title h3 a")
			titleLink, _ := t.Attr("href")
			title := t.Text()
			slug := strings.TrimSuffix(strings.TrimPrefix(titleLink, s.baseURL+"manga/"), "/")
			chapter := strings.TrimSpace(subSelection.Find(".chapter").Text())
			updateDate := strings.TrimSpace(subSelection.Find(".post-on").Text())
			if title != "" && slug != "" {
				updates = append(updates, Update{
					MangaTitle:    title,
					MangaSlug:     slug,
					ChapterNumber: chapter,
					UpdateDate:    updateDate,
				})
			}
		})
	})

	return updates, nil
}

// GetMangaDetails fetches and returns detailed information about a specific manga from MangaRead.org.
//
// This function scrapes the manga's dedicated page to extract various metadata such as title,
// description, author, status, cover image URL, and tags.
//
// Parameters:
//   - slug: A string representing the unique identifier of the manga in the URL.
//
// Returns:
//   - Manga: A Manga struct containing the scraped details of the manga.
//   - error: An error if the scraping process fails (ErrScrapeFailed), or nil if successful.
func (s *MangaReadScraper) GetMangaDetails(slug string) (Manga, error) {
	url := s.baseURL + "manga/" + slug + "/"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return Manga{}, ErrScrapeFailed
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil || doc.Text() == "" {
		return Manga{}, ErrScrapeFailed
	}

	title := strings.TrimSpace(doc.Find(".post-title h1").Text())
	description := strings.TrimSpace(doc.Find(".summary__content p").Text())
	author := strings.TrimSpace(doc.Find(".author-content a").Last().Text())
	status := strings.TrimSpace(doc.Find(".post-status .post-content_item .summary-content").Last().Text())
	coverURL, _ := doc.Find(".summary_image a img").Attr("src")

	var tags []string
	doc.Find(".genres-content a").Each(func(i int, selection *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(selection.Text()))
	})

	return Manga{
		Title:       title,
		Description: description,
		Author:      author,
		Status:      status,
		CoverURL:    coverURL,
		Tags:        tags,
	}, nil
}

// GetChapterList fetches the list of chapters for a specific manga from MangaRead.org.
//
// Parameters:
//   - slug: A string representing the unique identifier of the manga in the URL.
//
// Returns:
//   - []Chapter: A slice of Chapter structs, each containing information about a single chapter.
//     The chapters are sorted in descending order, with the most recent chapter first.
//   - error: An error if the scraping process fails, or nil if successful.
func (s *MangaReadScraper) GetChapterList(slug string) ([]Chapter, error) {
	url := s.baseURL + "manga/" + slug + "/"
	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrScrapeFailed
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, ErrScrapeFailed
	}

	var chapters []Chapter
	doc.Find(".chapters-list ul li").Each(func(i int, selection *goquery.Selection) {
		chapterLink := selection.Find("a")
		number := strings.TrimSpace(chapterLink.Text()) // e.g., "Chapter 1"
		href, _ := chapterLink.Attr("href")
		date := strings.TrimSpace(selection.Find(".chapter-release-date").Text())

		chapters = append(chapters, Chapter{
			Number: number,
			Title:  "", // Mangaread.org typically doesn't have chapter titles; can extend if needed
			Date:   date,
			URL:    href, // Full URL to original chapter
		})
	})

	// Reverse if needed to have latest first; assuming site lists oldest first
	for i, j := 0, len(chapters)-1; i < j; i, j = i+1, j-1 {
		chapters[i], chapters[j] = chapters[j], chapters[i]
	}

	return chapters, nil
}
