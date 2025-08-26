package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
)

// APIServer serveur API REST pour Fire Salamander
type APIServer struct {
	orchestrator *Orchestrator
	config       *config.Config
	server       *http.Server
	mux          *http.ServeMux
}

// APIResponse réponse API standardisée
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// AnalysisRequest requête d'analyse
type AnalysisRequest struct {
	URL     string          `json:"url"`
	Type    AnalysisType    `json:"type"`
	Options AnalysisOptions `json:"options"`
}

// HealthResponse réponse de santé du service
type HealthResponse struct {
	Status    string        `json:"status"`
	Version   string        `json:"version"`
	Uptime    time.Duration `json:"uptime"`
	Stats     *AnalysisStats `json:"stats"`
	Timestamp time.Time     `json:"timestamp"`
}

// NewAPIServer crée un nouveau serveur API
func NewAPIServer(orchestrator *Orchestrator, cfg *config.Config) *APIServer {
	mux := http.NewServeMux()
	
	server := &APIServer{
		orchestrator: orchestrator,
		config:       cfg,
		mux:          mux,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      mux,
			ReadTimeout:  constants.ServerReadTimeout,
			WriteTimeout: constants.ServerWriteTimeout,
			IdleTimeout:  constants.ServerIdleTimeout,
		},
	}

	// Enregistrer les routes
	server.registerRoutes()

	return server
}

// registerRoutes enregistre toutes les routes API
func (api *APIServer) registerRoutes() {
	// Middleware CORS et logging
	api.mux.HandleFunc("/", api.withMiddleware(api.handleRoot))
	
	// Routes d'analyse
	api.mux.HandleFunc("/" + constants.APIEndpointV1Analyze, api.withMiddleware(api.handleAnalyze))
	api.mux.HandleFunc("/" + constants.APIEndpointV1AnalyzeSemantic, api.withMiddleware(api.handleSemanticAnalysis))
	api.mux.HandleFunc("/" + constants.APIEndpointV1AnalyzeSEO, api.withMiddleware(api.handleSEOAnalysis))
	api.mux.HandleFunc("/" + constants.APIEndpointV1AnalyzeQuick, api.withMiddleware(api.handleQuickAnalysis))
	
	// Routes de monitoring
	api.mux.HandleFunc("/" + constants.APIEndpointV1Health, api.withMiddleware(api.handleHealth))
	api.mux.HandleFunc("/" + constants.APIEndpointV1Stats, api.withMiddleware(api.handleStats))
	api.mux.HandleFunc("/" + constants.APIEndpointV1Analyses, api.withMiddleware(api.handleAnalyses))
	api.mux.HandleFunc("/" + constants.APIEndpointV1Analysis, api.withMiddleware(api.handleAnalysisDetails))
	
	// Routes utilitaires
	api.mux.HandleFunc("/" + constants.APIEndpointV1Info, api.withMiddleware(api.handleInfo))
	api.mux.HandleFunc("/" + constants.APIEndpointV1Version, api.withMiddleware(api.handleVersion))
}

// withMiddleware applique les middlewares communs
func (api *APIServer) withMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Set(constants.APIHeaderAccessControlAllowOrigin, "*")
		w.Header().Set(constants.APIHeaderAccessControlAllowMethods, "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set(constants.APIHeaderAccessControlAllowHeaders, "Content-Type, Authorization, X-Requested-With")
		
		// Content-Type
		w.Header().Set(constants.APIHeaderContentType, constants.APIContentTypeJSON)
		
		// Version header
		w.Header().Set("X-Fire-Salamander-Version", config.Version())
		
		// Logging des requêtes
		start := time.Now()
		log.Printf("API Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		
		// Gestion OPTIONS pour CORS
		if r.Method == constants.APIMethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Appeler le handler
		handler(w, r)
		
		// Log de fin de requête
		duration := time.Since(start)
		log.Printf("API Response: %s %s completed in %v", r.Method, r.URL.Path, duration)
	}
}

// Start démarre le serveur API
func (api *APIServer) Start() error {
	log.Printf("Démarrage du serveur API Fire Salamander sur le port %d", api.config.Server.Port)
	
	// Démarrer le serveur dans une goroutine
	go func() {
		if err := api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur serveur API: %v", err)
		}
	}()
	
	log.Printf("Serveur API Fire Salamander démarré avec succès")
	log.Printf(constants.APIDocAvailableFormat, api.config.Server.Port)
	
	return nil
}

// Stop arrête le serveur API
func (api *APIServer) Stop(ctx context.Context) error {
	log.Printf("Arrêt du serveur API Fire Salamander")
	return api.server.Shutdown(ctx)
}

// Handlers

// handleRoot affiche la documentation de l'API
func (api *APIServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.APIMethodGet {
		api.sendError(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set(constants.APIHeaderContentType, constants.APIContentTypeHTML)
	html := `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Fire Salamander API</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        h1 { color: #ff6b35; }
        .endpoint { background: #f5f5f5; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .method { font-weight: bold; color: #fff; padding: 2px 8px; border-radius: 3px; }
        .get { background: #28a745; }
        .post { background: #007bff; }
        code { background: #e9ecef; padding: 2px 4px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>🔥 Fire Salamander API</h1>
    <p>API d'analyse SEO et sémantique avancée</p>
    
    <h2>Endpoints d'Analyse</h2>
    
    <div class="endpoint">
        <span class="method post">POST</span> <code>/api/v1/analyze</code>
        <p>Analyse complète (sémantique + SEO + crawling)</p>
        <pre>{constants.APIJSONFieldURL: "` + constants.TestExampleURL + `", constants.APIJSONFieldType: "full", "options": {...}}</pre>
    </div>
    
    <div class="endpoint">
        <span class="method post">POST</span> <code>/api/v1/analyze/semantic</code>
        <p>Analyse sémantique uniquement</p>
    </div>
    
    <div class="endpoint">
        <span class="method post">POST</span> <code>/api/v1/analyze/seo</code>
        <p>Analyse SEO technique uniquement</p>
    </div>
    
    <div class="endpoint">
        <span class="method post">POST</span> <code>/api/v1/analyze/quick</code>
        <p>Analyse rapide (SEO + base)</p>
    </div>
    
    <h2>Endpoints de Monitoring</h2>
    
    <div class="endpoint">
        <span class="method get">GET</span> <code>/api/v1/health</code>
        <p>État de santé du service</p>
    </div>
    
    <div class="endpoint">
        <span class="method get">GET</span> <code>/api/v1/stats</code>
        <p>Statistiques d'utilisation</p>
    </div>
    
    <div class="endpoint">
        <span class="method get">GET</span> <code>/api/v1/info</code>
        <p>Informations sur l'API</p>
    </div>
    
    <h2>Exemple d'utilisation</h2>
    <pre>
` + fmt.Sprintf(constants.APIExampleFormat, api.config.Server.Port) + `
    </pre>
</body>
</html>`
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// handleAnalyze traite les demandes d'analyse complète
func (api *APIServer) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.APIMethodPost {
		api.sendError(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.sendError(w, constants.APIErrorInvalidJSON + ": "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validation
	if req.URL == "" {
		api.sendError(w, "URL requise", http.StatusBadRequest)
		return
	}

	if req.Type == "" {
		req.Type = AnalysisTypeFull
	}

	// Valeurs par défaut pour les options
	if req.Options.Timeout == 0 {
		req.Options.Timeout = constants.LongRequestTimeout
	} else {
		// Si timeout envoyé en millisecondes depuis le frontend, convertir
		if req.Options.Timeout > 1000 {
			req.Options.Timeout = time.Duration(req.Options.Timeout) * time.Millisecond
		}
	}

	// Créer un contexte avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), req.Options.Timeout)
	defer cancel()

	// Lancer l'analyse
	result, err := api.orchestrator.AnalyzeURL(ctx, req.URL, req.Type, req.Options)
	if err != nil {
		api.sendError(w, "Erreur d'analyse: "+err.Error(), http.StatusInternalServerError)
		return
	}

	api.sendSuccess(w, result, "Analyse terminée avec succès")
}

// handleSemanticAnalysis traite les demandes d'analyse sémantique
func (api *APIServer) handleSemanticAnalysis(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.APIMethodPost {
		api.sendError(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.sendError(w, constants.APIErrorInvalidJSON + ": "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		api.sendError(w, "URL requise", http.StatusBadRequest)
		return
	}

	req.Type = AnalysisTypeSemantic
	
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
	defer cancel()

	result, err := api.orchestrator.AnalyzeURL(ctx, req.URL, req.Type, req.Options)
	if err != nil {
		api.sendError(w, "Erreur d'analyse sémantique: "+err.Error(), http.StatusInternalServerError)
		return
	}

	api.sendSuccess(w, result, "Analyse sémantique terminée")
}

// handleSEOAnalysis traite les demandes d'analyse SEO
func (api *APIServer) handleSEOAnalysis(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.APIMethodPost {
		api.sendError(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.sendError(w, constants.APIErrorInvalidJSON + ": "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		api.sendError(w, "URL requise", http.StatusBadRequest)
		return
	}

	req.Type = AnalysisTypeSEO
	
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
	defer cancel()

	result, err := api.orchestrator.AnalyzeURL(ctx, req.URL, req.Type, req.Options)
	if err != nil {
		api.sendError(w, "Erreur d'analyse SEO: "+err.Error(), http.StatusInternalServerError)
		return
	}

	api.sendSuccess(w, result, "Analyse SEO terminée")
}

// handleQuickAnalysis traite les demandes d'analyse rapide
func (api *APIServer) handleQuickAnalysis(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.APIMethodPost {
		api.sendError(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.sendError(w, constants.APIErrorInvalidJSON + ": "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		api.sendError(w, "URL requise", http.StatusBadRequest)
		return
	}

	req.Type = AnalysisTypeQuick
	
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultRequestTimeout)
	defer cancel()

	result, err := api.orchestrator.AnalyzeURL(ctx, req.URL, req.Type, req.Options)
	if err != nil {
		api.sendError(w, "Erreur d'analyse rapide: "+err.Error(), http.StatusInternalServerError)
		return
	}

	api.sendSuccess(w, result, "Analyse rapide terminée")
}

// handleHealth retourne l'état de santé du service
func (api *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.APIMethodGet {
		api.sendError(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	stats := api.orchestrator.GetStats()
	
	health := HealthResponse{
		Status:    constants.APIStatusHealthy,
		Version:   config.Version(),
		Uptime:    time.Since(stats.LastAnalysis),
		Stats:     stats,
		Timestamp: time.Now(),
	}

	// Déterminer le statut basé sur les statistiques
	if stats.TotalTasks > 0 {
		failureRate := float64(stats.FailedTasks) / float64(stats.TotalTasks)
		if failureRate > 0.5 {
			health.Status = constants.APIStatusDegraded
		} else if failureRate > 0.8 {
			health.Status = "unhealthy"
		}
	}

	api.sendSuccess(w, health, "Service opérationnel")
}

// handleStats retourne les statistiques détaillées
func (api *APIServer) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.APIMethodGet {
		api.sendError(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	stats := api.orchestrator.GetStats()
	api.sendSuccess(w, stats, "Statistiques récupérées")
}

// handleAnalyses retourne la liste des analyses récentes
func (api *APIServer) handleAnalyses(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.APIMethodGet {
		api.sendError(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Pour l'instant, retourner les analyses récentes depuis l'orchestrateur
	// TODO: Implémenter la persistance des analyses dans une base de données
	analyses := api.orchestrator.GetRecentAnalyses()
	api.sendSuccess(w, analyses, "Analyses récentes récupérées")
}

// handleAnalysisDetails retourne les détails complets d'une analyse spécifique
func (api *APIServer) handleAnalysisDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.APIMethodGet {
		api.sendError(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Extraire l'ID de l'analyse depuis l'URL
	path := strings.TrimPrefix(r.URL.Path, "/" + constants.APIEndpointV1Analysis)
	if path == "" {
		api.sendError(w, "ID d'analyse requis", http.StatusBadRequest)
		return
	}

	// Convertir l'ID en entier
	analysisID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		api.sendError(w, "ID d'analyse invalide", http.StatusBadRequest)
		return
	}

	// Récupérer les détails de l'analyse depuis le storage
	details, err := api.orchestrator.GetAnalysisDetails(analysisID)
	if err != nil {
		api.sendError(w, "Analyse non trouvée: "+err.Error(), http.StatusNotFound)
		return
	}

	api.sendSuccess(w, details, "Détails de l'analyse récupérés")
}

// handleInfo retourne les informations sur l'API
func (api *APIServer) handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.APIMethodGet {
		api.sendError(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	info := map[string]interface{}{
		constants.APIJSONFieldName:        api.config.App.Name,
		"version":     config.Version(),
		constants.APIJSONFieldDescription: "API d'analyse SEO et sémantique avancée",
		"endpoints": map[string]interface{}{
			"analyze":          "/" + constants.APIEndpointV1Analyze,
			constants.APIAgentSemantic:         "/" + constants.APIEndpointV1AnalyzeSemantic,
			constants.APIAgentSEO:              "/" + constants.APIEndpointV1AnalyzeSEO,
			"quick":            "/" + constants.APIEndpointV1AnalyzeQuick,
			"health":           "/" + constants.APIEndpointV1Health,
			"stats":            "/" + constants.APIEndpointV1Stats,
		},
		"supported_analysis_types": []string{
			string(AnalysisTypeFull),
			string(AnalysisTypeSemantic),
			string(AnalysisTypeSEO),
			string(AnalysisTypeQuick),
		},
		"features": []string{
			"Analyse sémantique hybride",
			"Audit SEO technique complet",
			"Core Web Vitals",
			"Recommandations intelligentes",
			"Scoring unifié",
			"Insights cross-modules",
		},
	}

	api.sendSuccess(w, info, "Informations API")
}

// handleVersion retourne la version de l'API
func (api *APIServer) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.APIMethodGet {
		api.sendError(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	version := map[string]string{
		"version":    config.Version(),
		constants.APIJSONFieldName:       api.config.App.Name,
		"build_time": time.Now().Format("2006-01-02 15:04:05"),
	}

	api.sendSuccess(w, version, "Version de l'API")
}

// Fonctions utilitaires

// sendSuccess envoie une réponse de succès
func (api *APIServer) sendSuccess(w http.ResponseWriter, data interface{}, message string) {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erreur encodage JSON: %v", err)
	}
}

// sendError envoie une réponse d'erreur
func (api *APIServer) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := APIResponse{
		Success:   false,
		Error:     message,
		Timestamp: time.Now(),
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erreur encodage JSON d'erreur: %v", err)
	}
}

// Public handler methods for web server integration

// HandleAnalyze handler public pour l'analyse complète
func (api *APIServer) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	api.handleAnalyze(w, r)
}

// HandleSemanticAnalysis handler public pour l'analyse sémantique
func (api *APIServer) HandleSemanticAnalysis(w http.ResponseWriter, r *http.Request) {
	api.handleSemanticAnalysis(w, r)
}

// HandleSEOAnalysis handler public pour l'analyse SEO
func (api *APIServer) HandleSEOAnalysis(w http.ResponseWriter, r *http.Request) {
	api.handleSEOAnalysis(w, r)
}

// HandleQuickAnalysis handler public pour l'analyse rapide
func (api *APIServer) HandleQuickAnalysis(w http.ResponseWriter, r *http.Request) {
	api.handleQuickAnalysis(w, r)
}

// HandleHealth handler public pour la santé
func (api *APIServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	api.handleHealth(w, r)
}

// HandleStats handler public pour les statistiques
func (api *APIServer) HandleStats(w http.ResponseWriter, r *http.Request) {
	api.handleStats(w, r)
}

// HandleAnalyses handler public pour la liste des analyses
func (api *APIServer) HandleAnalyses(w http.ResponseWriter, r *http.Request) {
	api.handleAnalyses(w, r)
}

// HandleAnalysisDetails handler public pour les détails d'une analyse
func (api *APIServer) HandleAnalysisDetails(w http.ResponseWriter, r *http.Request) {
	api.handleAnalysisDetails(w, r)
}

// HandleInfo handler public pour les informations
func (api *APIServer) HandleInfo(w http.ResponseWriter, r *http.Request) {
	api.handleInfo(w, r)
}

// HandleVersion handler public pour la version
func (api *APIServer) HandleVersion(w http.ResponseWriter, r *http.Request) {
	api.handleVersion(w, r)
}