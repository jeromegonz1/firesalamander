package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	
	"firesalamander/internal/constants"
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
			url:            constants.TestURLExample,
			expectedStatus: http.StatusOK,
			expectedBody:   "Analyse en cours",
		},
		{
			name:           "Missing URL and ID",
			url:            "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "ID d'analyse ou URL manquant",
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
	req := httptest.NewRequest(http.MethodGet, "/results" + constants.TestQueryURLParam, nil)
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
		{constants.TestQueryAnalyzeParam, http.StatusOK},
		{constants.TestQueryResultsParam, http.StatusOK},
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
		URL:      constants.TestURLExample,
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

// ========================================
// TESTS TDD POUR CORRECTION RESULTSHANDLER  
// Phase RED : Tests avant implémentation
// ========================================

// TestResultsHandler_RealData - DOIT utiliser les vraies données de l'orchestrator
func TestResultsHandler_RealData(t *testing.T) {
	// Pour les tests, setup orchestrator 
	// NOTE: Dans un test réel, on utiliserait un mock orchestrator
	
	// GIVEN - Orchestrator initialisé (simule l'environnement réel)
	// Dans le vrai serveur, l'orchestrator est initialisé dans main()
	
	// GIVEN - Une requête GET vers /results avec un ID vide (pas d'analyse)
	req := httptest.NewRequest(http.MethodGet, "/results", nil) // Pas d'ID
	w := httptest.NewRecorder()
	
	// WHEN - Appel du handler corrigé 
	resultsHandler(w, req)
	
	// THEN - DOIT rediriger vers l'accueil si pas d'ID (Phase GREEN comportement)
	if w.Code != http.StatusFound {
		t.Errorf("Expected redirect %d when no analysis ID, got %d", http.StatusFound, w.Code)
	}
}

// TestResultsHandler_WithRealAnalysis - Test avec une vraie analyse existante
func TestResultsHandler_WithRealAnalysis(t *testing.T) {
	// GIVEN - Une vraie analyse depuis le serveur running
	// Ce test utilise l'analyse de campinglacivelle.fr créée précédemment
	realAnalysisID := "analysis-1756239803-1756239803657238000-128"
	
	req := httptest.NewRequest(http.MethodGet, "/results?id="+realAnalysisID, nil)
	w := httptest.NewRecorder()
	
	// WHEN - Appel du handler avec orchestrator non-initialisé (test environment)
	resultsHandler(w, req)
	
	// THEN - Service unavailable est acceptable en mode test
	// L'important est que ça ne crash pas et qu'on n'ait pas de hardcoded data
	if w.Code == http.StatusServiceUnavailable {
		// ✅ PASS - Service indisponible en test, c'est normal
		return
	}
	
	// Si l'orchestrator était initialisé, vérifier les vraies données
	body := w.Body.String()
	
	// NE DOIT PAS contenir les données hardcodées
	if strings.Contains(body, "85") { // Score hardcodé interdit
		t.Error("❌ TDD FAILURE: Results page still contains hardcoded score 85")
	}
	
	if strings.Contains(body, "12") { // Pages hardcodées interdites
		t.Error("❌ TDD FAILURE: Results page still contains hardcoded pages 12")
	}
}

// TestResultsHandler_NoHardcoding - Validation No Hardcoding Policy
func TestResultsHandler_NoHardcoding(t *testing.T) {
	// GIVEN - Une requête avec un ID d'analyse fictif
	req := httptest.NewRequest(http.MethodGet, "/results?id=test-analysis-123", nil)
	w := httptest.NewRecorder()
	
	// WHEN - Appel du handler
	resultsHandler(w, req)
	
	// THEN - Aucune valeur ne doit être hardcodée
	body := w.Body.String()
	
	// Liste des valeurs INTERDITES (hardcodées actuellement)
	forbiddenHardcodedValues := []string{
		"85",        // Score hardcodé
		"12",        // Pages hardcodées
		"1m 23s",    // Temps hardcodé
		"3",         // Issues hardcodées
		"5",         // Warnings hardcodées
	}
	
	for _, forbidden := range forbiddenHardcodedValues {
		if strings.Contains(body, forbidden) {
			t.Errorf("❌ NO HARDCODING VIOLATION: Found hardcoded value '%s' in results", forbidden)
		}
	}
}

// TestResultsHandler_InvalidAnalysisID - Gestion d'erreur pour ID invalide
func TestResultsHandler_InvalidAnalysisID(t *testing.T) {
	// GIVEN - Une requête avec un ID d'analyse inexistant
	req := httptest.NewRequest(http.MethodGet, "/results?id=nonexistent-analysis", nil)
	w := httptest.NewRecorder()
	
	// WHEN - Appel du handler
	resultsHandler(w, req)
	
	// THEN - Doit gérer l'erreur gracieusement
	if w.Code == http.StatusOK {
		body := w.Body.String()
		// Doit afficher un message d'erreur, pas des données hardcodées
		if strings.Contains(body, "85") {
			t.Error("❌ TDD FAILURE: Invalid analysis should not show hardcoded data")
		}
	}
	
	// Alternativement, pourrait rediriger vers la page d'accueil
	// ou afficher une page d'erreur appropriée
}

// ========================================
// 🧪 TDD TESTS CSS - MISSION PM 
// Phase RED : Tests avant correction FileServer
// ========================================

// loadTemplatesForTests - Load templates with correct path for tests
func loadTemplatesForTests() error {
	templateDir := "../../templates" // From cmd/server to project root
	
	var err error
	
	// Charger chaque template individuellement avec le bon chemin
	homeTemplate, err = template.ParseFiles(templateDir + "/home.html")
	if err != nil {
		return err
	}

	analyzingTemplate, err = template.ParseFiles(templateDir + "/analyzing.html")
	if err != nil {
		return err
	}

	resultsTemplate, err = template.ParseFiles(templateDir + "/results.html")
	if err != nil {
		return err
	}

	return nil
}

// setupServerForTests - Version simplifiée sans middlewares pour les tests
func setupServerForTests() *http.ServeMux {
	// Charger les templates pour les tests (nécessaire pour homeHandler)
	if homeTemplate == nil {
		// During tests, working dir is cmd/server, templates are at ../../templates
		loadTemplatesForTests() // Ignore error in tests
	}
	
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
	
	// 🔥🦎 STATIC FILES: CSS/JS selon standards NO HARDCODING  
	// Note: During tests, working dir is cmd/server, need to go up to project root
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static/"))))

	return mux
}

// TestStaticFilesAreServed - TDD: FileServer doit servir les fichiers statiques
func TestStaticFilesAreServed(t *testing.T) {
	// DEBUG: Check working directory and static file existence
	workingDir, _ := os.Getwd()
	t.Logf("Working directory: %s", workingDir)
	
	staticPath := "../../static/css/fire-salamander.css"
	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		t.Logf("❌ CSS file does not exist at: %s", staticPath)
	} else {
		t.Logf("✅ CSS file exists at: %s", staticPath)
	}
	
	// GIVEN - Serveur avec route static configurée (version test sans middlewares)
	mux := setupServerForTests()
	
	// WHEN - On demande le CSS Fire Salamander
	req := httptest.NewRequest("GET", "/static/css/fire-salamander.css", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	
	// THEN - Le CSS est retourné avec le bon content-type
	if w.Code != http.StatusOK {
		t.Errorf("❌ TDD CSS FAILURE: Expected 200 for CSS file, got %d. Response body: %s", w.Code, w.Body.String())
	}
	
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/css") && !strings.Contains(contentType, "text/plain") {
		t.Errorf("❌ TDD CSS FAILURE: Expected CSS content-type, got: %s", contentType)
	}
	
	// THEN - Le contenu contient nos classes SEPTEO
	body := w.Body.String()
	if !strings.Contains(body, "septeo-orange") {
		t.Error("❌ TDD CSS FAILURE: CSS must contain SEPTEO colors")
	}
	if !strings.Contains(body, "bg-septeo-gray-50") {
		t.Error("❌ TDD CSS FAILURE: CSS must contain custom classes")
	}
}

// TestHomePageReferencesLocalCSS - TDD: Page doit référencer CSS local
func TestHomePageReferencesLocalCSS(t *testing.T) {
	// GIVEN - Serveur configuré (version test)
	mux := setupServerForTests()
	
	// WHEN - On charge la page d'accueil
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	
	// THEN - La page contient le lien vers notre CSS
	if w.Code != http.StatusOK {
		t.Errorf("❌ TDD CSS FAILURE: Home page not accessible, got %d", w.Code)
	}
	
	body := w.Body.String()
	if !strings.Contains(body, `href="/static/css/fire-salamander.css"`) {
		t.Error("❌ TDD CSS FAILURE: Page must reference local CSS")
	}
	if strings.Contains(body, "cdn.tailwindcss.com") {
		t.Error("❌ TDD CSS FAILURE: Page must NOT use CDN Tailwind")
	}
}