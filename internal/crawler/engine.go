package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func (c *Crawler) Crawl(ctx context.Context, seedURL string, outputDir string) (*CrawlResult, error) {
	startTime := time.Now()

	// Initialize
	c.mutex.Lock()
	c.Visited = make(map[string]bool)
	c.Queue = []CrawlTask{{URL: seedURL, Depth: 0}}
	c.Results = make([]PageData, 0)
	c.mutex.Unlock()

	// Create semaphore for concurrent requests
	semaphore := make(chan struct{}, c.Config.Performance.ConcurrentRequests)

	var wg sync.WaitGroup
	maxDepthReached := 0

	for len(c.Queue) > 0 && len(c.Results) < c.Config.Limits.MaxURLs {
		c.mutex.Lock()
		if len(c.Queue) == 0 {
			c.mutex.Unlock()
			break
		}

		task := c.Queue[0]
		c.Queue = c.Queue[1:]
		c.mutex.Unlock()

		// Check if already visited
		c.mutex.RLock()
		if c.Visited[task.URL] {
			c.mutex.RUnlock()
			continue
		}
		c.mutex.RUnlock()

		// Check depth limit
		if !c.RespectDepthLimit(task.Depth) {
			continue
		}

		// Check URL should be crawled
		if !c.ShouldCrawlURL(task.URL) {
			continue
		}

		// Mark as visited
		c.mutex.Lock()
		c.Visited[task.URL] = true
		c.mutex.Unlock()

		if task.Depth > maxDepthReached {
			maxDepthReached = task.Depth
		}

		// Crawl page concurrently
		wg.Add(1)
		go func(task CrawlTask) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire
			defer func() { <-semaphore }() // Release

			if err := c.crawlPage(ctx, task); err != nil {
				fmt.Printf("Error crawling %s: %v\n", task.URL, err)
			}
		}(task)
	}

	wg.Wait()

	// Create result
	result := &CrawlResult{
		Pages: c.Results,
		Metadata: Metadata{
			TotalPages:      len(c.Results),
			MaxDepthReached: maxDepthReached,
			DurationMs:      int(time.Since(startTime).Milliseconds()),
			RobotsRespected: c.Config.Respect.RobotsTxt,
			SitemapFound:    false, // TODO: implement sitemap detection
		},
	}

	// Save to file
	if outputDir != "" {
		if err := c.saveToFile(result, outputDir); err != nil {
			return result, fmt.Errorf("failed to save results: %w", err)
		}
	}

	return result, nil
}

func (c *Crawler) crawlPage(ctx context.Context, task CrawlTask) error {
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", task.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.Config.UserAgent)

	// Make request with retry
	var resp *http.Response
	for attempt := 0; attempt <= c.Config.Performance.RetryAttempts; attempt++ {
		resp, err = c.client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}
		if attempt < c.Config.Performance.RetryAttempts {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", task.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d for %s", resp.StatusCode, task.URL)
	}

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	// Extract content
	page, err := ExtractContent(task.URL, string(body), task.Depth)
	if err != nil {
		return fmt.Errorf("failed to extract content: %w", err)
	}

	// Add to results
	c.mutex.Lock()
	c.Results = append(c.Results, *page)

	// Add new URLs to queue
	for _, anchor := range page.Anchors {
		newURL := c.resolveURL(task.URL, anchor.Href)
		if newURL != "" && !c.Visited[newURL] && c.RespectDepthLimit(task.Depth+1) {
			c.Queue = append(c.Queue, CrawlTask{
				URL:    newURL,
				Depth:  task.Depth + 1,
				Parent: task.URL,
			})
		}
	}
	c.mutex.Unlock()

	return nil
}

func (c *Crawler) resolveURL(baseURL, href string) string {
	if href == "" {
		return ""
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	link, err := url.Parse(href)
	if err != nil {
		return ""
	}

	resolved := base.ResolveReference(link)
	
	// Only return if same host
	if resolved.Host != base.Host {
		return ""
	}

	return NormalizeURL(resolved.String())
}

func (c *Crawler) saveToFile(result *CrawlResult, outputDir string) error {
	// Ensure directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Write crawl_index.json
	filePath := filepath.Join(outputDir, "crawl_index.json")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

