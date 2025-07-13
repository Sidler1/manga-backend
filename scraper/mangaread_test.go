package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMangaReadScraper_GetMangaDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<div class="post-title">
						<h1>Test Manga</h1>
					</div>
					<div class="summary__content">
						<p>This is a test manga description.</p>
					</div>
					<div class="author-content">
						<span>Author:</span>
						<a>Test Author</a>
					</div>
					<div class="post-status">
						<div class="post-content_item">
							<span>Status:</span>
							<span class="summary-content">Ongoing</span>
						</div>
					</div>
					<div class="summary_image">
						<a><img src="https://example.com/cover.jpg" /></a>
					</div>
					<div class="genres-content">
						<a>Action</a>
						<a>Adventure</a>
					</div>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	scraper := &MangaReadScraper{baseURL: server.URL + "/"}
	manga, err := scraper.GetMangaDetails("test-manga")

	assert.NoError(t, err)
	assert.Equal(t, "Test Manga", manga.Title)
	assert.Equal(t, "This is a test manga description.", manga.Description)
	assert.Equal(t, "Test Author", manga.Author)
	assert.Equal(t, "Ongoing", manga.Status)
	assert.Equal(t, "https://example.com/cover.jpg", manga.CoverURL)
	assert.Equal(t, []string{"Action", "Adventure"}, manga.Tags)
}

func TestMangaReadScraper_GetMangaDetails_HTTPRequestFails(t *testing.T) {
	// Create a server that always fails
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	scraper := &MangaReadScraper{baseURL: server.URL + "/"}
	_, err := scraper.GetMangaDetails("test-manga")

	assert.Error(t, err)
	assert.Equal(t, ErrScrapeFailed, err)
}

func TestMangaReadScraper_GetMangaDetails_HTMLParsingFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><div class="post-title"><h1`))
	}))
	defer server.Close()

	scraper := &MangaReadScraper{baseURL: server.URL + "/"}
	_, err := scraper.GetMangaDetails("test-manga")

	assert.Error(t, err)
	assert.Equal(t, ErrScrapeFailed, err)
}

func TestMangaReadScraper_GetMangaDetails_NoDescription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<div class="post-title">
						<h1>Test Manga</h1>
					</div>
					<div class="summary__content">
						<p></p>
					</div>
					<div class="author-content">
						<span>Author:</span>
						<a>Test Author</a>
					</div>
					<div class="post-status">
						<div class="post-content_item">
							<span>Status:</span>
							<span class="summary-content">Ongoing</span>
						</div>
					</div>
					<div class="summary_image">
						<a><img src="https://example.com/cover.jpg" /></a>
					</div>
					<div class="genres-content">
						<a>Action</a>
						<a>Adventure</a>
					</div>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	scraper := &MangaReadScraper{baseURL: server.URL + "/"}
	manga, err := scraper.GetMangaDetails("test-manga")

	assert.NoError(t, err)
	assert.Equal(t, "Test Manga", manga.Title)
	assert.Equal(t, "", manga.Description)
	assert.Equal(t, "Test Author", manga.Author)
	assert.Equal(t, "Ongoing", manga.Status)
	assert.Equal(t, "https://example.com/cover.jpg", manga.CoverURL)
	assert.Equal(t, []string{"Action", "Adventure"}, manga.Tags)
}

func TestMangaReadScraper_GetMangaDetails_NoAuthor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<div class="post-title">
						<h1>Test Manga</h1>
					</div>
					<div class="summary__content">
						<p>This is a test manga description.</p>
					</div>
					<div class="author-content">
						<span>Author:</span>
						<a></a>
					</div>
					<div class="post-status">
						<div class="post-content_item">
							<span>Status:</span>
							<span class="summary-content">Ongoing</span>
						</div>
					</div>
					<div class="summary_image">
						<a><img src="https://example.com/cover.jpg" /></a>
					</div>
					<div class="genres-content">
						<a>Action</a>
						<a>Adventure</a>
					</div>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	scraper := &MangaReadScraper{baseURL: server.URL + "/"}
	manga, err := scraper.GetMangaDetails("test-manga")

	assert.NoError(t, err)
	assert.Equal(t, "Test Manga", manga.Title)
	assert.Equal(t, "This is a test manga description.", manga.Description)
	assert.Equal(t, "", manga.Author)
	assert.Equal(t, "Ongoing", manga.Status)
	assert.Equal(t, "https://example.com/cover.jpg", manga.CoverURL)
	assert.Equal(t, []string{"Action", "Adventure"}, manga.Tags)
}

func TestMangaReadScraper_GetMangaDetails_NoStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<div class="post-title">
						<h1>Test Manga</h1>
					</div>
					<div class="summary__content">
						<p>This is a test manga description.</p>
					</div>
					<div class="author-content">
						<span>Author:</span>
						<a>Test Author</a>
					</div>
					<div class="post-status">
						<div class="post-content_item">
							<span>Status:</span>
							<span class="summary-content"></span>
						</div>
					</div>
					<div class="summary_image">
						<a><img src="https://example.com/cover.jpg" /></a>
					</div>
					<div class="genres-content">
						<a>Action</a>
						<a>Adventure</a>
					</div>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	scraper := &MangaReadScraper{baseURL: server.URL + "/"}
	manga, err := scraper.GetMangaDetails("test-manga")

	assert.NoError(t, err)
	assert.Equal(t, "Test Manga", manga.Title)
	assert.Equal(t, "This is a test manga description.", manga.Description)
	assert.Equal(t, "Test Author", manga.Author)
	assert.Equal(t, "", manga.Status)
	assert.Equal(t, "https://example.com/cover.jpg", manga.CoverURL)
	assert.Equal(t, []string{"Action", "Adventure"}, manga.Tags)
}

func TestMangaReadScraper_GetMangaDetails_NoCoverImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<div class="post-title">
						<h1>Test Manga</h1>
					</div>
					<div class="summary__content">
						<p>This is a test manga description.</p>
					</div>
					<div class="author-content">
						<span>Author:</span>
						<a>Test Author</a>
					</div>
					<div class="post-status">
						<div class="post-content_item">
							<span>Status:</span>
							<span class="summary-content">Ongoing</span>
						</div>
					</div>
					<div class="summary_image">
						<a><img  /></a>
					</div>
					<div class="genres-content">
						<a>Action</a>
						<a>Adventure</a>
					</div>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	scraper := &MangaReadScraper{baseURL: server.URL + "/"}
	manga, err := scraper.GetMangaDetails("test-manga")

	assert.NoError(t, err)
	assert.Equal(t, "Test Manga", manga.Title)
	assert.Equal(t, "This is a test manga description.", manga.Description)
	assert.Equal(t, "Test Author", manga.Author)
	assert.Equal(t, "Ongoing", manga.Status)
	assert.Equal(t, "", manga.CoverURL)
	assert.Equal(t, []string{"Action", "Adventure"}, manga.Tags)
}

func TestMangaReadScraper_GetMangaDetails_MultipleTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<div class="post-title">
						<h1>Test Manga</h1>
					</div>
					<div class="summary__content">
						<p>This is a test manga description.</p>
					</div>
					<div class="author-content">
						<span>Author:</span>
						<a>Test Author</a>
					</div>
					<div class="post-status">
						<div class="post-content_item">
							<span>Status:</span>
							<span class="summary-content">Ongoing</span>
						</div>
					</div>
					<div class="summary_image">
						<a><img src="https://example.com/cover.jpg" /></a>
					</div>
					<div class="genres-content">
						<a>Action</a>
						<a>Adventure</a>
						<a>Comedy</a>
						<a>Drama</a>
						<a>Fantasy</a>
					</div>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	scraper := &MangaReadScraper{baseURL: server.URL + "/"}
	manga, err := scraper.GetMangaDetails("test-manga")

	assert.NoError(t, err)
	assert.Equal(t, "Test Manga", manga.Title)
	assert.Equal(t, "This is a test manga description.", manga.Description)
	assert.Equal(t, "Test Author", manga.Author)
	assert.Equal(t, "Ongoing", manga.Status)
	assert.Equal(t, "https://example.com/cover.jpg", manga.CoverURL)
	assert.Equal(t, []string{"Action", "Adventure", "Comedy", "Drama", "Fantasy"}, manga.Tags)
}

func TestMangaReadScraper_GetMangaDetails_NoTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<div class="post-title">
						<h1>Test Manga</h1>
					</div>
					<div class="summary__content">
						<p>This is a test manga description.</p>
					</div>
					<div class="author-content">
						<span>Author:</span>
						<a>Test Author</a>
					</div>
					<div class="post-status">
						<div class="post-content_item">
							<span>Status:</span>
							<span class="summary-content">Ongoing</span>
						</div>
					</div>
					<div class="summary_image">
						<a><img src="https://example.com/cover.jpg" /></a>
					</div>
					<div class="genres-content">
					</div>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	scraper := &MangaReadScraper{baseURL: server.URL + "/"}
	manga, err := scraper.GetMangaDetails("test-manga")

	assert.NoError(t, err)
	assert.Equal(t, "Test Manga", manga.Title)
	assert.Equal(t, "This is a test manga description.", manga.Description)
	assert.Equal(t, "Test Author", manga.Author)
	assert.Equal(t, "Ongoing", manga.Status)
	assert.Equal(t, "https://example.com/cover.jpg", manga.CoverURL)
	assert.Empty(t, manga.Tags)
}

func TestMangaReadScraper_GetMangaDetails_RealHTML(t *testing.T) {
	scraper := &MangaReadScraper{baseURL: "https://www.mangaread.org/"}
	manga, err := scraper.GetMangaDetails("healing-life-through-camping-in-another-world")

	assert.NoError(t, err)
	assert.Equal(t, "Healing Life Through Camping In Another World", manga.Title)
	assert.Equal(t, "The Star chef, KangHyun, hid in a quite countryside after losing his sense of taste where he found A pathway to another world in his grandfather’s house. Since he was on the run anyway, he planned on enjoying a relaxing camp life, but… the people in the other world keep growing interested in KangHyun! Will KangHyun really be able to heal through experiencing a slow life?", manga.Description)
	assert.Equal(t, "OnGoing", manga.Status)
	assert.Equal(t, "https://www.mangaread.org/wp-content/uploads/2024/04/Read-Manhwa-2-193x278.jpg", manga.CoverURL)
	assert.NotEmpty(t, manga.Tags)
}

func TestNewMangaReadScraper(t *testing.T) {
	scraper := NewMangaReadScraper()
	assert.Equal(t, "https://www.mangaread.org/", scraper.GetBaseUrl())
}

func TestMangaReadScraper_GetLatestUpdates_MultipleUpdates(t *testing.T) {
	scraper := &MangaReadScraper{baseURL: "https://www.mangaread.org/"}
	updates, err := scraper.GetLatestUpdates()

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(updates), 2)

	println("Manga Title:", updates[0].MangaTitle)

	assert.NotEmpty(t, updates[0].MangaTitle)
	assert.NotEmpty(t, updates[0].MangaSlug)
	assert.NotEmpty(t, updates[0].ChapterNumber)
	assert.NotEmpty(t, updates[0].UpdateDate)

	assert.NotEmpty(t, updates[1].MangaTitle)
	assert.NotEmpty(t, updates[1].MangaSlug)
	assert.NotEmpty(t, updates[1].ChapterNumber)
	assert.NotEmpty(t, updates[1].UpdateDate)
}
