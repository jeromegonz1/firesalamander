package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"firesalamander/internal/constants"
)

// TestCrawlerCreation teste la création d'un nouveau crawler
func TestCrawlerCreation(t *testing.T) {
	config := DefaultConfig()
	crawler, err := New(config)

	if err != nil {
		t.Fatalf("Failed to create crawler: %v", err)
	}

	if crawler == nil {
		t.Fatal("Crawler is nil")
	}

	if crawler.config != config {
		t.Error("Config not properly assigned")
	}

	if crawler.fetcher == nil {
		t.Error("Fetcher not initialized")
	}

	if crawler.robotsCache == nil {
		t.Error("Robots cache not initialized")
	}
}

// TestFetcherBasic teste le fetcher HTTP basique
func TestFetcherBasic(t *testing.T) {
	// Créer un serveur de test
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.CrawlerResponseHeaderContentType, constants.CrawlerContentTypeHTML)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "<html><body>Test Page</body></html>")
	}))
	defer server.Close()

	config := DefaultConfig()
	fetcher := NewFetcher(config)

	ctx := context.Background()
	result, err := fetcher.Fetch(ctx, server.URL)

	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", result.StatusCode)
	}

	if !strings.Contains(result.Body, "Test Page") {
		t.Error("Body doesn't contain expected content")
	}

	if result.ContentType != "text/html" {
		t.Errorf("Expected content-type text/html, got %s", result.ContentType)
	}
}

// TestFetcherRetry teste le mécanisme de retry
func TestFetcherRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Success after retries")
	}))
	defer server.Close()

	config := DefaultConfig()
	config.RetryAttempts = 3
	config.RetryDelay = 100 * time.Millisecond
	fetcher := NewFetcher(config)

	ctx := context.Background()
	result, err := fetcher.Fetch(ctx, server.URL)

	if err != nil {
		t.Fatalf("Fetch failed after retries: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected final status 200, got %d", result.StatusCode)
	}
}

// TestRobotsTxtParsing teste le parsing de robots.txt
func TestRobotsTxtParsing(t *testing.T) {
	robotsContent := `# Test robots.txt
User-agent: *
Disallow: /admin/
Disallow: /private/
Allow: /public/

User-agent: FireSalamander
Allow: /
Crawl-delay: 1

User-agent: BadBot
Disallow: /

Sitemap: https://example.com/sitemap.xml
Sitemap: https://example.com/sitemap2.xml`

	robots, err := ParseRobotsTxt(robotsContent)
	if err != nil {
		t.Fatalf("Failed to parse robots.txt: %v", err)
	}

	// Test pour user-agent général
	if robots.IsAllowed("SomeBot", "/public/page.html") != true {
		t.Error("Should allow /public/ for general bots")
	}

	if robots.IsAllowed("SomeBot", "/admin/secret.html") != false {
		t.Error("Should disallow /admin/ for general bots")
	}

	// Test pour FireSalamander
	if robots.IsAllowed("FireSalamander", "/admin/page.html") != true {
		t.Error("Should allow everything for FireSalamander")
	}

	// Test pour BadBot
	if robots.IsAllowed("BadBot", "/any/page.html") != false {
		t.Error("Should disallow everything for BadBot")
	}

	// Test sitemaps
	if len(robots.Sitemaps) != 2 {
		t.Errorf("Expected 2 sitemaps, got %d", len(robots.Sitemaps))
	}

	// Test crawl-delay
	delay := robots.GetCrawlDelay("FireSalamander")
	if delay != 1*time.Second {
		t.Errorf("Expected 1s crawl delay for FireSalamander, got %v", delay)
	}
}

// TestSitemapParsing teste le parsing de sitemap XML
func TestSitemapParsing(t *testing.T) {
	sitemapContent := `<?xml version="1.0" encoding=constants.CrawlerEncodingUTF8?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/page1.html</loc>
    <lastmod>2024-01-01</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
  <url>
    <loc>https://example.com/page2.html</loc>
    <lastmod>2024-01-02</lastmod>
    <changefreq>daily</changefreq>
    <priority>0.6</priority>
  </url>
</urlset>`

	parser := NewSitemapParser()
	sitemap, err := parser.Parse(sitemapContent)

	if err != nil {
		t.Fatalf("Failed to parse sitemap: %v", err)
	}

	if len(sitemap.URLs) != 2 {
		t.Errorf("Expected 2 URLs, got %d", len(sitemap.URLs))
	}

	// Test première URL
	url1 := sitemap.URLs[0]
	if url1.Loc != "https://example.com/page1.html" {
		t.Errorf("Wrong URL location: %s", url1.Loc)
	}
	if url1.Priority != 0.8 {
		t.Errorf("Wrong priority: %f", url1.Priority)
	}
	if url1.Changefreq != "weekly" {
		t.Errorf("Wrong changefreq: %s", url1.Changefreq)
	}
}

// TestPageCache teste le cache de pages
func TestPageCache(t *testing.T) {
	cache := NewPageCache(1 * time.Hour)

	// Test Set et Get
	result := &CrawlResult{
		URL:        "https://example.com/test",
		StatusCode: 200,
		Body:       "Test content",
	}

	cache.Set("test-key", result)

	cached, found := cache.Get("test-key")
	if !found {
		t.Error("Should find cached item")
	}

	if cached.URL != result.URL {
		t.Error("Cached item doesn't match")
	}

	// Test suppression
	cache.Remove("test-key")
	_, found = cache.Get("test-key")
	if found {
		t.Error("Should not find removed item")
	}

	// Test taille
	cache.Set("key1", result)
	cache.Set("key2", result)
	if cache.Size() != 2 {
		t.Errorf("Expected cache size 2, got %d", cache.Size())
	}

	// Test clear
	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("Cache should be empty after clear, got size %d", cache.Size())
	}
}

// TestRateLimiter teste le rate limiter
func TestRateLimiter(t *testing.T) {
	rl, err := NewRateLimiter("5/s")
	if err != nil {
		t.Fatalf("Failed to create rate limiter: %v", err)
	}
	defer rl.Stop()

	ctx := context.Background()
	start := time.Now()

	// Faire 10 requêtes
	for i := 0; i < 10; i++ {
		err := rl.Wait(ctx)
		if err != nil {
			t.Errorf("Rate limiter wait failed: %v", err)
		}
	}

	elapsed := time.Since(start)
	// Avec 5 req/s, 10 requêtes devraient prendre au moins 1 seconde
	if elapsed < 1*time.Second {
		t.Errorf("Rate limiting too fast: %v", elapsed)
	}
}

// TestCrawlQueue teste la queue de crawl
func TestCrawlQueue(t *testing.T) {
	queue := NewCrawlQueue(100)

	// Test ajout
	added := queue.Add("https://example.com/page1", 0)
	if !added {
		t.Error("Should add new URL")
	}

	// Test doublon
	added = queue.Add("https://example.com/page1", 0)
	if added {
		t.Error("Should not add duplicate URL")
	}

	// Test Next
	item, ok := queue.Next()
	if !ok {
		t.Error("Should get item from queue")
	}
	if item.URL != "https://example.com/page1" {
		t.Errorf("Wrong URL from queue: %s", item.URL)
	}

	// Test queue vide
	_, ok = queue.Next()
	if ok {
		t.Error("Should not get item from empty queue")
	}

	// Test HasSeen
	if !queue.HasSeen("https://example.com/page1") {
		t.Error("Should have seen URL")
	}
}

// TestCrawlSiteIntegration teste l'intégration complète
func TestCrawlSiteIntegration(t *testing.T) {
	// Créer un serveur de test avec plusieurs pages
	mux := http.NewServeMux()
	
	// Page principale
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.CrawlerResponseHeaderContentType, constants.CrawlerContentTypeHTML)
		fmt.Fprintf(w, `<html>
			<head><title>Test Site</title></head>
			<body>
				<a href="/page1.html">Page 1</a>
				<a href="/page2.html">Page 2</a>
			</body>
		</html>`)
	})

	// robots.txt
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `User-agent: *
Allow: /
Sitemap: /sitemap.xml`)
	})

	// sitemap.xml
	mux.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.CrawlerResponseHeaderContentType, constants.CrawlerContentTypeXMLApp)
		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>/</loc>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>/page1.html</loc>
    <priority>0.8</priority>
  </url>
</urlset>`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	// Configurer et lancer le crawler
	config := DefaultConfig()
	config.MaxPages = 10
	config.Workers = 2
	
	crawler, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create crawler: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	report, err := crawler.CrawlSite(ctx, server.URL)
	if err != nil {
		t.Fatalf("Crawl failed: %v", err)
	}

	// Vérifier les résultats
	if len(report.Pages) == 0 {
		t.Error("No pages crawled")
	}

	if report.RobotsTxt == nil {
		t.Error("robots.txt not fetched")
	}

	if len(report.Sitemaps) == 0 {
		t.Error("No sitemaps discovered")
	}

	stats := crawler.GetStats()
	if stats.TotalPages == 0 {
		t.Error("No pages in stats")
	}

	t.Logf("Crawled %d pages, %d successful, %d failed",
		stats.TotalPages, stats.SuccessfulPages, stats.FailedPages)
}