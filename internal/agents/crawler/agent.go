package crawler

import (
	"context"
	"fmt"
	"time"

	"firesalamander/internal/agents"
)

// Ensure Crawler implements the Agent interface
var _ agents.Agent = (*Crawler)(nil)

// Name returns the agent name
func (c *Crawler) Name() string {
	return "web-crawler"
}

// Process implements the Agent interface for the Crawler
func (c *Crawler) Process(ctx context.Context, data interface{}) (*agents.AgentResult, error) {
	startTime := time.Now()

	// Validate and parse input data
	request, err := c.parseInput(data)
	if err != nil {
		return nil, fmt.Errorf("invalid input data: %w", err)
	}

	// Perform the crawling operation
	crawlResult, err := c.Crawl(ctx, request.SeedURL, request.OutputDir)
	
	// Create agent result even if crawling had issues
	agentResult := &agents.AgentResult{
		AgentName: c.Name(),
		Status:    "success",
		Data: map[string]interface{}{
			"crawl_result": crawlResult,
		},
		Duration: time.Since(startTime).Milliseconds(),
	}

	if err != nil {
		// If there was an error but we got some results, mark as partial success
		if crawlResult != nil && len(crawlResult.Pages) > 0 {
			agentResult.Status = "partial"
			agentResult.Errors = []string{err.Error()}
		} else {
			agentResult.Status = "error"
			agentResult.Errors = []string{err.Error()}
			return agentResult, nil // Return result with error info, don't fail completely
		}
	}

	return agentResult, nil
}

// HealthCheck implements the Agent interface
func (c *Crawler) HealthCheck() error {
	// Check if crawler configuration is valid
	if c.Config.Performance.RequestTimeout <= 0 {
		return fmt.Errorf("invalid request timeout: %v", c.Config.Performance.RequestTimeout)
	}

	if c.Config.Limits.MaxURLs <= 0 {
		return fmt.Errorf("invalid max URLs limit: %d", c.Config.Limits.MaxURLs)
	}

	if c.Config.Limits.MaxDepth < 0 {
		return fmt.Errorf("invalid max depth: %d", c.Config.Limits.MaxDepth)
	}

	if c.Config.Performance.ConcurrentRequests <= 0 {
		return fmt.Errorf("invalid concurrent requests: %d", c.Config.Performance.ConcurrentRequests)
	}

	// Check if HTTP client is initialized
	if c.client == nil {
		return fmt.Errorf("HTTP client not initialized")
	}

	// All checks passed
	return nil
}

// parseInput validates and parses the input data into a CrawlRequest
func (c *Crawler) parseInput(data interface{}) (*CrawlRequest, error) {
	if data == nil {
		return nil, fmt.Errorf("input data is nil")
	}

	request, ok := data.(CrawlRequest)
	if !ok {
		return nil, fmt.Errorf("expected CrawlRequest, got %T", data)
	}

	if request.SeedURL == "" {
		return nil, fmt.Errorf("seed URL is required")
	}

	return &request, nil
}