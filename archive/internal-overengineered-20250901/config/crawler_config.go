package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"firesalamander/internal/constants"
)

// ========================================
// CRAWLER CONFIGURATION EXTERNALISÉE 
// Conformité Standards SEPTEO - Zéro Hardcoding
// ========================================

// AdvancedCrawlerConfig configuration externalisée pour le crawler
type AdvancedCrawlerConfig struct {
	// Parallel Worker Configuration
	MinWorkers               int           `yaml:"min_workers" env:"CRAWLER_MIN_WORKERS"`
	MaxWorkers              int           `yaml:"max_workers" env:"CRAWLER_MAX_WORKERS"`
	InitialWorkers          int           `yaml:"initial_workers" env:"CRAWLER_INITIAL_WORKERS"`
	
	// Performance Adaptation
	FastThresholdMs         int           `yaml:"fast_threshold_ms" env:"CRAWLER_FAST_THRESHOLD_MS"`
	SlowThresholdMs         int           `yaml:"slow_threshold_ms" env:"CRAWLER_SLOW_THRESHOLD_MS"`
	ErrorThresholdPercent   int           `yaml:"error_threshold_percent" env:"CRAWLER_ERROR_THRESHOLD_PERCENT"`
	AdaptIntervalSeconds    int           `yaml:"adapt_interval_seconds" env:"CRAWLER_ADAPT_INTERVAL_SECONDS"`
	
	// Limits & Timeouts
	MaxPages                int           `yaml:"max_pages" env:"CRAWLER_MAX_PAGES"`
	MaxDepth               int           `yaml:"max_depth" env:"CRAWLER_MAX_DEPTH"`
	TimeoutSeconds         int           `yaml:"timeout_seconds" env:"CRAWLER_TIMEOUT_SECONDS"`
	DelayMs                int           `yaml:"delay_ms" env:"CRAWLER_DELAY_MS"`
	
	// Network Configuration
	UserAgent              string        `yaml:"user_agent" env:"CRAWLER_USER_AGENT"`
	MaxRedirects           int           `yaml:"max_redirects" env:"CRAWLER_MAX_REDIRECTS"`
	MaxBodySizeMB          int           `yaml:"max_body_size_mb" env:"CRAWLER_MAX_BODY_SIZE_MB"`
	RetryAttempts          int           `yaml:"retry_attempts" env:"CRAWLER_RETRY_ATTEMPTS"`
	RetryDelayMs           int           `yaml:"retry_delay_ms" env:"CRAWLER_RETRY_DELAY_MS"`
	
	// Monitoring & Metrics
	MonitoringIntervalMs   int           `yaml:"monitoring_interval_ms" env:"CRAWLER_MONITORING_INTERVAL_MS"`
	MetricsHistorySize     int           `yaml:"metrics_history_size" env:"CRAWLER_METRICS_HISTORY_SIZE"`
	MinMetricsForAdaptation int          `yaml:"min_metrics_for_adaptation" env:"CRAWLER_MIN_METRICS_FOR_ADAPTATION"`
	
	// Politeness & Compliance
	RespectRobotsTxt       bool          `yaml:"respect_robots_txt" env:"CRAWLER_RESPECT_ROBOTS_TXT"`
	FollowSitemaps         bool          `yaml:"follow_sitemaps" env:"CRAWLER_FOLLOW_SITEMAPS"`
	EnableCache            bool          `yaml:"enable_cache" env:"CRAWLER_ENABLE_CACHE"`
	CacheDurationHours     int           `yaml:"cache_duration_hours" env:"CRAWLER_CACHE_DURATION_HOURS"`
	
	// Sitemap Configuration
	SitemapMaxURLs         int           `yaml:"sitemap_max_urls" env:"CRAWLER_SITEMAP_MAX_URLS"`
	SitemapDefaultPriority float64       `yaml:"sitemap_default_priority" env:"CRAWLER_SITEMAP_DEFAULT_PRIORITY"`
	SitemapDefaultChangeFreq string      `yaml:"sitemap_default_change_freq" env:"CRAWLER_SITEMAP_DEFAULT_CHANGE_FREQ"`
}

// LoadAdvancedCrawlerConfig charge la configuration depuis les variables d'environnement
func LoadAdvancedCrawlerConfig() (*AdvancedCrawlerConfig, error) {
	cfg := &AdvancedCrawlerConfig{}
	
	// Worker Configuration
	cfg.MinWorkers = getEnvInt("CRAWLER_MIN_WORKERS", constants.DefaultMinWorkers)
	cfg.MaxWorkers = getEnvInt("CRAWLER_MAX_WORKERS", constants.DefaultMaxWorkers)
	cfg.InitialWorkers = getEnvInt("CRAWLER_INITIAL_WORKERS", constants.DefaultInitialWorkers)
	
	// Performance Thresholds
	cfg.FastThresholdMs = getEnvInt("CRAWLER_FAST_THRESHOLD_MS", constants.DefaultFastThresholdMs)
	cfg.SlowThresholdMs = getEnvInt("CRAWLER_SLOW_THRESHOLD_MS", constants.DefaultSlowThresholdMs)
	cfg.ErrorThresholdPercent = getEnvInt("CRAWLER_ERROR_THRESHOLD_PERCENT", constants.DefaultErrorThresholdPercent)
	cfg.AdaptIntervalSeconds = getEnvInt("CRAWLER_ADAPT_INTERVAL_SECONDS", constants.DefaultAdaptIntervalSeconds)
	
	// Limits & Timeouts
	cfg.MaxPages = getEnvInt("CRAWLER_MAX_PAGES", constants.DefaultMaxPages)
	cfg.MaxDepth = getEnvInt("CRAWLER_MAX_DEPTH", constants.DefaultMaxDepth)
	cfg.TimeoutSeconds = getEnvInt("CRAWLER_TIMEOUT_SECONDS", constants.DefaultTimeoutSeconds)
	cfg.DelayMs = getEnvInt("CRAWLER_DELAY_MS", constants.DefaultDelayMs)
	
	// Network Configuration
	cfg.UserAgent = getEnvString("CRAWLER_USER_AGENT", constants.ParallelCrawlerUserAgent)
	cfg.MaxRedirects = getEnvInt("CRAWLER_MAX_REDIRECTS", constants.DefaultMaxRedirects)
	cfg.MaxBodySizeMB = getEnvInt("CRAWLER_MAX_BODY_SIZE_MB", constants.DefaultMaxBodySize10MB/(1024*1024))
	cfg.RetryAttempts = getEnvInt("CRAWLER_RETRY_ATTEMPTS", constants.DefaultRequestRetries)
	cfg.RetryDelayMs = getEnvInt("CRAWLER_RETRY_DELAY_MS", int(constants.DefaultRetryDelay.Milliseconds()))
	
	// Monitoring Configuration
	cfg.MonitoringIntervalMs = getEnvInt("CRAWLER_MONITORING_INTERVAL_MS", constants.DefaultMonitoringIntervalMs)
	cfg.MetricsHistorySize = getEnvInt("CRAWLER_METRICS_HISTORY_SIZE", constants.DefaultMetricsHistorySize)
	cfg.MinMetricsForAdaptation = getEnvInt("CRAWLER_MIN_METRICS_FOR_ADAPTATION", constants.DefaultMinMetricsForAdaptation)
	
	// Compliance Configuration
	cfg.RespectRobotsTxt = getEnvBool("CRAWLER_RESPECT_ROBOTS_TXT", true)
	cfg.FollowSitemaps = getEnvBool("CRAWLER_FOLLOW_SITEMAPS", true)
	cfg.EnableCache = getEnvBool("CRAWLER_ENABLE_CACHE", true)
	cfg.CacheDurationHours = getEnvInt("CRAWLER_CACHE_DURATION_HOURS", int(constants.RobotsCacheDuration.Hours()))
	
	// Sitemap Configuration
	cfg.SitemapMaxURLs = getEnvInt("CRAWLER_SITEMAP_MAX_URLS", constants.DefaultSitemapMaxURLs)
	cfg.SitemapDefaultPriority = getEnvFloat("CRAWLER_SITEMAP_DEFAULT_PRIORITY", constants.DefaultSitemapPriority)
	cfg.SitemapDefaultChangeFreq = getEnvString("CRAWLER_SITEMAP_DEFAULT_CHANGE_FREQ", constants.DefaultSitemapChangeFreq)
	
	// Valider la configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration crawler invalide: %w", err)
	}
	
	return cfg, nil
}

// Validate valide la configuration
func (c *AdvancedCrawlerConfig) Validate() error {
	// Valider les workers
	if c.MinWorkers <= 0 {
		return fmt.Errorf("min_workers doit être positif, got %d", c.MinWorkers)
	}
	if c.MaxWorkers <= 0 {
		return fmt.Errorf("max_workers doit être positif, got %d", c.MaxWorkers)
	}
	if c.InitialWorkers <= 0 {
		return fmt.Errorf("initial_workers doit être positif, got %d", c.InitialWorkers)
	}
	if c.MinWorkers > c.MaxWorkers {
		return fmt.Errorf("min_workers (%d) ne peut pas être supérieur à max_workers (%d)", 
			c.MinWorkers, c.MaxWorkers)
	}
	if c.InitialWorkers < c.MinWorkers || c.InitialWorkers > c.MaxWorkers {
		return fmt.Errorf("initial_workers (%d) doit être entre min_workers (%d) et max_workers (%d)", 
			c.InitialWorkers, c.MinWorkers, c.MaxWorkers)
	}
	
	// Valider les seuils de performance
	if c.FastThresholdMs <= 0 {
		return fmt.Errorf("fast_threshold_ms doit être positif, got %d", c.FastThresholdMs)
	}
	if c.SlowThresholdMs <= 0 {
		return fmt.Errorf("slow_threshold_ms doit être positif, got %d", c.SlowThresholdMs)
	}
	if c.FastThresholdMs >= c.SlowThresholdMs {
		return fmt.Errorf("fast_threshold_ms (%d) doit être inférieur à slow_threshold_ms (%d)", 
			c.FastThresholdMs, c.SlowThresholdMs)
	}
	
	// Valider les pourcentages
	if c.ErrorThresholdPercent < 0 || c.ErrorThresholdPercent > 100 {
		return fmt.Errorf("error_threshold_percent doit être entre 0 et 100, got %d", c.ErrorThresholdPercent)
	}
	
	// Valider les limites
	if c.MaxPages <= 0 {
		return fmt.Errorf("max_pages doit être positif, got %d", c.MaxPages)
	}
	if c.MaxDepth < 0 {
		return fmt.Errorf("max_depth ne peut pas être négatif, got %d", c.MaxDepth)
	}
	if c.TimeoutSeconds <= 0 {
		return fmt.Errorf("timeout_seconds doit être positif, got %d", c.TimeoutSeconds)
	}
	
	// Valider la configuration réseau
	if c.UserAgent == "" {
		return fmt.Errorf("user_agent ne peut pas être vide")
	}
	if c.MaxRedirects < 0 {
		return fmt.Errorf("max_redirects ne peut pas être négatif, got %d", c.MaxRedirects)
	}
	if c.MaxBodySizeMB <= 0 {
		return fmt.Errorf("max_body_size_mb doit être positif, got %d", c.MaxBodySizeMB)
	}
	
	// Valider la priorité du sitemap
	if c.SitemapDefaultPriority < 0 || c.SitemapDefaultPriority > 1 {
		return fmt.Errorf("sitemap_default_priority doit être entre 0 et 1, got %f", c.SitemapDefaultPriority)
	}
	
	// Valider les fréquences de changement valides
	validChangeFreqs := map[string]bool{
		"always": true, "hourly": true, "daily": true, "weekly": true,
		"monthly": true, "yearly": true, "never": true,
	}
	if !validChangeFreqs[c.SitemapDefaultChangeFreq] {
		return fmt.Errorf("sitemap_default_change_freq invalide '%s', doit être: always, hourly, daily, weekly, monthly, yearly, never", 
			c.SitemapDefaultChangeFreq)
	}
	
	return nil
}

// ToTimeouts convertit la configuration en timeouts utilisables
func (c *AdvancedCrawlerConfig) ToTimeouts() (time.Duration, time.Duration, time.Duration) {
	timeout := time.Duration(c.TimeoutSeconds) * time.Second
	delay := time.Duration(c.DelayMs) * time.Millisecond
	retryDelay := time.Duration(c.RetryDelayMs) * time.Millisecond
	return timeout, delay, retryDelay
}

// ToMonitoringInterval retourne l'intervalle de monitoring
func (c *AdvancedCrawlerConfig) ToMonitoringInterval() time.Duration {
	return time.Duration(c.MonitoringIntervalMs) * time.Millisecond
}

// ToCacheDuration retourne la durée de cache
func (c *AdvancedCrawlerConfig) ToCacheDuration() time.Duration {
	return time.Duration(c.CacheDurationHours) * time.Hour
}

// ========================================
// FONCTIONS UTILITAIRES POUR ENV VARS
// ========================================

// getEnvInt récupère une variable d'environnement int avec valeur par défaut
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvString récupère une variable d'environnement string avec valeur par défaut
func getEnvString(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool récupère une variable d'environnement bool avec valeur par défaut
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}

// getEnvFloat récupère une variable d'environnement float64 avec valeur par défaut
func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}