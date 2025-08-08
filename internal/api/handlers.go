package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"firesalamander/internal/constants"
	"firesalamander/internal/messages"
)

// sendJSONError - Helper pour envoyer des erreurs JSON
func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	errorJSON := fmt.Sprintf(`{"error":"%s"}`, message)
	http.Error(w, errorJSON, statusCode)
}

// AnalyzeHandler - POST /api/analyze
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendJSONError(w, messages.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Parser la requête JSON
	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, messages.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	// Valider l'URL
	if req.URL == "" {
		sendJSONError(w, messages.ErrURLRequired, http.StatusBadRequest)
		return
	}

	parsedURL, err := url.Parse(req.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		sendJSONError(w, messages.ErrInvalidURL, http.StatusBadRequest)
		return
	}

	// Générer un ID unique
	analysisID := generateAnalysisID()

	// Créer l'état initial de l'analyse
	analysis := &AnalysisState{
		ID:             analysisID,
		URL:            req.URL,
		Status:         constants.StatusProcessing,
		Progress:       constants.DefaultProgressStart,
		PagesFound:     0,
		PagesAnalyzed:  0,
		IssuesFound:    0,
		EstimatedTime:  messages.TimeEstimateCalculating,
		StartTime:      time.Now(),
	}

	// Stocker l'analyse
	Store.Set(analysisID, analysis)

	// Démarrer la simulation en arrière-plan
	go SimulateAnalysis(analysisID)

	// Retourner la réponse
	response := AnalyzeResponse{
		ID:     analysisID,
		Status: constants.StatusProcessing,
	}

	json.NewEncoder(w).Encode(response)
}

// StatusHandler - GET /api/status/{id}
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Extraire l'ID depuis l'URL
	analysisID := extractAnalysisID(r.URL.Path)
	if analysisID == "" {
		http.Error(w, `{"error":"Invalid analysis ID"}`, http.StatusBadRequest)
		return
	}

	// Récupérer l'analyse
	analysis, exists := Store.Get(analysisID)
	if !exists {
		sendJSONError(w, messages.ErrAnalysisNotFound, http.StatusNotFound)
		return
	}

	// Construire la réponse
	response := StatusResponse{
		ID:             analysis.ID,
		URL:            analysis.URL,
		Status:         analysis.Status,
		Progress:       analysis.Progress,
		PagesFound:     analysis.PagesFound,
		PagesAnalyzed:  analysis.PagesAnalyzed,
		IssuesFound:    analysis.IssuesFound,
		EstimatedTime:  analysis.EstimatedTime,
	}

	json.NewEncoder(w).Encode(response)
}

// ResultsHandler - GET /api/results/{id}
func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Extraire l'ID depuis l'URL
	analysisID := extractAnalysisID(r.URL.Path)
	if analysisID == "" {
		http.Error(w, `{"error":"Invalid analysis ID"}`, http.StatusBadRequest)
		return
	}

	// Récupérer l'analyse
	analysis, exists := Store.Get(analysisID)
	if !exists {
		sendJSONError(w, messages.ErrAnalysisNotFound, http.StatusNotFound)
		return
	}

	// Vérifier que l'analyse est terminée
	if analysis.Status != "complete" {
		http.Error(w, `{"error":"Analysis not complete"}`, http.StatusBadRequest)
		return
	}

	// Retourner les résultats ou les générer si non disponibles
	if analysis.Results == nil {
		analysis.Results = GenerateTestResults(analysis.URL)
	}

	json.NewEncoder(w).Encode(analysis.Results)
}

// generateAnalysisID - Générer un ID unique d'analyse
func generateAnalysisID() string {
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("analysis-%s", timestamp)
}

// extractAnalysisID - Extraire l'ID d'analyse depuis le path
func extractAnalysisID(urlPath string) string {
	// Extraire l'ID depuis /api/status/{id} ou /api/results/{id}
	parts := strings.Split(path.Clean(urlPath), "/")
	if len(parts) >= 3 {
		return parts[len(parts)-1]
	}
	return ""
}

// GenerateTestResults - Générer des résultats de test
func GenerateTestResults(analysisURL string) *ResultsResponse {
	parsedURL, _ := url.Parse(analysisURL)
	domain := parsedURL.Host

	return &ResultsResponse{
		Score:      72,
		PagesCount: 47,
		Issues: []ResultIssue{
			{
				Title:       "Balises title manquantes",
				Count:       5,
				Description: "Certaines pages n'ont pas de balise title ou celle-ci est vide.",
				Pages:       []string{"/contact", "/about", "/services", "/blog", "/pricing"},
				Solution:    "Ajoutez une balise title unique et descriptive pour chaque page.",
			},
			{
				Title:       "Images sans attribut alt",
				Count:       12,
				Description: "Des images n'ont pas d'attribut alt pour l'accessibilité.",
				Pages:       []string{"/home", "/gallery", "/products"},
				Solution:    "Ajoutez des attributs alt descriptifs à toutes vos images.",
			},
		},
		Warnings: []ResultWarning{
			{
				Title:       "Meta descriptions trop courtes",
				Count:       8,
				Description: "Certaines meta descriptions font moins de 120 caractères.",
			},
		},
		Analysis: AnalysisResult{
			Domain:         domain,
			Date:           time.Now().Format("02/01/2006"),
			Score:          72,
			PagesAnalyzed:  47,
			AnalysisTime:   "2m 15s",
			CriticalIssues: 2,
			Warnings:       8,
			AISuggestions: []AISuggestion{
				{
					Title:       "Optimisation des mots-clés",
					Description: "Concentrez-vous sur ces mots-clés pour améliorer votre référencement.",
					Keywords:    []string{"SEO", "analyse", "optimisation", domain},
				},
			},
		},
	}
}

// Fonctions utilitaires pour les tests

// CreateTestAnalysis - Créer une analyse pour les tests
func CreateTestAnalysis(id, url string) {
	analysis := &AnalysisState{
		ID:             id,
		URL:            url,
		Status:         "analyzing",
		Progress:       25,
		PagesFound:     10,
		PagesAnalyzed:  3,
		IssuesFound:    2,
		EstimatedTime:  "45s",
		StartTime:      time.Now(),
	}
	Store.Set(id, analysis)
}

// CompleteTestAnalysis - Marquer une analyse comme terminée
func CompleteTestAnalysis(id string) {
	Store.Update(id, func(analysis *AnalysisState) {
		analysis.Status = "complete"
		analysis.Progress = 100
		analysis.Results = GenerateTestResults(analysis.URL)
	})
}

// GetAnalysisProgress - Récupérer le progrès d'une analyse
func GetAnalysisProgress(id string) int {
	if analysis, exists := Store.Get(id); exists {
		return analysis.Progress
	}
	return 0
}