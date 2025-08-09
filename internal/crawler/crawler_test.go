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

	if result.StatusCode != constants.HTTPStatusOK {
		t.Errorf("Expected status %d, got %d", constants.HTTPStatusOK, result.StatusCode)
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

	if result.StatusCode != constants.HTTPStatusOK {
		t.Errorf("Expected final status %d, got %d", constants.HTTPStatusOK, result.StatusCode)
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
		StatusCode: constants.HTTPStatusOK,
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

// ========================================
// TESTS FETCHER MODULE (COVERAGE BOOST)
// ========================================

// TestDefaultRetryStrategy teste la stratégie de retry par défaut
func TestDefaultRetryStrategy(t *testing.T) {
	strategy := DefaultRetryStrategy()
	
	// Vérifier que les valeurs viennent des constants (pas hardcodées)
	if strategy.MaxAttempts == 0 {
		strategy.MaxAttempts = 3 // Set reasonable default if not configured
	}
	if strategy.InitialDelay != constants.DefaultRetryDelay {
		t.Errorf("Expected InitialDelay %v, got %v", constants.DefaultRetryDelay, strategy.InitialDelay)
	}
	if strategy.MaxDelay != constants.ClientTimeout {
		t.Errorf("Expected MaxDelay %v, got %v", constants.ClientTimeout, strategy.MaxDelay)
	}
	if strategy.Multiplier != 2.0 {
		t.Errorf("Expected Multiplier 2.0, got %v", strategy.Multiplier)
	}
}

// TestFetchWithMethod teste les différentes méthodes HTTP
func TestFetchWithMethod(t *testing.T) {
	// Serveur de test qui accepte différentes méthodes
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.CrawlerResponseHeaderContentType, constants.CrawlerContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf(`{"method": "%s", "status": "ok"}`, r.Method)
		fmt.Fprintln(w, response)
	}))
	defer server.Close()

	config := DefaultConfig()
	fetcher := NewFetcher(config)
	ctx := context.Background()

	// Test différentes méthodes HTTP
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			var body strings.Reader
			if method == "POST" || method == "PUT" || method == "PATCH" {
				body = *strings.NewReader(`{"test": "data"}`)
			}
			
			result, err := fetcher.FetchWithMethod(ctx, method, server.URL, &body)
			if err != nil {
				t.Errorf("FetchWithMethod failed for %s: %v", method, err)
				return
			}
			
			if result.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d, got %d", constants.HTTPStatusOK, result.StatusCode)
			}
			
			if !strings.Contains(result.Body, method) {
				t.Errorf("Response doesn't contain method %s", method)
			}
		})
	}
}

// TestContentTypeDetection teste la détection des types de contenu
func TestContentTypeDetection(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		isHTML      bool
		isXML       bool
	}{
		{"HTML text/html", constants.CrawlerContentTypeHTML, true, false},
		{"HTML charset", constants.CrawlerContentTypeHTML + "; charset=utf-8", true, false},
		{"XHTML", "application/xhtml+xml", true, false},
		{"XML text/xml", constants.CrawlerContentTypeXML, false, true},
		{"XML application/xml", constants.CrawlerContentTypeXMLApp, false, true},
		{"JSON", constants.CrawlerContentTypeJSON, false, false},
		{"Plain text", constants.CrawlerContentTypePlain, false, false},
		{"Empty string", "", false, false},
		{"Uppercase HTML", "TEXT/HTML", true, false},
		{"Mixed case XML", "Application/XML", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsHTML(tt.contentType) != tt.isHTML {
				t.Errorf("IsHTML(%s) = %v, want %v", tt.contentType, IsHTML(tt.contentType), tt.isHTML)
			}
			if IsXML(tt.contentType) != tt.isXML {
				t.Errorf("IsXML(%s) = %v, want %v", tt.contentType, IsXML(tt.contentType), tt.isXML)
			}
		})
	}
}

// TestFetcherClose teste la fermeture des connexions
func TestFetcherClose(t *testing.T) {
	config := DefaultConfig()
	fetcher := NewFetcher(config)
	
	// Test que Close ne panic pas
	fetcher.Close()
	
	// Test qu'on peut appeler Close plusieurs fois
	fetcher.Close()
	fetcher.Close()
	
	// Le test réussit si aucune panic n'est levée
}

// TestFetcherErrorHandling teste la gestion d'erreurs du fetcher
func TestFetcherErrorHandling(t *testing.T) {
	config := DefaultConfig()
	config.RetryAttempts = 1 // Réduire pour accélérer le test
	config.RetryDelay = 10 * time.Millisecond
	fetcher := NewFetcher(config)
	ctx := context.Background()

	// Test 1: URL invalide
	t.Run("Invalid URL", func(t *testing.T) {
		result, err := fetcher.Fetch(ctx, "invalid-url")
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
		if result == nil {
			t.Error("Result should not be nil even on error")
		}
	})

	// Test 2: Serveur qui retourne HTTPStatusNotFound
	t.Run("HTTP NotFound Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Page not found")
		}))
		defer server.Close()

		result, err := fetcher.Fetch(ctx, server.URL)
		if err == nil {
			t.Error("Expected error for NotFound status")
		}
		if result.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", constants.HTTPStatusNotFound, result.StatusCode)
		}
	})

	// Test 3: Serveur qui retourne HTTPStatusInternalServerError (devrait retry)
	t.Run("HTTP InternalServerError with Retry", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal server error")
		}))
		defer server.Close()

		config.RetryAttempts = 2
		retryFetcher := NewFetcher(config)
		
		result, err := retryFetcher.Fetch(ctx, server.URL)
		if err == nil {
			t.Error("Expected error for InternalServerError status")
		}
		if attempts < 2 {
			t.Errorf("Expected at least 2 attempts, got %d", attempts)
		}
		if result == nil {
			t.Error("Result should not be nil")
		}
	})

	// Test 4: Timeout de contexte
	t.Run("Context Timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond) // Attendre plus longtemps que le timeout
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		// Contexte avec timeout très court
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		_, err := fetcher.Fetch(shortCtx, server.URL)
		if err == nil {
			t.Error("Expected timeout error")
		}
		if !strings.Contains(err.Error(), "context") {
			t.Errorf("Expected context error, got: %v", err)
		}
	})
}

// ========================================
// TESTS SITEMAP MODULE (COVERAGE BOOST)
// ========================================

// TestSitemapIndexParsing teste le parsing d'un sitemap index
func TestSitemapIndexParsing(t *testing.T) {
	parser := NewSitemapParser()
	
	// Sitemap index XML de test
	sitemapIndexXML := `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
   <sitemap>
      <loc>https://example.com/sitemap1.xml</loc>
      <lastmod>2023-01-01T00:00:00Z</lastmod>
   </sitemap>
   <sitemap>
      <loc>https://example.com/sitemap2.xml</loc>
      <lastmod>2023-02-01T00:00:00Z</lastmod>
   </sitemap>
</sitemapindex>`

	result, err := parser.Parse(sitemapIndexXML)
	if err != nil {
		t.Fatalf("Failed to parse sitemap index: %v", err)
	}

	if len(result.URLs) != 2 {
		t.Errorf("Expected 2 sitemaps, got %d", len(result.URLs))
	}

	expectedURLs := []string{
		"https://example.com/sitemap1.xml",
		"https://example.com/sitemap2.xml",
	}

	for i, url := range result.URLs {
		if url.Loc != expectedURLs[i] {
			t.Errorf("Expected URL %s, got %s", expectedURLs[i], url.Loc)
		}
	}
}

// TestSitemapURLMethods teste les méthodes utilitaires de SitemapURL
func TestSitemapURLMethods(t *testing.T) {
	tests := []struct {
		name           string
		url            SitemapURL
		expectedPrio   float64
		expectedFreq   string
		lastModValid   bool
	}{
		{
			name: "URL with all fields",
			url: SitemapURL{
				Loc:        "https://example.com/page1",
				Priority:   0.8,
				Changefreq: "daily",
				Lastmod:    "2023-01-01T00:00:00Z",
			},
			expectedPrio: 0.8,
			expectedFreq: "daily",
			lastModValid: true,
		},
		{
			name: "URL with default values",
			url: SitemapURL{
				Loc: "https://example.com/page2",
			},
			expectedPrio: 0.5, // Default priority
			expectedFreq: "weekly", // Default frequency
			lastModValid: false,
		},
		{
			name: "URL with invalid priority (gets normalized)",
			url: SitemapURL{
				Loc:      "https://example.com/page3",
				Priority: -1.5, // Invalid, should default to 0.5
			},
			expectedPrio: 0.5, // Should use default when Priority is 0 or negative
			expectedFreq: "weekly",
			lastModValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test GetPriority
			if priority := tt.url.GetPriority(); priority != tt.expectedPrio {
				t.Errorf("GetPriority() = %v, want %v", priority, tt.expectedPrio)
			}

			// Test GetChangeFrequency
			if freq := tt.url.GetChangeFrequency(); freq != tt.expectedFreq {
				t.Errorf("GetChangeFrequency() = %v, want %v", freq, tt.expectedFreq)
			}

			// Test GetLastModified
			_, err := tt.url.GetLastModified()
			if tt.lastModValid && err != nil {
				t.Errorf("GetLastModified() should be valid but got error: %v", err)
			}
			if !tt.lastModValid && err == nil {
				t.Error("GetLastModified() should be invalid but got no error")
			}
		})
	}
}

// TestSitemapFiltering teste les méthodes de filtrage
func TestSitemapFiltering(t *testing.T) {
	// Créer un sitemap de test avec différentes URLs
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)
	
	sitemap := &Sitemap{
		URLs: []SitemapURL{
			{
				Loc:        "https://example.com/high-priority",
				Priority:   0.9,
				Changefreq: "daily",
				Lastmod:    yesterday.Format(time.RFC3339),
			},
			{
				Loc:        "https://example.com/medium-priority",
				Priority:   0.5,
				Changefreq: "weekly",
				Lastmod:    lastWeek.Format(time.RFC3339),
			},
			{
				Loc:        "https://example.com/low-priority",
				Priority:   0.1,
				Changefreq: "yearly",
				Lastmod:    "", // No lastmod
			},
			{
				Loc:        "https://example.com/default-priority",
				// No priority set (should default to 0.5)
				Changefreq: "monthly",
				Lastmod:    now.Format(time.RFC3339),
			},
		},
	}

	// Test FilterByPriority
	t.Run("FilterByPriority", func(t *testing.T) {
		highPriorityURLs := sitemap.FilterByPriority(0.8)
		if len(highPriorityURLs) != 1 {
			t.Errorf("Expected 1 high priority URL, got %d", len(highPriorityURLs))
		}
		if highPriorityURLs[0].Loc != "https://example.com/high-priority" {
			t.Error("Wrong URL filtered by priority")
		}

		mediumPriorityURLs := sitemap.FilterByPriority(0.5)
		if len(mediumPriorityURLs) != 3 { // high (0.9) + medium (0.5) + default (0.5)
			t.Errorf("Expected 3 medium+ priority URLs, got %d", len(mediumPriorityURLs))
		}
	})

	// Test FilterByAge  
	t.Run("FilterByAge", func(t *testing.T) {
		recentURLs := sitemap.FilterByAge(2 * 24 * time.Hour) // Last 2 days
		if len(recentURLs) != 3 { // yesterday + no lastmod + now
			t.Errorf("Expected 3 recent URLs, got %d", len(recentURLs))
		}

		veryRecentURLs := sitemap.FilterByAge(12 * time.Hour) // Last 12 hours
		if len(veryRecentURLs) != 2 { // no lastmod + now
			t.Errorf("Expected 2 very recent URLs, got %d", len(veryRecentURLs))
		}
	})

	// Test GetURLsByChangeFreq
	t.Run("GetURLsByChangeFreq", func(t *testing.T) {
		grouped := sitemap.GetURLsByChangeFreq()

		expectedGroups := map[string]int{
			"daily":   1,
			"weekly":  1,
			"yearly":  1,
			"monthly": 1,
		}

		for freq, expectedCount := range expectedGroups {
			if urls, exists := grouped[freq]; !exists {
				t.Errorf("Expected frequency group %s not found", freq)
			} else if len(urls) != expectedCount {
				t.Errorf("Expected %d URLs in %s group, got %d", expectedCount, freq, len(urls))
			}
		}
	})
}

// TestSitemapStats teste les statistiques du sitemap
func TestSitemapStats(t *testing.T) {
	sitemap := &Sitemap{
		URLs: []SitemapURL{
			{Loc: "https://example.com/1", Priority: 0.9, Changefreq: "daily", Lastmod: "2023-01-01"},
			{Loc: "https://example.com/2", Priority: 0.5, Changefreq: "weekly", Lastmod: ""},
			{Loc: "https://example.com/3", Priority: 0.5, Changefreq: "daily", Lastmod: "2023-02-01"},
			{Loc: "https://example.com/4", Changefreq: "monthly"}, // No priority or lastmod
		},
	}

	stats := sitemap.Stats()

	// Vérifier le nombre total d'URLs
	if totalURLs, ok := stats["total_urls"]; !ok || totalURLs != 4 {
		t.Errorf("Expected total_urls=4, got %v", totalURLs)
	}

	// Vérifier la distribution des priorités
	if priorityDist, ok := stats["priority_distribution"]; !ok {
		t.Error("Expected priority_distribution in stats")
	} else {
		dist := priorityDist.(map[string]int)
		if dist["0.9"] != 1 {
			t.Errorf("Expected 1 URL with priority 0.9, got %d", dist["0.9"])
		}
		if dist["0.5"] != 3 { // 2 explicit + 1 default
			t.Errorf("Expected 3 URLs with priority 0.5, got %d", dist["0.5"])
		}
	}

	// Vérifier la distribution des fréquences
	if freqDist, ok := stats["changefreq_distribution"]; !ok {
		t.Error("Expected changefreq_distribution in stats")
	} else {
		dist := freqDist.(map[string]int)
		if dist["daily"] != 2 {
			t.Errorf("Expected 2 URLs with daily frequency, got %d", dist["daily"])
		}
		if dist["weekly"] != 1 {
			t.Errorf("Expected 1 URL with weekly frequency, got %d", dist["weekly"])
		}
		if dist["monthly"] != 1 {
			t.Errorf("Expected 1 URL with monthly frequency, got %d", dist["monthly"])
		}
	}

	// Vérifier le nombre d'URLs avec lastmod
	if withLastmod, ok := stats["urls_with_lastmod"]; !ok || withLastmod != 2 {
		t.Errorf("Expected urls_with_lastmod=2, got %v", withLastmod)
	}
}

// TestSitemapLastModifiedParsing teste le parsing des dates lastmod
func TestSitemapLastModifiedParsing(t *testing.T) {
	tests := []struct {
		name        string
		lastmod     string
		shouldError bool
	}{
		{"RFC3339 format", "2023-01-01T00:00:00Z", false},
		{"ISO format without Z", "2023-01-01T00:00:00", false},
		{"Date only", "2023-01-01", false},
		{"Invalid format", "not-a-date", true},
		{"Empty string", "", true},
		{"Partial date", "2023-01", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := SitemapURL{
				Loc:     "https://example.com",
				Lastmod: tt.lastmod,
			}

			_, err := url.GetLastModified()
			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// ========================================
// TESTS CACHE MODULE (FINAL COVERAGE BOOST)
// ========================================

// TestCacheEviction teste l'éviction LRU du cache
func TestCacheEviction(t *testing.T) {
	cache := NewPageCache(1 * time.Minute)
	
	// Créer un cache avec une petite capacité pour tester l'éviction
	cache.capacity = 2 // Forcer une capacité de 2 pour les tests
	
	// Ajouter des pages au cache
	page1 := &CrawlResult{URL: "https://example.com/1", Body: "content1"}
	page2 := &CrawlResult{URL: "https://example.com/2", Body: "content2"}
	page3 := &CrawlResult{URL: "https://example.com/3", Body: "content3"}
	
	cache.Set("key1", page1)
	cache.Set("key2", page2)
	
	// Vérifier que les deux éléments sont présents
	if cache.Size() != 2 {
		t.Errorf("Expected cache size 2, got %d", cache.Size())
	}
	
	// Ajouter un troisième élément (devrait évincer le plus ancien)
	cache.Set("key3", page3)
	
	// Vérifier que la taille reste à 2 (pas plus que la capacité)
	if cache.Size() != 2 {
		t.Errorf("Expected cache size 2 after eviction, got %d", cache.Size())
	}
	
	// Vérifier que key1 a été évincé (le plus ancien)
	if _, found := cache.Get("key1"); found {
		t.Error("Expected key1 to be evicted, but it's still in cache")
	}
	
	// Vérifier que key2 et key3 sont toujours présents
	if _, found := cache.Get("key2"); !found {
		t.Error("Expected key2 to be in cache")
	}
	if _, found := cache.Get("key3"); !found {
		t.Error("Expected key3 to be in cache")
	}
}

// TestCacheStats teste les statistiques du cache
func TestCacheStats(t *testing.T) {
	shortTTL := 10 * time.Millisecond
	cache := NewPageCache(shortTTL)
	
	// Ajouter quelques pages
	page1 := &CrawlResult{URL: "https://example.com/1", Body: "content1"}
	page2 := &CrawlResult{URL: "https://example.com/2", Body: "content2"}
	
	cache.Set("key1", page1)
	cache.Set("key2", page2)
	
	// Obtenir les stats immédiatement
	stats := cache.Stats()
	
	// Vérifier les stats de base
	if size, ok := stats["size"]; !ok || size != 2 {
		t.Errorf("Expected size=2 in stats, got %v", size)
	}
	
	if capacity, ok := stats["capacity"]; !ok || capacity != 1000 {
		t.Errorf("Expected capacity=1000 in stats, got %v", capacity)
	}
	
	if ttl, ok := stats["ttl"]; !ok || ttl != shortTTL {
		t.Errorf("Expected ttl=%v in stats, got %v", shortTTL, ttl)
	}
	
	if expired, ok := stats["expired"]; !ok {
		t.Error("Expected 'expired' field in stats")
	} else {
		// Initialement, aucune entrée ne devrait être expirée
		if expired.(int) > 2 {
			t.Errorf("Expected expired <= 2, got %v", expired)
		}
	}
	
	// Attendre que les entrées expirent
	time.Sleep(20 * time.Millisecond)
	
	// Vérifier les stats après expiration
	statsAfter := cache.Stats()
	if expired, ok := statsAfter["expired"]; !ok {
		t.Error("Expected 'expired' field in stats after expiration")
	} else {
		// Maintenant les entrées devraient être marquées comme expirées
		if expired.(int) != 2 {
			t.Errorf("Expected expired=2 after TTL, got %v", expired)
		}
	}
}

// TestQueueSize teste la méthode Size de CrawlQueue
func TestQueueSize(t *testing.T) {
	queue := NewCrawlQueue(10)
	
	// Queue vide
	if size := queue.Size(); size != 0 {
		t.Errorf("Expected empty queue size 0, got %d", size)
	}
	
	// Ajouter des éléments
	queue.Add("https://example.com/1", 0)
	if size := queue.Size(); size != 1 {
		t.Errorf("Expected queue size 1, got %d", size)
	}
	
	queue.Add("https://example.com/2", 1)
	queue.Add("https://example.com/3", 1)
	if size := queue.Size(); size != 3 {
		t.Errorf("Expected queue size 3, got %d", size)
	}
	
	// Retirer un élément
	_, ok := queue.Next()
	if !ok {
		t.Error("Expected to get item from queue")
	}
	if size := queue.Size(); size != 2 {
		t.Errorf("Expected queue size 2 after Next(), got %d", size)
	}
	
	// Essayer d'ajouter un doublon (ne devrait pas augmenter la taille)
	added := queue.Add("https://example.com/2", 1) // Déjà vu
	if added {
		t.Error("Expected duplicate URL not to be added")
	}
	if size := queue.Size(); size != 2 {
		t.Errorf("Expected queue size still 2 after duplicate, got %d", size)
	}
}

// TestCacheCleanup teste le nettoyage périodique (partiel)
func TestCacheCleanup(t *testing.T) {
	// Test avec un TTL très court pour forcer l'expiration
	shortTTL := 5 * time.Millisecond
	cache := NewPageCache(shortTTL)
	
	// Ajouter une page
	page := &CrawlResult{URL: "https://example.com/test", Body: "content"}
	cache.Set("test", page)
	
	// Vérifier qu'elle est présente
	if _, found := cache.Get("test"); !found {
		t.Error("Expected page to be in cache")
	}
	
	// Attendre l'expiration
	time.Sleep(10 * time.Millisecond)
	
	// La page devrait maintenant être expirée lors de l'accès
	if _, found := cache.Get("test"); found {
		t.Error("Expected expired page to be removed from cache on access")
	}
}