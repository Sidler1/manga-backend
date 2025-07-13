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

// GetLatestUpdates fetches the latest chapter updates from the homepage.
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
	doc.Find("div.lastest-update .update-item").Each(func(i int, selection *goquery.Selection) { // Adjust selector if site changes; based on standard manga site patterns
		titleLink := selection.Find(".manga-info h3 a")
		title := strings.TrimSpace(titleLink.Text())
		href, _ := titleLink.Attr("href")
		slug := strings.TrimSuffix(strings.TrimPrefix(href, s.baseURL+"manga/"), "/")
		chapter := strings.TrimSpace(selection.Find(".chapter a").Text()) // e.g., "Chapter 123"
		updateDate := strings.TrimSpace(selection.Find(".time").Text())   // e.g., "2 hours ago"

		if title != "" && slug != "" {
			updates = append(updates, Update{
				MangaTitle:    title,
				MangaSlug:     slug,
				ChapterNumber: chapter,
				UpdateDate:    updateDate,
			})
		}
	})

	return updates, nil
}

// GetMangaDetails fetches metadata for a specific manga using its slug.
func (s *MangaReadScraper) GetMangaDetails(slug string) (Manga, error) {
	url := s.baseURL + "manga/" + slug + "/"
	resp, err := http.Get(url)
	if err != nil {
		return Manga{}, ErrScrapeFailed
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Manga{}, ErrScrapeFailed
	}

	title := strings.TrimSpace(doc.Find(".manga-info h1").Text())
	description := strings.TrimSpace(doc.Find(".summary-content p").Text())
	author := strings.TrimSpace(doc.Find(".author-row span").Last().Text())
	status := strings.TrimSpace(doc.Find(".status-row span").Last().Text())
	coverURL, _ := doc.Find(".manga-poster img").Attr("src")

	var tags []string
	doc.Find(".genres a").Each(func(i int, selection *goquery.Selection) {
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
