package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration
// Following Single Responsibility Principle
type Config struct {
	// Server configuration
	Port int
	Host string
	
	// Environment
	Env      string
	LogLevel string
	
	// SEO Analysis limits (NO HARDCODING!)
	MaxPagesCrawl   int
	TimeoutSeconds  int
	MaxConcurrent   int
	
	// Database
	DBPath string
	
	// External services
	OpenAIAPIKey string
	OpenAIModel  string
	
	// Reports
	ReportsDir     string
	EnablePDFExport bool
}

// Load loads configuration from environment variables
// Returns error with context (NO PANIC!)
func Load() (*Config, error) {
	cfg := &Config{}
	
	// Load PORT with default fallback
	if port := os.Getenv("PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PORT environment variable: %w", err)
		}
		cfg.Port = p
	} else {
		cfg.Port = 8080 // Default from .env.example
	}
	
	// Load HOST with default fallback
	if host := os.Getenv("HOST"); host != "" {
		cfg.Host = host
	} else {
		cfg.Host = "localhost" // Default from .env.example
	}
	
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
	// Validate port range
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}
	
	// Validate host is not empty
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
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