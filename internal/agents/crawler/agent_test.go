package crawler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"firesalamander/internal/agents"
	appconfig "firesalamander/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that Crawler implements the Agent interface
func TestCrawlerImplementsAgent(t *testing.T) {
	config := appconfig.CrawlerConfig{
		UserAgent: "Test-Bot/1.0",
		Limits: appconfig.Limits{
			MaxURLs:  5,
			MaxDepth: 2,
		},
		Performance: appconfig.Performance{
			ConcurrentRequests: 2,
			RetryAttempts:      1,
			RequestTimeout:     5 * time.Second,
		},
	}

	crawler := NewCrawler(config)
	
	// Verify that Crawler implements Agent interface
	var agent agents.Agent = crawler
	assert.NotNil(t, agent)
}

// Test Agent.Name() implementation
func TestCrawlerName(t *testing.T) {
	config := appconfig.CrawlerConfig{}
	crawler := NewCrawler(config)
	
	name := crawler.Name()
	assert.Equal(t, "web-crawler", name)
	assert.NotEmpty(t, name)
}

// Test Agent.HealthCheck() implementation
func TestCrawlerHealthCheck(t *testing.T) {
	config := appconfig.CrawlerConfig{
		Limits: appconfig.Limits{
			MaxURLs:  10,
			MaxDepth: 2,
		},
		Performance: appconfig.Performance{
			ConcurrentRequests: 2,
			RequestTimeout:     5 * time.Second,
		},
	}
	crawler := NewCrawler(config)
	
	err := crawler.HealthCheck()
	assert.NoError(t, err)
}

// Test Agent.Process() implementation with valid input
func TestCrawlerProcess(t *testing.T) {
	// Create mock HTTP server
	testHTML := `<!DOCTYPE html>
<html lang="en">
<head>
	<title>Test Page</title>
	<meta name="description" content="Test description">
</head>
<body>
	<h1>Test Title</h1>
	<p>Test content for semantic analysis</p>
	<a href="/page2">Internal Link</a>
</body>
</html>`

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

	config := appconfig.CrawlerConfig{
		UserAgent: "Test-Bot/1.0",
		Limits: appconfig.Limits{
			MaxURLs:  3,
			MaxDepth: 1,
		},
		Performance: appconfig.Performance{
			ConcurrentRequests: 2,
			RetryAttempts:      1,
			RequestTimeout:     5 * time.Second,
		},
	}

	crawler := NewCrawler(config)

	// Create CrawlRequest as input data
	input := CrawlRequest{
		SeedURL:   server.URL,
		OutputDir: "", // Empty for in-memory only
	}

	ctx := context.Background()
	result, err := crawler.Process(ctx, input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "web-crawler", result.AgentName)
	assert.Equal(t, "success", result.Status)
	assert.Greater(t, result.Duration, int64(0))

	// Verify the crawl data is in the result
	assert.Contains(t, result.Data, "crawl_result")
	
	// Extract and verify crawl result
	crawlData, ok := result.Data["crawl_result"]
	assert.True(t, ok)
	
	// Convert to CrawlResult for validation
	jsonData, err := json.Marshal(crawlData)
	require.NoError(t, err)
	
	var crawlResult CrawlResult
	err = json.Unmarshal(jsonData, &crawlResult)
	require.NoError(t, err)
	
	assert.Greater(t, len(crawlResult.Pages), 0)
	assert.Equal(t, len(crawlResult.Pages), crawlResult.Metadata.TotalPages)
}

// Test Agent.Process() with invalid input
func TestCrawlerProcessInvalidInput(t *testing.T) {
	config := appconfig.CrawlerConfig{
		Performance: appconfig.Performance{
			RequestTimeout: 5 * time.Second,
		},
	}
	crawler := NewCrawler(config)

	ctx := context.Background()

	// Test with nil input
	result, err := crawler.Process(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Test with wrong input type
	result, err = crawler.Process(ctx, "invalid input")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// Test Agent.Process() with unreachable URL
func TestCrawlerProcessUnreachableURL(t *testing.T) {
	config := appconfig.CrawlerConfig{
		UserAgent: "Test-Bot/1.0",
		Limits: appconfig.Limits{
			MaxURLs:  1,
			MaxDepth: 0,
		},
		Performance: appconfig.Performance{
			ConcurrentRequests: 1,
			RetryAttempts:      0,
			RequestTimeout:     1 * time.Second,
		},
	}

	crawler := NewCrawler(config)

	input := CrawlRequest{
		SeedURL:   "http://localhost:99999/nonexistent",
		OutputDir: "",
	}

	ctx := context.Background()
	result, err := crawler.Process(ctx, input)

	require.NoError(t, err) // Process should not fail, but return partial results
	assert.NotNil(t, result)
	assert.Equal(t, "web-crawler", result.AgentName)
	// Status might be "partial" or "success" with no pages crawled
	assert.Contains(t, []string{"success", "partial"}, result.Status)
}

// Test context cancellation during Process
func TestCrawlerProcessContextCancellation(t *testing.T) {
	// Create slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond) // Slow response
		w.Write([]byte("<html><body>Test</body></html>"))
	}))
	defer server.Close()

	config := appconfig.CrawlerConfig{
		UserAgent: "Test-Bot/1.0",
		Limits: appconfig.Limits{
			MaxURLs:  1,
			MaxDepth: 0,
		},
		Performance: appconfig.Performance{
			ConcurrentRequests: 1,
			RetryAttempts:      0,
			RequestTimeout:     5 * time.Second,
		},
	}

	crawler := NewCrawler(config)

	input := CrawlRequest{
		SeedURL:   server.URL,
		OutputDir: "",
	}

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result, err := crawler.Process(ctx, input)

	// Should handle cancellation gracefully
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "web-crawler", result.AgentName)
}

// Benchmark Agent.Process() performance
func BenchmarkCrawlerProcess(b *testing.B) {
	// Create simple test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><head><title>Benchmark</title></head><body><p>Content</p></body></html>"))
	}))
	defer server.Close()

	config := appconfig.CrawlerConfig{
		UserAgent: "Benchmark-Bot/1.0",
		Limits: appconfig.Limits{
			MaxURLs:  1,
			MaxDepth: 0,
		},
		Performance: appconfig.Performance{
			ConcurrentRequests: 1,
			RetryAttempts:      0,
			RequestTimeout:     5 * time.Second,
		},
	}

	crawler := NewCrawler(config)
	input := CrawlRequest{
		SeedURL:   server.URL,
		OutputDir: "",
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crawler.Process(ctx, input)
		if err != nil {
			b.Fatal(err)
		}
	}
}