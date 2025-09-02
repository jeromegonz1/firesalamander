package config

import (
	"fmt"
	"os"
	"time"

	"firesalamander/internal/constants"
	"gopkg.in/yaml.v3"
)

type CrawlerConfig struct {
	Performance Performance `yaml:"performance"`
	Respect     Respect     `yaml:"respect"`
	Limits      Limits      `yaml:"limits"`
	UserAgent   string      `yaml:"user_agent"`
	Exclusions  Exclusions  `yaml:"exclusions"`
}

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
	MaxURLs  int    `yaml:"max_urls"`
	MaxDepth int    `yaml:"max_depth"`
	Strategy string `yaml:"strategy"`
}

type Exclusions struct {
	Extensions []string `yaml:"extensions"`
	Patterns   []string `yaml:"patterns"`
}

func LoadCrawlerConfig(path string) (*CrawlerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var wrapper struct {
		Crawler CrawlerConfig `yaml:"crawler"`
	}

	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults
	if wrapper.Crawler.Performance.ConcurrentRequests == 0 {
		wrapper.Crawler.Performance.ConcurrentRequests = 5
	}
	if wrapper.Crawler.Performance.RequestTimeout == 0 {
		wrapper.Crawler.Performance.RequestTimeout = 10 * time.Second
	}
	if wrapper.Crawler.Limits.MaxURLs == 0 {
		wrapper.Crawler.Limits.MaxURLs = constants.DefaultMaxURLs
	}
	if wrapper.Crawler.Limits.MaxDepth == 0 {
		wrapper.Crawler.Limits.MaxDepth = 3
	}

	return &wrapper.Crawler, nil
}

// Config represents the main application configuration
type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

// Load loads the main application configuration
func Load() (*Config, error) {
	// Return default config for CI/testing
	return &Config{
		Server: ServerConfig{
			Port: constants.DefaultServerPort,
		},
	}, nil
}