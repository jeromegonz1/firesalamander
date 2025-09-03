package crawler

import (
	"time"
)

// Use config.CrawlerConfig instead of local type

type Performance struct {
	ConcurrentRequests int           `yaml:"concurrent_requests"`
	RequestTimeout     time.Duration `yaml:"request_timeout"`
	RetryAttempts      int           `yaml:"retry_attempts"`
	CacheTTL           time.Duration `yaml:"cache_ttl"`
}

type Respect struct {
	RobotsTxt  bool `yaml:"robots_txt"`
	CrawlDelay bool `yaml:"crawl_delay"`
	Sitemap    bool `yaml:"sitemap"`
}

type Limits struct {
	MaxURLs int    `yaml:"max_urls"`
	MaxDepth int   `yaml:"max_depth"`
	Strategy string `yaml:"strategy"`
}

type Exclusions struct {
	Extensions []string `yaml:"extensions"`
	Patterns   []string `yaml:"patterns"`
}

type CrawlTask struct {
	URL    string
	Depth  int
	Parent string
}

type PageData struct {
	URL           string            `json:"url"`
	Lang          string            `json:"lang"`
	Title         string            `json:"title"`
	H1            string            `json:"h1"`
	H2            []string          `json:"h2"`
	H3            []string          `json:"h3"`
	Anchors       []Anchor          `json:"anchors"`
	Canonical     string            `json:"canonical"`
	MetaIndex     bool              `json:"meta_index"`
	Depth         int               `json:"depth"`
	OutgoingLinks []string          `json:"outgoing_links"`
	IncomingLinks []string          `json:"incoming_links"`
	Content       string            `json:"content"`
}

type Anchor struct {
	Text string `json:"text"`
	Href string `json:"href"`
}

type CrawlResult struct {
	Pages    []PageData `json:"pages"`
	Metadata Metadata   `json:"metadata"`
}

type Metadata struct {
	TotalPages      int  `json:"total_pages"`
	MaxDepthReached int  `json:"max_depth_reached"`
	DurationMs      int  `json:"duration_ms"`
	RobotsRespected bool `json:"robots_respected"`
	SitemapFound    bool `json:"sitemap_found"`
}

// CrawlRequest represents the input data for the Crawler agent
type CrawlRequest struct {
	SeedURL   string `json:"seed_url"`
	OutputDir string `json:"output_dir,omitempty"`
}