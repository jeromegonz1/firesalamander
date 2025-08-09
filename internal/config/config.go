package config

import (
	"fmt"
	"os"
	"strconv"

	"firesalamander/internal/constants"
)

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int `yaml:"port"`
	Host string `yaml:"host"`
}

// AppConfig holds application configuration
type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

// CrawlerConfig holds crawler configuration for Sprint 3 Parallel Crawler
type CrawlerConfig struct {
	// Basic crawler config
	UserAgent string `yaml:"user_agent"`
	Workers   int    `yaml:"workers"`
	RateLimit string `yaml:"rate_limit"` // Format: "10/s" or "60/m"
	
	// SPRINT 3: Parallel Crawler Configuration
	// Worker Pool Configuration
	MinWorkers     int `yaml:"min_workers"`
	MaxWorkers     int `yaml:"max_workers"`
	InitialWorkers int `yaml:"initial_workers"`
	
	// Performance Adaptation Thresholds
	FastThresholdMs      int `yaml:"fast_threshold_ms"`       // < 500ms = fast site
	SlowThresholdMs      int `yaml:"slow_threshold_ms"`       // > 2000ms = slow site
	ErrorThresholdPercent int `yaml:"error_threshold_percent"` // > 10% errors = problematic
	AdaptIntervalSeconds  int `yaml:"adapt_interval_seconds"`  // Check every 5s
	
	// Crawler Limits
	MaxPages       int `yaml:"max_pages"`
	MaxDepth       int `yaml:"max_depth"`
	TimeoutSeconds int `yaml:"timeout_seconds"`
	DelayMs        int `yaml:"delay_ms"`
	
	// Politeness
	RespectRobotsTxt bool `yaml:"respect_robots_txt"`
}

// Config holds application configuration
// Following Single Responsibility Principle
type Config struct {
	// Nested configurations
	Server  ServerConfig  `yaml:"server"`
	App     AppConfig     `yaml:"app"`
	Crawler CrawlerConfig `yaml:"crawler"`
	
	// Environment
	Env      string `yaml:"env"`
	LogLevel string `yaml:"log_level"`
	
	// SEO Analysis limits (NO HARDCODING!)
	MaxPagesCrawl   int `yaml:"max_pages_crawl"`
	TimeoutSeconds  int `yaml:"timeout_seconds"`
	MaxConcurrent   int `yaml:"max_concurrent"`
	
	// Database
	DBPath string `yaml:"db_path"`
	
	// External services
	OpenAIAPIKey string `yaml:"openai_api_key"`
	OpenAIModel  string `yaml:"openai_model"`
	
	// Reports
	ReportsDir     string `yaml:"reports_dir"`
	EnablePDFExport bool   `yaml:"enable_pdf_export"`
}

// Load loads configuration from environment variables
// Returns error with context (NO PANIC!)
func Load() (*Config, error) {
	cfg := &Config{}
	
	// Load Server configuration
	if port := os.Getenv("PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PORT environment variable: %w", err)
		}
		cfg.Server.Port = p
	} else {
		cfg.Server.Port = constants.DefaultPortInt // Default from constants
	}
	
	if host := os.Getenv("HOST"); host != "" {
		cfg.Server.Host = host
	} else {
		cfg.Server.Host = constants.ServerDefaultHost // Default from .env.example
	}
	
	// Load App configuration
	cfg.App.Name = constants.AppName
	cfg.App.Version = constants.AppVersion
	cfg.App.Description = constants.AppDescription
	
	// Load Crawler configuration
	cfg.Crawler.UserAgent = constants.SEOBotUserAgent
	cfg.Crawler.Workers = 5 // Default
	cfg.Crawler.RateLimit = "10/s" // Default rate limit
	
	// Load ENV with default fallback
	if env := os.Getenv("ENV"); env != "" {
		cfg.Env = env
	} else {
		cfg.Env = "development" // Default from .env.example
	}
	
	// Load LOG_LEVEL with default fallback
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	} else {
		cfg.LogLevel = "debug" // Default from .env.example
	}
	
	// Load MAX_PAGES_CRAWL with default fallback and validation
	if maxPages := os.Getenv("MAX_PAGES_CRAWL"); maxPages != "" {
		mp, err := strconv.Atoi(maxPages)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MAX_PAGES_CRAWL environment variable: %w", err)
		}
		if mp <= 0 {
			return nil, fmt.Errorf("max pages crawl must be positive, got %d", mp)
		}
		cfg.MaxPagesCrawl = mp
	} else {
		cfg.MaxPagesCrawl = 20 // Default from .env.example
	}
	
	// Load TIMEOUT_SECONDS with default fallback
	if timeout := os.Getenv("TIMEOUT_SECONDS"); timeout != "" {
		t, err := strconv.Atoi(timeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse TIMEOUT_SECONDS environment variable: %w", err)
		}
		cfg.TimeoutSeconds = t
	} else {
		cfg.TimeoutSeconds = 120 // Default from .env.example
	}
	
	// Load MAX_CONCURRENT_REQUESTS with default fallback
	if concurrent := os.Getenv("MAX_CONCURRENT_REQUESTS"); concurrent != "" {
		c, err := strconv.Atoi(concurrent)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MAX_CONCURRENT_REQUESTS environment variable: %w", err)
		}
		cfg.MaxConcurrent = c
	} else {
		cfg.MaxConcurrent = 5 // Default from .env.example
	}
	
	// Load DB_PATH with default fallback
	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		cfg.DBPath = dbPath
	} else {
		cfg.DBPath = "./fire_salamander.db" // Default from .env.example
	}
	
	// Load OPENAI_API_KEY (no default for security)
	cfg.OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
	
	// Load OPENAI_MODEL with default fallback
	if model := os.Getenv("OPENAI_MODEL"); model != "" {
		cfg.OpenAIModel = model
	} else {
		cfg.OpenAIModel = "gpt-3.5-turbo" // Default from .env.example
	}
	
	// Load REPORTS_DIR with default fallback
	if reportsDir := os.Getenv("REPORTS_DIR"); reportsDir != "" {
		cfg.ReportsDir = reportsDir
	} else {
		cfg.ReportsDir = "./reports" // Default from .env.example
	}
	
	// Load ENABLE_PDF_EXPORT with default fallback
	if enablePDF := os.Getenv("ENABLE_PDF_EXPORT"); enablePDF != "" {
		cfg.EnablePDFExport = enablePDF == "true"
	} else {
		cfg.EnablePDFExport = true // Default from .env.example
	}
	
	// Validate configuration before returning
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	
	return cfg, nil
}

// Validate validates the configuration
// Returns wrapped error with context
func (c *Config) Validate() error {
	// Validate server configuration
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Server.Port)
	}
	
	if c.Server.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	
	// Validate app configuration
	if c.App.Name == "" {
		return fmt.Errorf("app name cannot be empty")
	}
	
	// Validate crawler configuration
	if c.Crawler.Workers <= 0 {
		return fmt.Errorf("crawler workers must be positive, got %d", c.Crawler.Workers)
	}
	
	if c.Crawler.RateLimit == "" {
		return fmt.Errorf("crawler rate limit cannot be empty")
	}
	
	// Validate positive values
	if c.MaxPagesCrawl <= 0 {
		return fmt.Errorf("max pages crawl must be positive, got %d", c.MaxPagesCrawl)
	}
	
	if c.TimeoutSeconds <= 0 {
		return fmt.Errorf("timeout seconds must be positive, got %d", c.TimeoutSeconds)
	}
	
	if c.MaxConcurrent <= 0 {
		return fmt.Errorf("max concurrent requests must be positive, got %d", c.MaxConcurrent)
	}
	
	// Validate required paths are not empty
	if c.DBPath == "" {
		return fmt.Errorf("database path cannot be empty")
	}
	
	if c.ReportsDir == "" {
		return fmt.Errorf("reports directory cannot be empty")
	}
	
	// Validate environment enum
	validEnvs := map[string]bool{
		"development": true,
		"production":  true,
		"test":        true,
		"staging":     true,
	}
	if !validEnvs[c.Env] {
		return fmt.Errorf("invalid environment '%s', must be one of: development, production, test, staging", c.Env)
	}
	
	// Validate log level enum
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level '%s', must be one of: debug, info, warn, error", c.LogLevel)
	}
	
	return nil
}