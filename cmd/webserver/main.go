package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// WebServer sert l'interface web MVP
type WebServer struct {
	port       string
	staticDir  string
	templatesDir string
}

// AuditRequest représente une demande d'audit depuis l'interface web
type AuditRequest struct {
	SiteURL   string `json:"siteUrl"`
	AuditType string `json:"auditType"`
	MaxPages  int    `json:"maxPages"`
	Timestamp string `json:"timestamp"`
}

// AuditResponse représente la réponse lors du démarrage d'un audit
type AuditResponse struct {
	AuditID string `json:"auditId"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NewWebServer crée une nouvelle instance du serveur web
func NewWebServer(port string) *WebServer {
	return &WebServer{
		port:         port,
		staticDir:    "web/static",
		templatesDir: "web",
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := NewWebServer(port)
	server.Start()
}

// Start démarre le serveur web
func (ws *WebServer) Start() {
	mux := http.NewServeMux()

	// Servir les fichiers statiques
	fs := http.FileServer(http.Dir(ws.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route pour la page d'accueil
	mux.HandleFunc("/", ws.handleHome)

	// Routes API pour le MVP
	mux.HandleFunc("/api/v1/audits", ws.handleAudits)
	mux.HandleFunc("/api/v1/audits/", ws.handleAuditDetails)

	// Middleware pour les logs et CORS
	handler := ws.loggingMiddleware(ws.corsMiddleware(mux))

	log.Printf("🦎 Fire Salamander Web Interface démarrée sur http://localhost:%s", ws.port)
	log.Printf("📁 Servir les fichiers statiques depuis: %s", ws.staticDir)
	
	if err := http.ListenAndServe(":"+ws.port, handler); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}

// handleHome sert la page d'accueil
func (ws *WebServer) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	indexPath := filepath.Join(ws.templatesDir, "index.html")
	http.ServeFile(w, r, indexPath)
}

// handleAudits gère les requêtes d'audit
func (ws *WebServer) handleAudits(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ws.handleStartAudit(w, r)
	case "GET":
		ws.handleListAudits(w, r)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// handleStartAudit démarre un nouvel audit
func (ws *WebServer) handleStartAudit(w http.ResponseWriter, r *http.Request) {
	var auditReq AuditRequest
	if err := json.NewDecoder(r.Body).Decode(&auditReq); err != nil {
		http.Error(w, "Requête JSON invalide", http.StatusBadRequest)
		return
	}

	// Validation basique
	if auditReq.SiteURL == "" {
		http.Error(w, "URL du site requise", http.StatusBadRequest)
		return
	}

	// Pour le MVP, simuler le démarrage d'audit
	// En production, ceci appellerait l'Orchestrateur V2
	auditID := fmt.Sprintf("aud_%d", time.Now().Unix())
	
	response := AuditResponse{
		AuditID: auditID,
		Status:  "started",
		Message: "Audit démarré avec succès",
	}

	log.Printf("Nouvel audit démarré: %s pour %s", auditID, auditReq.SiteURL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleListAudits liste les audits (pour le MVP, retourne des données simulées)
func (ws *WebServer) handleListAudits(w http.ResponseWriter, r *http.Request) {
	// Pour le MVP, retourner des données d'exemple
	audits := []map[string]interface{}{
		{
			"id":        "aud_1234567890",
			"siteUrl":   "https://example.com",
			"status":    "completed",
			"createdAt": time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		},
		{
			"id":        "aud_1234567891",
			"siteUrl":   "https://test-site.com",
			"status":    "running",
			"createdAt": time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(audits)
}

// handleAuditDetails gère les détails d'un audit spécifique
func (ws *WebServer) handleAuditDetails(w http.ResponseWriter, r *http.Request) {
	auditID := r.URL.Path[len("/api/v1/audits/"):]
	
	if auditID == "" {
		http.Error(w, "ID d'audit requis", http.StatusBadRequest)
		return
	}

	// Pour le MVP, retourner des données simulées
	auditDetails := map[string]interface{}{
		"id":       auditID,
		"status":   "completed",
		"progress": 100,
		"results": map[string]interface{}{
			"summary": map[string]interface{}{
				"totalPages":      15,
				"totalKeywords":   45,
				"brokenLinks":     2,
				"averageSeoScore": 78,
				"duration":        "2m 45s",
			},
			"agents": map[string]interface{}{
				"keyword_extractor": map[string]interface{}{
					"status": "completed",
					"data": map[string]interface{}{
						"keywords_count": 45,
					},
				},
				"technical_auditor": map[string]interface{}{
					"status": "completed",
					"data": map[string]interface{}{
						"performance_score": 82,
						"accessibility_score": 74,
						"seo_score": 78,
					},
				},
				"linking_mapper": map[string]interface{}{
					"status": "completed",
					"data": map[string]interface{}{
						"total_links":    120,
						"internal_links": 85,
						"external_links": 35,
					},
				},
				"broken_links_detector": map[string]interface{}{
					"status": "completed",
					"data": map[string]interface{}{
						"broken_count": 2,
						"total_checked": 120,
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auditDetails)
}

// corsMiddleware ajoute les headers CORS pour le développement
func (ws *WebServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware ajoute les logs des requêtes
func (ws *WebServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrapper pour capturer le status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapper, r)
		
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapper.statusCode, time.Since(start))
	})
}

// responseWriter wrapper pour capturer le status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// HealthCheck endpoint pour vérifier que le serveur fonctionne
func (ws *WebServer) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "fire-salamander-web",
		"version":   "1.0.0-mvp",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}