package crawler

import (
	"context"
	"testing"
	"time"
	
	"firesalamander/internal/config"
)

// TestMustCrawlMultiplePages validates that the crawler actually crawls multiple pages
func TestMustCrawlMultiplePages(t *testing.T) {
	// Simple config - 5 pages, 2 workers
	crawler := NewParallelCrawler(&config.CrawlerConfig{
		MaxPages:         5,
		InitialWorkers:   2,
		MinWorkers:       1,
		TimeoutSeconds:   10,
		UserAgent:        "TestCrawler/1.0",
		RespectRobotsTxt: false, // Skip robots for test
	})
	
	ctx := context.Background()
	report, err := crawler.CrawlWithContext(ctx, "https://example.com")
	
	if err != nil {
		t.Fatalf("Crawl failed: %v", err)
	}
	
	if len(report.Pages) < 2 {
		t.Fatalf("ÉCHEC : Seulement %d pages crawlées, attendu au moins 2", len(report.Pages))
	}
	
	t.Logf("✅ SUCCESS: Crawled %d pages", len(report.Pages))
	for _, page := range report.Pages {
		t.Logf("  - %s (status: %d)", page.URL, page.StatusCode)
	}
}

// TestRaceConditionFixed validates no deadlock occurs
func TestRaceConditionFixed(t *testing.T) {
	crawler := NewParallelCrawler(&config.CrawlerConfig{
		MaxPages:         10,
		InitialWorkers:   5,
		MinWorkers:       1,
		TimeoutSeconds:   5,
		UserAgent:        "TestCrawler/1.0",
		RespectRobotsTxt: false,
	})
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	done := make(chan bool)
	go func() {
		_, _ = crawler.CrawlWithContext(ctx, "https://example.com")
		done <- true
	}()
	
	select {
	case <-done:
		t.Log("✅ No deadlock - crawler completed")
	case <-time.After(6 * time.Second):
		t.Fatal("❌ DEADLOCK DETECTED - crawler stuck")
	}
}