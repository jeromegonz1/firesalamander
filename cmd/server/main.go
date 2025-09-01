package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"firesalamander/internal/orchestrator"
)

type HomeData struct {
	Title string
	URL   string
}

var (
	homeTemplate *template.Template
	orch         *orchestrator.Orchestrator
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
	auditRequest := orchestrator.AuditRequest{
		SeedURL: req.URL,
		AuditID: auditID,
		Options: make(map[string]interface{}),
	}

	if err := orch.StartAudit(auditRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	status, err := orch.GetStatus(auditID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func main() {
	// Initialize orchestrator
	var err error
	orch, err = orchestrator.NewOrchestrator()
	if err != nil {
		log.Fatalf("Failed to create orchestrator: %v", err)
	}

	// Load templates
	homeTemplate, err = template.ParseFiles(filepath.Join("templates", "home.html"))
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// Setup routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/analyze", analyzeHandler)
	http.HandleFunc("/api/status", statusHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Println("🔥 Fire Salamander starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}