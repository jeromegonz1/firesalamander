package crawler

import (
	"testing"

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