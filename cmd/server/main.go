package main

import (
	"firesalamander/internal/constants"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"firesalamander/internal/api"
	"firesalamander/internal/config"
	"firesalamander/internal/logging"
	"firesalamander/internal/monitoring"
)

// HomeData - Structure pour la page d'accueil
type HomeData struct {
	Title string
	URL   string
}

// AnalyzingData - Structure pour la page d'analyse
type AnalyzingData struct {
	Title    string
	URL      string
	Progress int
	Analysis AnalysisProgress
}

// AnalysisProgress - Progression de l'analyse
type AnalysisProgress struct {
	PagesFound    int    `json:"pages_found"`
	PagesAnalyzed int    `json:"pages_analyzed"`
	IssuesFound   int    `json:"issues_found"`
	EstimatedTime string `json:"estimated_time"`
}

// ResultsData - Structure pour les résultats
type ResultsData struct {
	Title    string
	Analysis Analysis
}

// Analysis - Résultats d'analyse SEO
type Analysis struct {
	Domain         string        `json:"domain"`
	Date           string        `json:"date"`
	Score          int           `json:"score"`
	PagesAnalyzed  int           `json:"pages_analyzed"`
	AnalysisTime   string        `json:"analysis_time"`
	CriticalIssues int           `json:"critical_issues"`
	Warnings       int           `json:"warnings"`
	Issues         []Issue       `json:"issues"`
	WarningsList   []Warning     `json:"warnings_list"`
	AISuggestions  []AISuggestion `json:"ai_suggestions"`
}

// Issue - Problème SEO détecté
type Issue struct {
	Title       string   `json:"title"`
	Count       int      `json:"count"`
	Description string   `json:"description"`
	Pages       []string `json:"pages"`
	Solution    string   `json:"solution"`
}

// Warning - Avertissement SEO
type Warning struct {
	Title       string `json:"title"`
	Count       int    `json:"count"`
	Description string `json:"description"`
}

// AISuggestion - Suggestion IA
type AISuggestion struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
}

var (
	homeTemplate      *template.Template
	analyzingTemplate *template.Template
	resultsTemplate   *template.Template
	appLogger         logging.Logger
)

// loadTemplates - Charger les 3 templates
func loadTemplates() error {
	templateDir := filepath.Join(".", "templates")
	
	// Vérifier que le répertoire templates existe
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("templates directory not found: %w", err)
	}

	var err error
	
	// Charger chaque template individuellement
	homeTemplate, err = template.ParseFiles(filepath.Join(templateDir, "home.html"))
	if err != nil {
		return fmt.Errorf("failed to parse home template: %w", err)
	}

	analyzingTemplate, err = template.ParseFiles(filepath.Join(templateDir, "analyzing.html"))
	if err != nil {
		return fmt.Errorf("failed to parse analyzing template: %w", err)
	}

	resultsTemplate, err = template.ParseFiles(filepath.Join(templateDir, "results.html"))
	if err != nil {
		return fmt.Errorf("failed to parse results template: %w", err)
	}

	return nil
}

// homeHandler - Handler pour la page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := HomeData{
		Title: "Accueil",
		URL:   "",
	}

	// Mode test - utiliser HTML simple
	if homeTemplate == nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>Fire Salamander</title></head><body><h1>Analysez votre SEO</h1></body></html>`)
		return
	}

	err := homeTemplate.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur template: %v", err), http.StatusInternalServerError)
		return
	}
}

// analyzeHandler - Handler pour l'analyse
func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID d'analyse depuis les paramètres
	analysisID := r.URL.Query().Get("id")
	if analysisID == "" {
		// Fallback : récupérer l'URL directement (ancienne méthode)
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			http.Error(w, "ID d'analyse ou URL manquant", http.StatusBadRequest)
			return
		}

		// Valider l'URL pour le fallback
		parsedURL, err := url.Parse(urlParam)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			http.Error(w, "URL invalide", http.StatusBadRequest)
			return
		}

		// Simuler une progression d'analyse (ancien comportement)
		data := AnalyzingData{
			Title:    "Analyse",
			URL:      urlParam,
			Progress: 75,
			Analysis: AnalysisProgress{
				PagesFound:    12,
				PagesAnalyzed: 9,
				IssuesFound:   3,
				EstimatedTime: "45s",
			},
		}

		renderAnalyzingTemplate(w, data)
		return
	}

	// Récupérer l'analyse depuis le RealOrchestrator (SPRINT 5)
	status, err := api.GetRealOrchestrator().GetStatus(analysisID)
	if err != nil {
		http.Error(w, constants.AnalysisNotFound, http.StatusNotFound)
		return
	}

	// Calculer le pourcentage de progression
	progress := 0
	if status.PagesFound > 0 {
		progress = (status.PagesAnalyzed * 100) / status.PagesFound
	}
	
	// Calculer le temps écoulé
	elapsed := time.Since(status.StartTime)
	elapsedStr := fmt.Sprintf("%.0fs", elapsed.Seconds())
	
	// Construire les données pour le template depuis RealOrchestrator
	data := AnalyzingData{
		Title:    "Analyse",
		URL:      status.URL,
		Progress: progress,
		Analysis: AnalysisProgress{
			PagesFound:    status.PagesFound,
			PagesAnalyzed: status.PagesAnalyzed,
			IssuesFound:   len(status.TopIssues), // Utiliser le nombre d'issues trouvées
			EstimatedTime: elapsedStr,
		},
	}

	renderAnalyzingTemplate(w, data)
}

// renderAnalyzingTemplate - Fonction helper pour rendre le template d'analyse
func renderAnalyzingTemplate(w http.ResponseWriter, data AnalyzingData) {
	// Mode test - utiliser HTML simple
	if analyzingTemplate == nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>Analyse en cours</title></head><body><h1>Analyse en cours</h1><p>%s (%d%%)</p></body></html>`, data.URL, data.Progress)
		return
	}

	err := analyzingTemplate.ExecuteTemplate(w, "analyzing.html", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur template: %v", err), http.StatusInternalServerError)
		return
	}
}

// resultsHandler - Handler pour les résultats
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'URL depuis les paramètres
	urlParam := r.URL.Query().Get("url")
	if urlParam == "" {
		urlParam = "example.com"
	}

	// Simuler des résultats d'analyse
	data := ResultsData{
		Title: "Résultats",
		Analysis: Analysis{
			Domain:         extractDomain(urlParam),
			Date:           time.Now().Format("02/01/2006"),
			Score:          85,
			PagesAnalyzed:  12,
			AnalysisTime:   "1m 23s",
			CriticalIssues: 3,
			Warnings:       5,
			Issues: []Issue{
				{
					Title:       "Balises title manquantes",
					Count:       3,
					Description: "Certaines pages n'ont pas de balise title ou celle-ci est vide.",
					Pages:       []string{"/contact", "/about", "/services"},
					Solution:    "Ajoutez une balise title unique et descriptive pour chaque page.",
				},
				{
					Title:       "Images sans attribut alt",
					Count:       7,
					Description: "Des images n'ont pas d'attribut alt pour l'accessibilité.",
					Pages:       []string{"/home", "/gallery"},
					Solution:    "Ajoutez des attributs alt descriptifs à toutes vos images.",
				},
			},
			WarningsList: []Warning{
				{
					Title:       "Meta descriptions trop courtes",
					Count:       4,
					Description: "Certaines meta descriptions font moins de 120 caractères.",
				},
			},
			AISuggestions: []AISuggestion{
				{
					Title:       "Optimisation des mots-clés",
					Description: "Concentrez-vous sur ces mots-clés pour améliorer votre référencement.",
					Keywords:    []string{"SEO", "analyse", "optimisation", "référencement"},
				},
			},
		},
	}

	// Mode test - utiliser HTML simple
	if resultsTemplate == nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>Résultats SEO</title></head><body><h1>Score Global SEO</h1><p>%s</p></body></html>`, extractDomain(urlParam))
		return
	}

	err := resultsTemplate.ExecuteTemplate(w, "results.html", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur template: %v", err), http.StatusInternalServerError)
		return
	}
}

// extractDomain - Extraire le domaine d'une URL
func extractDomain(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return parsedURL.Host
}

// setupServer - Configuration du serveur HTTP avec logging
func setupServer() http.Handler {
	mux := http.NewServeMux()

	// Routes principales (pages web)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		homeHandler(w, r)
	})
	
	// Route pour la page d'analyse
	mux.HandleFunc("/analyze", analyzeHandler)
	
	// Route pour la page de résultats
	mux.HandleFunc("/results", resultsHandler)
	
	// 🔥🦎 SPRINT 5: REAL Fire Salamander API Routes (ZERO HARDCODING)
	mux.HandleFunc(constants.APIEndpointAnalyze, api.RealAnalyzeHandler)
	mux.HandleFunc(constants.APIEndpointStatus + "/", api.RealStatusHandler) 
	mux.HandleFunc(constants.APIEndpointResults + "/", api.RealResultsHandler)
	
	// 🔥🦎 MONITORING V2.0: Endpoints de surveillance anti-boucle infinie
	mux.HandleFunc("/debug/metrics", monitoring.MetricsHandler)
	mux.HandleFunc("/health", monitoring.HealthHandler)
	mux.HandleFunc("/api/health", monitoring.HealthHandler)

	// Legacy routes avec données fake (pour debug/comparaison)
	mux.HandleFunc("/api/fake/analyze", analyzeHandler)
	mux.HandleFunc("/api/fake/results", resultsHandler)
	mux.HandleFunc("/api/legacy/analyze", api.AnalyzeHandler)
	mux.HandleFunc("/api/legacy/status/", api.StatusHandler)
	mux.HandleFunc("/api/legacy/results/", api.ResultsHandler)

	// Appliquer les middlewares de logging
	var handler http.Handler = mux
	
	// Middleware de logging HTTP (pour toutes les requêtes)
	handler = logging.HTTPLoggingMiddleware(appLogger)(handler)
	
	// Middleware de logging API (pour les requêtes /api/*)  
	handler = logging.APILoggingMiddleware(appLogger)(handler)
	
	// Middleware de métriques de performance
	handler = logging.MetricsMiddleware(appLogger)(handler)
	
	// Middleware de recovery pour capturer les panics
	handler = logging.RecoveryMiddleware(appLogger)(handler)

	return handler
}

func main() {
	// Initialiser le système de logging en premier
	logConfig := logging.LoadConfigFromEnv()
	var err error
	appLogger, err = logging.NewLogger(logConfig)
	if err != nil {
		log.Fatalf("🔥 Erreur initialisation logging: %v", err)
	}
	defer appLogger.Close()
	
	appLogger.Info(constants.LogCategorySystem, constants.LogMsgServerStarting, map[string]interface{}{
		"version": "1.0.0",
		"pid":     os.Getpid(),
	})

	// Charger la configuration
	cfg, err := config.Load()
	if err != nil {
		appLogger.Fatal(constants.LogCategorySystem, "Erreur chargement config", err)
	}

	// Initialiser le RealOrchestrator pour les API réelles
	appLogger.Info(constants.LogCategorySystem, "Initializing RealOrchestrator")
	api.InitRealOrchestrator()

	// Charger les templates
	appLogger.Info(constants.LogCategorySystem, "Loading templates")
	if err := loadTemplates(); err != nil {
		appLogger.Fatal(constants.LogCategorySystem, "Erreur chargement templates", err)
	}

	// Configuration du serveur avec middlewares de logging
	appLogger.Info(constants.LogCategorySystem, "Setting up HTTP server with logging middlewares")
	server := setupServer()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	
	appLogger.Info(constants.LogCategorySystem, constants.LogMsgServerStarted, map[string]interface{}{
		"address": addr,
		"host":    cfg.Server.Host,
		"port":    cfg.Server.Port,
	})

	// Log des URLs disponibles
	appLogger.Info(constants.LogCategorySystem, "Fire Salamander interfaces available", map[string]interface{}{
		"web_interface": fmt.Sprintf("http://%s", addr),
		"api_endpoint":  fmt.Sprintf("http://%s/api", addr),
		"health_check":  fmt.Sprintf("http://%s/api/health", addr),
	})

	// Démarrage du serveur avec gestion d'erreur via logging
	appLogger.Info(constants.LogCategorySystem, "Starting HTTP server", map[string]interface{}{
		"address": addr,
	})
	
	if err := http.ListenAndServe(addr, server); err != nil {
		appLogger.Fatal(constants.LogCategorySystem, "Erreur serveur HTTP", err)
	}
}