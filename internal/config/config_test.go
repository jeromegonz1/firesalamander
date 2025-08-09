package config

import (
	"os"
	"testing"
	
	"firesalamander/internal/constants"
)

// TestLoadConfig_Success tests successful configuration loading from environment
func TestLoadConfig_Success(t *testing.T) {
	// GIVEN - Environment variables set
	originalEnv := preserveEnvironment()
	defer restoreEnvironment(originalEnv)
	
	os.Setenv("PORT", "3000")
	os.Setenv("HOST", constants.TestIP)
	os.Setenv("ENV", "test")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("MAX_PAGES_CRAWL", "50")
	os.Setenv("TIMEOUT_SECONDS", "180")
	
	// WHEN - Loading configuration
	cfg, err := Load()
	
	// THEN - Configuration should be loaded correctly
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if cfg.Server.Port != 3000 {
		t.Errorf("Expected port 3000, got %d", cfg.Server.Port)
	}
	
	if cfg.Server.Host != constants.TestIP {
		t.Errorf("Expected host '%s', got '%s'", constants.TestIP, cfg.Server.Host)
	}
	
	if cfg.Env != "test" {
		t.Errorf("Expected env 'test', got '%s'", cfg.Env)
	}
	
	if cfg.LogLevel != "info" {
		t.Errorf("Expected log level 'info', got '%s'", cfg.LogLevel)
	}
	
	if cfg.MaxPagesCrawl != 50 {
		t.Errorf("Expected max pages crawl 50, got %d", cfg.MaxPagesCrawl)
	}
	
	if cfg.TimeoutSeconds != 180 {
		t.Errorf("Expected timeout 180, got %d", cfg.TimeoutSeconds)
	}
}

// TestLoadConfig_Defaults tests default values when environment variables are not set
func TestLoadConfig_Defaults(t *testing.T) {
	// GIVEN - Clean environment
	originalEnv := preserveEnvironment()
	defer restoreEnvironment(originalEnv)
	clearEnvironment()
	
	// WHEN - Loading configuration
	cfg, err := Load()
	
	// THEN - Should use default values
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Default values as per .env.example
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.Server.Port)
	}
	
	if cfg.Server.Host != constants.TestDomainLocalhost {
		t.Errorf("Expected default host '%s', got '%s'", constants.TestDomainLocalhost, cfg.Server.Host)
	}
	
	if cfg.Env != "development" {
		t.Errorf("Expected default env 'development', got '%s'", cfg.Env)
	}
	
	if cfg.MaxPagesCrawl != 20 {
		t.Errorf("Expected default max pages crawl 20, got %d", cfg.MaxPagesCrawl)
	}
}

// TestLoadConfig_InvalidPort tests error handling for invalid port
func TestLoadConfig_InvalidPort(t *testing.T) {
	// GIVEN - Invalid port in environment
	originalEnv := preserveEnvironment()
	defer restoreEnvironment(originalEnv)
	
	os.Setenv("PORT", "invalid")
	
	// WHEN - Loading configuration
	_, err := Load()
	
	// THEN - Should return error with context
	if err == nil {
		t.Fatal("Expected error for invalid port, got nil")
	}
	
	// Error should be wrapped with context (no naked errors!)
	expectedErrorSubstring := "failed to parse PORT"
	if !containsString(err.Error(), expectedErrorSubstring) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedErrorSubstring, err.Error())
	}
}

// TestLoadConfig_InvalidMaxPages tests error handling for invalid max pages
func TestLoadConfig_InvalidMaxPages(t *testing.T) {
	// GIVEN - Invalid max pages in environment
	originalEnv := preserveEnvironment()
	defer restoreEnvironment(originalEnv)
	
	os.Setenv("MAX_PAGES_CRAWL", "-1")
	
	// WHEN - Loading configuration
	_, err := Load()
	
	// THEN - Should return error (negative values not allowed)
	if err == nil {
		t.Fatal("Expected error for negative max pages, got nil")
	}
	
	expectedErrorSubstring := "max pages crawl must be positive"
	if !containsString(err.Error(), expectedErrorSubstring) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedErrorSubstring, err.Error())
	}
}

// TestConfig_Validate tests configuration validation
func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Host: constants.TestDomainLocalhost,
				},
				App: AppConfig{
					Name: "Test App",
					Version: "1.0.0",
				},
				Crawler: CrawlerConfig{
					UserAgent: "Test Agent",
					Workers: 5,
					RateLimit: "10/s",
				},
				Env:             "development",
				LogLevel:        "info",
				MaxPagesCrawl:   20,
				TimeoutSeconds:  120,
				MaxConcurrent:   5,
				DBPath:          "./test.db",
				ReportsDir:      "./reports",
				EnablePDFExport: true,
			},
			wantErr: false,
		},
		{
			name: "invalid port too low",
			config: Config{
				Server: ServerConfig{Port: 0},
			},
			wantErr: true,
			errMsg:  "port must be between 1 and 65535",
		},
		{
			name: "invalid port too high",
			config: Config{
				Server: ServerConfig{Port: 70000},
			},
			wantErr: true,
			errMsg:  "port must be between 1 and 65535",
		},
		{
			name: "empty host",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Host: "",
				},
			},
			wantErr: true,
			errMsg:  "host cannot be empty",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			
			if tt.wantErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			
			if tt.wantErr && err != nil && !containsString(err.Error(), tt.errMsg) {
				t.Errorf("Expected error to contain '%s', got '%s'", tt.errMsg, err.Error())
			}
		})
	}
}

// Helper functions for test environment management

func preserveEnvironment() map[string]string {
	env := make(map[string]string)
	envVars := []string{"PORT", "HOST", "ENV", "LOG_LEVEL", "MAX_PAGES_CRAWL", "TIMEOUT_SECONDS"}
	
	for _, key := range envVars {
		if value := os.Getenv(key); value != "" {
			env[key] = value
		}
	}
	
	return env
}

func restoreEnvironment(env map[string]string) {
	clearEnvironment()
	for key, value := range env {
		os.Setenv(key, value)
	}
}

func clearEnvironment() {
	envVars := []string{"PORT", "HOST", "ENV", "LOG_LEVEL", "MAX_PAGES_CRAWL", "TIMEOUT_SECONDS"}
	for _, key := range envVars {
		os.Unsetenv(key)
	}
}

func containsString(haystack, needle string) bool {
	return len(haystack) >= len(needle) && 
		   haystack != "" && 
		   needle != "" &&
		   findSubstring(haystack, needle)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}