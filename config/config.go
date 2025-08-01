package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Branding BrandingConfig `yaml:"branding"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Crawler  CrawlerConfig  `yaml:"crawler"`
	AI       AIConfig       `yaml:"ai"`
}

type AppConfig struct {
	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	Icon      string `yaml:"icon"`
	PoweredBy string `yaml:"powered_by"`
}

type BrandingConfig struct {
	PrimaryColor   string   `yaml:"primary_color"`
	SecondaryColor string   `yaml:"secondary_color"`
	NeutralColors  []string `yaml:"neutral_colors"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Path     string `yaml:"path,omitempty"`
	Host     string `yaml:"host,omitempty"`
	Name     string `yaml:"name,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type CrawlerConfig struct {
	Workers   int    `yaml:"workers"`
	RateLimit string `yaml:"rate_limit"`
	UserAgent string `yaml:"user_agent"`
}

type AIConfig struct {
	Enabled  bool   `yaml:"enabled"`
	MockMode bool   `yaml:"mock_mode,omitempty"`
	APIKey   string `yaml:"api_key,omitempty"`
}

func Load(env string) (*Config, error) {
	var configFile string
	switch env {
	case "production", "prod":
		configFile = "config/config.prod.yaml"
	default:
		configFile = "config/config.dev.yaml"
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s not found", configFile)
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Expand environment variables
	config.Database.Name = os.ExpandEnv(config.Database.Name)
	config.Database.User = os.ExpandEnv(config.Database.User)
	config.Database.Password = os.ExpandEnv(config.Database.Password)
	config.AI.APIKey = os.ExpandEnv(config.AI.APIKey)

	// Ensure data directory exists for SQLite
	if config.Database.Type == "sqlite" && config.Database.Path != "" {
		dir := filepath.Dir(config.Database.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	return &config, nil
}