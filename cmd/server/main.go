package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"firesalamander/internal/config"
)

// HomeData - Structure pour la page d'accueil
type HomeData struct {
	Title       string
	URL         string
	CurrentStep string
	Progress    int
}

// AnalyzingData - Structure pour la page d'analyse
type AnalyzingData struct {
	Title       string
	URL         string
	CurrentStep string
	Progress    int
	Analysis    AnalysisProgress
}

// AnalysisProgress - Progression de l'analyse
type AnalysisProgress struct {
	PagesFound     int    `json:"pages_found"`
	PagesAnalyzed  int    `json:"pages_analyzed"`
	IssuesFound    int    `json:"issues_found"`
	EstimatedTime  string `json:"estimated_time"`
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

var templates *template.Template

// loadTemplates - Charger tous les templates
func loadTemplates() error {
	templateDir := filepath.Join(".", "templates")
	
	// Vérifier que le répertoire templates existe
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("templates directory not found: %w", err)
	}

	// Parser tous les templates
	var err error
	templates, err = template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	return nil
}

// homeHandler - Handler pour la page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := HomeData{
		Title:       "Accueil",
		URL:         "",
		CurrentStep: "home",
		Progress:    0,
	}

	// Mode test - utiliser HTML simple
	if templates == nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>Fire Salamander</title></head><body><h1>Analysez votre SEO</h1></body></html>`)
		return
	}

	err := templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur template: %v", err), http.StatusInternalServerError)
		return
	}
}

// analyzeHandler - Handler pour l'analyse
func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'URL depuis les paramètres
	urlParam := r.URL.Query().Get("url")
	if urlParam == "" {
		http.Error(w, "URL manquante", http.StatusBadRequest)
		return
	}

	// Valider l'URL
	parsedURL, err := url.Parse(urlParam)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		http.Error(w, "URL invalide", http.StatusBadRequest)
		return
	}

	// Simuler une progression d'analyse
	data := AnalyzingData{
		Title:       "Analyse",
		URL:         urlParam,
		CurrentStep: "analyzing",
		Progress:    75,
		Analysis: AnalysisProgress{
			PagesFound:     12,
			PagesAnalyzed:  9,
			IssuesFound:    3,
			EstimatedTime:  "45s",
		},
	}

	// Mode test - utiliser HTML simple
	if templates == nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>Analyse en cours</title></head><body><h1>Analyse en cours</h1><p>%s</p></body></html>`, urlParam)
		return
	}

	err = templates.ExecuteTemplate(w, "base.html", data)
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
	if templates == nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>Résultats SEO</title></head><body><h1>Score Global SEO</h1><p>%s</p></body></html>`, extractDomain(urlParam))
		return
	}

	err := templates.ExecuteTemplate(w, "base", data)
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

// setupServer - Configuration du serveur HTTP
func setupServer() *http.ServeMux {
	mux := http.NewServeMux()

	// Routes principales
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		homeHandler(w, r)
	})
	mux.HandleFunc("/analyze", analyzeHandler)
	mux.HandleFunc("/results", resultsHandler)

	return mux
}

func main() {
	// Charger la configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erreur chargement config: %v", err)
	}

	// Charger les templates
	if err := loadTemplates(); err != nil {
		log.Fatalf("Erreur chargement templates: %v", err)
	}

	// Configuration du serveur
	server := setupServer()

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("🔥 Fire Salamander démarré sur http://%s", addr)
	log.Printf("📊 Interface SEO disponible à l'adresse ci-dessus")

	// Démarrage du serveur
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("Erreur serveur: %v", err)
	}
}