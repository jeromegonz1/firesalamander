package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestHomeHandler - Test du handler de la page d'accueil
func TestHomeHandler(t *testing.T) {
	// GIVEN - Une requête GET vers la page d'accueil
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// WHEN - Appel du handler
	homeHandler(w, req)

	// THEN - Vérifier la réponse
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Fire Salamander") {
		t.Error("Expected page to contain 'Fire Salamander'")
	}

	if !strings.Contains(body, "Analysez votre SEO") {
		t.Error("Expected page to contain main heading")
	}
}

// TestAnalyzeHandler - Test du handler d'analyse
func TestAnalyzeHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid URL",
			url:            "https://example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   "Analyse en cours",
		},
		{
			name:           "Missing URL",
			url:            "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "URL manquante",
		},
		{
			name:           "Invalid URL",
			url:            "invalid-url",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "URL invalide",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN - Une requête avec l'URL spécifiée
			reqURL := "/analyze"
			if tt.url != "" {
				reqURL += "?url=" + url.QueryEscape(tt.url)
			}
			
			req := httptest.NewRequest(http.MethodGet, reqURL, nil)
			w := httptest.NewRecorder()

			// WHEN - Appel du handler
			analyzeHandler(w, req)

			// THEN - Vérifier la réponse
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			body := w.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("Expected body to contain '%s', got '%s'", tt.expectedBody, body)
			}
		})
	}
}

// TestResultsHandler - Test du handler des résultats
func TestResultsHandler(t *testing.T) {
	// GIVEN - Une requête GET vers les résultats
	req := httptest.NewRequest(http.MethodGet, "/results?url=https://example.com", nil)
	w := httptest.NewRecorder()

	// WHEN - Appel du handler
	resultsHandler(w, req)

	// THEN - Vérifier la réponse
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Score Global SEO") {
		t.Error("Expected results page to contain SEO score section")
	}

	if !strings.Contains(body, "example.com") {
		t.Error("Expected results page to contain analyzed domain")
	}
}

// TestServer - Test d'intégration du serveur
func TestServer(t *testing.T) {
	// GIVEN - Un serveur de test
	server := setupServer()
	ts := httptest.NewServer(server)
	defer ts.Close()

	// Test routes principales
	routes := []struct {
		path           string
		expectedStatus int
	}{
		{"/", http.StatusOK},
		{"/analyze?url=https://example.com", http.StatusOK},
		{"/results?url=https://example.com", http.StatusOK},
		{"/nonexistent", http.StatusNotFound},
	}

	for _, route := range routes {
		t.Run("Route "+route.path, func(t *testing.T) {
			// WHEN - Requête vers la route
			resp, err := http.Get(ts.URL + route.path)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// THEN - Vérifier le status
			if resp.StatusCode != route.expectedStatus {
				t.Errorf("Expected status %d for %s, got %d", 
					route.expectedStatus, route.path, resp.StatusCode)
			}
		})
	}
}

// TestTemplateData - Test des structures de données pour templates
func TestTemplateData(t *testing.T) {
	// Test HomeData
	homeData := HomeData{
		Title: "Accueil",
		URL:   "",
	}

	if homeData.Title != "Accueil" {
		t.Errorf("Expected title 'Accueil', got '%s'", homeData.Title)
	}

	// Test AnalyzingData
	analyzingData := AnalyzingData{
		Title:    "Analyse",
		URL:      "https://example.com",
		Progress: 50,
	}

	if analyzingData.Progress != 50 {
		t.Errorf("Expected progress 50, got %d", analyzingData.Progress)
	}

	// Test ResultsData
	resultsData := ResultsData{
		Title: "Résultats",
		Analysis: Analysis{
			Domain: "example.com",
			Score:  85,
		},
	}

	if resultsData.Analysis.Score != 85 {
		t.Errorf("Expected score 85, got %d", resultsData.Analysis.Score)
	}
}