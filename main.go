// Package main - Fire Salamander MVP
// Architecte : Claude Code
// Principe : Simplicité et solidité
package main

import (
	"log"
	"os"

	"firesalamander/internal/config"
)

func main() {
	// Load configuration from environment (NO HARDCODING!)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Check if we're in CI mode (just validate and exit)
	if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
		log.Printf("CI mode detected - configuration loaded successfully")
		os.Exit(0)
	}

	log.Printf("🔥 Fire Salamander starting on :%d", cfg.Server.Port)
	
	// TODO: Start HTTP server with TDD
	// For now, just exit cleanly
	log.Printf("Fire Salamander initialized successfully")
	os.Exit(0)
}