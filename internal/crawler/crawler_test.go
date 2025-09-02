package crawler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	appconfig "firesalamander/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCrawler(t *testing.T) {
	cfg := appconfig.CrawlerConfig{
		Limits: appconfig.Limits{
			MaxURLs:  10,
			MaxDepth: 2,
		},
		Performance: appconfig.Performance{
			ConcurrentRequests: 2,
		},
	}

	crawler := NewCrawler(cfg)
	assert.NotNil(t, crawler)
	assert.Equal(t, 10, crawler.Config.Limits.MaxURLs)
	assert.Equal(t, 2, crawler.Config.Limits.MaxDepth)
}

func TestShouldCrawlURL(t *testing.T) {
	crawler := NewCrawler(appconfig.CrawlerConfig{
		Exclusions: appconfig.Exclusions{
			Extensions: []string{".pdf", ".jpg"},
			Patterns:   []string{"/admin/", "/wp-admin/"},
		},
	})

	tests := []struct {
		url      string
		expected bool
		name     string
	}{
		{"https://example.com/page", true, "valid page"},
		{"https://example.com/doc.pdf", false, "PDF excluded"},
		{"https://example.com/admin/panel", false, "admin pattern excluded"},
		{"https://example.com/image.jpg", false, "image excluded"},
		{"https://example.com/contact", true, "contact page valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := crawler.ShouldCrawlURL(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{"https://example.com/page#section", "https://example.com/page", "remove fragment"},
		{"https://example.com/page?utm_source=test", "https://example.com/page", "remove UTM params"},
		{"https://example.com/page/", "https://example.com/page", "remove trailing slash"},
		{"https://example.com/page?id=123", "https://example.com/page?id=123", "keep relevant params"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html lang="fr">
<head>
	<title>Test Page</title>
	<meta name="description" content="Test description">
	<link rel="canonical" href="https://example.com/test">
</head>
<body>
	<h1>Main Title</h1>
	<h2>Section 1</h2>
	<h2>Section 2</h2>
	<h3>Subsection</h3>
	<p>Main content text for analysis.</p>
	<a href="/contact">Contact Us</a>
	<a href="/about">About Page</a>
</body>
</html>`

	page, err := ExtractContent("https://example.com/test", html, 1)
	require.NoError(t, err)

	assert.Equal(t, "https://example.com/test", page.URL)
	assert.Equal(t, "fr", page.Lang)
	assert.Equal(t, "Test Page", page.Title)
	assert.Equal(t, "Main Title", page.H1)
	assert.Len(t, page.H2, 2)
	assert.Contains(t, page.H2, "Section 1")
	assert.Contains(t, page.H2, "Section 2")
	assert.Len(t, page.H3, 1)
	assert.Contains(t, page.H3, "Subsection")
	assert.Equal(t, "https://example.com/test", page.Canonical)
	assert.Equal(t, 1, page.Depth)
	assert.Len(t, page.Anchors, 2)
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		text     string
		expected string
		name     string
	}{
		{"Bonjour, nous sommes une entreprise française", "fr", "French text"},
		{"Hello, we are an English company", "en", "English text"},
		{"", "unknown", "empty text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectLanguage(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRespectDepthLimit(t *testing.T) {
	crawler := NewCrawler(appconfig.CrawlerConfig{
		Limits: appconfig.Limits{MaxDepth: 2},
	})

	assert.True(t, crawler.RespectDepthLimit(0))
	assert.True(t, crawler.RespectDepthLimit(1))
	assert.True(t, crawler.RespectDepthLimit(2))
	assert.False(t, crawler.RespectDepthLimit(3))
}

// Test isInternalLink function to improve coverage
func TestIsInternalLink(t *testing.T) {
	tests := []struct {
		href     string
		baseURL  string
		expected bool
		name     string
	}{
		{"/contact", "https://example.com", true, "relative path"},
		{"https://example.com/page", "https://example.com", true, "same domain"},
		{"https://other.com/page", "https://example.com", false, "external domain"},
		{"", "https://example.com", false, "empty href"},
		{"invalid-url", "invalid-base", true, "invalid base URL treats invalid as relative"},
		{"invalid-url", "https://example.com", false, "invalid href URL"},
		{"https://subdomain.example.com/page", "https://example.com", false, "different subdomain"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInternalLink(tt.href, tt.baseURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test resolveURL function
func TestResolveURL(t *testing.T) {
	crawler := NewCrawler(appconfig.CrawlerConfig{})
	
	tests := []struct {
		baseURL  string
		href     string
		expected string
		name     string
	}{
		{"https://example.com/page", "/contact", "https://example.com/contact", "relative path"},
		{"https://example.com/page", "about", "https://example.com/about", "relative file"},
		{"https://example.com/page", "", "", "empty href"},
		{"invalid-base", "/contact", "/contact", "invalid base URL returns href as-is"},
		{"https://example.com/page", "invalid-href with spaces", "https://example.com/invalid-href%20with%20spaces", "invalid href gets URL encoded"},
		{"https://example.com/page", "https://other.com/external", "", "external domain filtered"},
		{"https://example.com/page?param=1", "../other", "https://example.com/other", "relative with normalization"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := crawler.resolveURL(tt.baseURL, tt.href)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test saveToFile function
func TestSaveToFile(t *testing.T) {
	crawler := NewCrawler(appconfig.CrawlerConfig{})
	
	// Create temp directory
	tempDir, err := ioutil.TempDir("", "crawler_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create test result
	result := &CrawlResult{
		Pages: []PageData{
			{
				URL:   "https://example.com",
				Title: "Test Page",
				Depth: 0,
			},
		},
		Metadata: Metadata{
			TotalPages:      1,
			MaxDepthReached: 0,
			DurationMs:      100,
			RobotsRespected: true,
			SitemapFound:    false,
		},
	}
	
	// Test successful save
	err = crawler.saveToFile(result, tempDir)
	require.NoError(t, err)
	
	// Verify file was created
	filePath := filepath.Join(tempDir, "crawl_index.json")
	assert.FileExists(t, filePath)
	
	// Verify content
	content, err := ioutil.ReadFile(filePath)
	require.NoError(t, err)
	
	var savedResult CrawlResult
	err = json.Unmarshal(content, &savedResult)
	require.NoError(t, err)
	
	assert.Equal(t, result.Pages[0].URL, savedResult.Pages[0].URL)
	assert.Equal(t, result.Metadata.TotalPages, savedResult.Metadata.TotalPages)
}

// Test Crawl function with HTTP server mock
func TestCrawl(t *testing.T) {
	// Create test HTML content
	testHTML := `<!DOCTYPE html>
<html lang="en">
<head>
	<title>Test Page</title>
	<meta name="description" content="Test description">
</head>
<body>
	<h1>Test Title</h1>
	<p>Test content</p>
	<a href="/page2">Internal Link</a>
	<a href="https://external.com">External Link</a>
</body>
</html>`

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(testHTML))
		case "/page2":
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<html><head><title>Page 2</title></head><body><h1>Page 2</h1></body></html>"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create crawler with limited config to avoid infinite crawling
	crawler := NewCrawler(appconfig.CrawlerConfig{
		UserAgent: "Test-Bot/1.0",
		Limits: appconfig.Limits{
			MaxURLs:  5,
			MaxDepth: 1,
		},
		Performance: appconfig.Performance{
			ConcurrentRequests: 2,
			RetryAttempts:      1,
			RequestTimeout:     5 * time.Second,
		},
		Respect: appconfig.Respect{
			RobotsTxt: false,
		},
	})

	// Create temp directory for output
	tempDir, err := ioutil.TempDir("", "crawl_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test crawling
	ctx := context.Background()
	result, err := crawler.Crawl(ctx, server.URL, tempDir)
	
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Greater(t, len(result.Pages), 0)
	assert.Equal(t, len(result.Pages), result.Metadata.TotalPages)
	assert.GreaterOrEqual(t, result.Metadata.MaxDepthReached, 0)
	assert.Greater(t, result.Metadata.DurationMs, 0)
	
	// Verify file was saved
	filePath := filepath.Join(tempDir, "crawl_index.json")
	assert.FileExists(t, filePath)
}

// Test Crawl with error scenarios
func TestCrawlErrorScenarios(t *testing.T) {
	t.Run("server returns 404", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		crawler := NewCrawler(appconfig.CrawlerConfig{
			Limits: appconfig.Limits{MaxURLs: 1, MaxDepth: 0},
			Performance: appconfig.Performance{ConcurrentRequests: 1, RetryAttempts: 0},
		})

		ctx := context.Background()
		result, err := crawler.Crawl(ctx, server.URL, "")
		
		require.NoError(t, err) // Crawl doesn't fail, just logs errors
		assert.Equal(t, 0, len(result.Pages)) // But no pages are added
	})

	t.Run("server returns 500 with retry", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		crawler := NewCrawler(appconfig.CrawlerConfig{
			Limits: appconfig.Limits{MaxURLs: 1, MaxDepth: 0},
			Performance: appconfig.Performance{ConcurrentRequests: 1, RetryAttempts: 2},
		})

		ctx := context.Background()
		result, err := crawler.Crawl(ctx, server.URL, "")
		
		require.NoError(t, err)
		assert.Equal(t, 0, len(result.Pages))
		assert.Greater(t, attempts, 1) // Verify retries occurred
	})

	t.Run("context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond) // Slow response
			w.Write([]byte("<html><body>Test</body></html>"))
		}))
		defer server.Close()

		crawler := NewCrawler(appconfig.CrawlerConfig{
			Limits: appconfig.Limits{MaxURLs: 1, MaxDepth: 0},
			Performance: appconfig.Performance{ConcurrentRequests: 1, RetryAttempts: 0},
		})

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		
		result, err := crawler.Crawl(ctx, server.URL, "")
		
		require.NoError(t, err) // Crawl handles context cancellation gracefully
		assert.Equal(t, 0, len(result.Pages))
	})
}