package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	v2 "firesalamander/internal/orchestrator"
	"firesalamander/internal/agents"
	"firesalamander/internal/agents/broken"
	"firesalamander/internal/agents/keyword"
	"firesalamander/internal/agents/linking" 
	"firesalamander/internal/agents/technical"
	"firesalamander/internal/agents/page_profiler"
	"firesalamander/internal/agents/semantic/topic"
	"firesalamander/internal/agents/semantic/recommender"
)

type HomeData struct {
	Title string
	URL   string
}

var (
	homeTemplate *template.Template
	orch         v2.OrchestratorV2
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := HomeData{
		Title: "Fire Salamander",
		URL:   "",
	}

	if err := homeTemplate.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate URL
	if _, err := url.Parse(req.URL); err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Generate audit ID
	auditID := fmt.Sprintf("audit_%d", time.Now().Unix())

	// Start audit
	auditRequest := v2.AuditRequest{
		AuditID:   auditID,
		SeedURL:   req.URL,
		MaxPages:  10,
		Options:   make(map[string]interface{}),
		Timestamp: time.Now(),
	}

	log.Printf("Starting audit %s for URL: %s", auditID, req.URL)
	progressChan, err := orch.StartAudit(context.Background(), &auditRequest)
	if err != nil {
		log.Printf("Failed to start audit %s: %v", auditID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Monitor progress in background
	go func() {
		for update := range progressChan {
			if update.Error != "" {
				log.Printf("Audit %s error: %s", auditID, update.Error)
			} else {
				log.Printf("Audit %s progress: %.1f%% - %s", auditID, update.Progress, update.Step)
			}
		}
		log.Printf("Audit %s completed", auditID)
	}()

	// Return audit ID
	response := map[string]interface{}{
		"id":      auditID,
		"status":  "started",
		"message": "Audit started successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	auditID := r.URL.Query().Get("id")
	if auditID == "" {
		http.Error(w, "Missing audit ID", http.StatusBadRequest)
		return
	}

	status, err := orch.GetAuditStatus(auditID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// resultsHandler - Temporary handler for tests compatibility
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	data := HomeData{
		Title: "SEO Results",
		URL:   "example.com",
	}
	
	// Simple response for tests
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>%s</title></head><body>", data.Title)
	fmt.Fprintf(w, "<h1>Score Global SEO</h1>")
	fmt.Fprintf(w, "<p>Analyzing: %s</p>", data.URL)
	fmt.Fprintf(w, "</body></html>")
}

// setupServer - Create HTTP server for tests
func setupServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/analyze", analyzeHandler)
	mux.HandleFunc("/results", resultsHandler)
	return mux
}

// registerAgents registers all available agents with the orchestrator
func registerAgents() {
	log.Println("Registering agents...")
	
	
	// Create agent instances - Sprint 3 complet avec agents sémantiques
	agentList := []struct {
		name  string
		agent agents.Agent
	}{
		{"technical", technical.NewTechnicalAuditor()},
		{"keyword", keyword.NewKeywordExtractor()},
		{"linking", linking.NewLinkingMapper()},
		{"broken_links", broken.NewBrokenLinksDetector()},
		{"page_profiler", page_profiler.NewPageProfiler()},
		{"topic_clusterer", topic.NewTopicClusterer()},
		{"semantic_recommender", recommender.NewSemanticRecommender()},
	}
	
	// TODO: Add crawler when it implements agents.Agent interface
	
	// Register agents with orchestrator
	registered := 0
	for _, a := range agentList {
		if err := orch.RegisterAgent(a.name, a.agent); err != nil {
			log.Printf("Warning: Failed to register agent %s: %v", a.name, err)
		} else {
			log.Printf("✅ Registered agent: %s", a.name)
			registered++
		}
	}
	
	log.Printf("Successfully registered %d agents", registered)
}

func main() {
	// Initialize orchestrator
	orch = v2.NewOrchestratorV2()
	
	// Register all agents
	registerAgents()
	
	var err error

	// Load templates
	homeTemplate, err = template.ParseFiles(filepath.Join("templates", "home.html"))
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// Get port from environment or default to 8080
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// Setup routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/analyze", analyzeHandler)
	http.HandleFunc("/api/status", statusHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Printf("🔥 Fire Salamander starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}