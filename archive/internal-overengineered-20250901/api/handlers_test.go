package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	
	"firesalamander/internal/constants"
)

// TestAnalyzeEndpoint - Test TDD pour démarrer une analyse
func TestAnalyzeEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "Valid URL analysis request",
			requestBody:    constants.TestRequestBodyExample,
			expectedStatus: http.StatusOK,
			expectedFields: []string{"id", "status"},
		},
		{
			name:           "Invalid JSON",
			requestBody:    `invalid json`,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{},
		},
		{
			name:           "Missing URL",
			requestBody:    `{"other":"field"}`,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{},
		},
		{
			name:           "Invalid URL format",
			requestBody:    `{"url":"not-a-valid-url"}`,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN - Une requête POST vers /api/analyze
			req := httptest.NewRequest(http.MethodPost, "/api/analyze", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// WHEN - Appel du handler
			AnalyzeHandler(w, req)

			// THEN - Vérifier le status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// THEN - Vérifier les champs de réponse pour les cas de succès
			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				for _, field := range tt.expectedFields {
					if _, exists := response[field]; !exists {
						t.Errorf("Expected field '%s' in response", field)
					}
				}

				// Vérifier le format de l'ID
				if id, ok := response["id"].(string); ok {
					if !strings.HasPrefix(id, "analysis-") {
						t.Errorf("Expected ID to start with 'analysis-', got '%s'", id)
					}
				}

				// Vérifier le status initial
				if status, ok := response["status"].(string); ok {
					if status != "started" {
						t.Errorf("Expected status 'started', got '%s'", status)
					}
				}
			}
		})
	}
}

// TestStatusEndpoint - Test TDD pour récupérer le status d'une analyse
func TestStatusEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		analysisID     string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "Valid analysis ID",
			analysisID:     "analysis-20240108-143022",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"id", "url", "status", "progress", "pages_found", "pages_analyzed"},
		},
		{
			name:           "Non-existent analysis ID",
			analysisID:     "analysis-nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedFields: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN - Une analyse existante (pour les cas valides)
			if tt.expectedStatus == http.StatusOK {
				// Créer une analyse de test
				CreateTestAnalysis(tt.analysisID, constants.TestURLExample)
			}

			// GIVEN - Une requête GET vers /api/status/{id}
			req := httptest.NewRequest(http.MethodGet, constants.RouteAPIStatus+"/"+tt.analysisID, nil)
			w := httptest.NewRecorder()

			// WHEN - Appel du handler
			StatusHandler(w, req)

			// THEN - Vérifier le status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// THEN - Vérifier les champs pour les cas de succès
			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				for _, field := range tt.expectedFields {
					if _, exists := response[field]; !exists {
						t.Errorf("Expected field '%s' in response", field)
					}
				}

				// Vérifier les types de données
				if progress, ok := response["progress"].(float64); ok {
					if progress < 0 || progress > 100 {
						t.Errorf("Expected progress between 0-100, got %f", progress)
					}
				}
			}
		})
	}
}

// TestResultsEndpoint - Test TDD pour récupérer les résultats d'une analyse
func TestResultsEndpoint(t *testing.T) {
	// GIVEN - Une analyse complétée
	analysisID := "analysis-completed-test"
	CreateTestAnalysis(analysisID, constants.TestURLExample)
	CompleteTestAnalysis(analysisID)

	// GIVEN - Une requête GET vers /api/results/{id}
	req := httptest.NewRequest(http.MethodGet, constants.RouteAPIResults+"/"+analysisID, nil)
	w := httptest.NewRecorder()

	// WHEN - Appel du handler
	ResultsHandler(w, req)

	// THEN - Vérifier le status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// THEN - Vérifier la structure des résultats
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedFields := []string{"score", "pages_count", "issues", "warnings"}
	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Expected field '%s' in response", field)
		}
	}

	// Vérifier le score SEO
	if score, ok := response["score"].(float64); ok {
		if score < 0 || score > 100 {
			t.Errorf("Expected score between 0-100, got %f", score)
		}
	}
}

// TestProgressionSimulation - Test TDD pour la simulation de progression
func TestProgressionSimulation(t *testing.T) {
	// GIVEN - Une nouvelle analyse créée via l'API (pas par CreateTestAnalysis)
	analysisID := "analysis-simulation-test"
	analysis := &AnalysisState{
		ID:             analysisID,
		URL:            constants.TestURLExample,
		Status:         "started",
		Progress:       0, // Commencer à 0 pour ce test
		PagesFound:     0,
		PagesAnalyzed:  0,
		IssuesFound:    0,
		EstimatedTime:  "Calcul en cours...",
	}
	Store.Set(analysisID, analysis)
	
	// THEN - Vérifier que la progression démarre à 0
	initialProgress := GetAnalysisProgress(analysisID)
	if initialProgress != 0 {
		t.Errorf("Expected initial progress 0, got %d", initialProgress)
	}

	// WHEN - Démarrage d'une simulation courte (test seulement)
	// Note: Pour un vrai test, on pourrait démarrer la goroutine et attendre
	// mais pour l'instant on teste juste l'état initial
}