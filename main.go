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

	log.Printf("Fire Salamander MVP starting on port %d...", cfg.Server.Port)
	
	// TODO: Start HTTP server with TDD
	os.Exit(0)
}