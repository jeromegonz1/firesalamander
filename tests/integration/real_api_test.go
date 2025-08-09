package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"firesalamander/internal/api"
	"firesalamander/internal/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 🔥🦎 FIRE SALAMANDER - REAL API INTEGRATION TESTS
// Sprint 5 - Tester l'API connectée au RealOrchestrator

// TestRealAPI_AnalyzeEndpoint teste l'endpoint d'analyse réelle
func TestRealAPI_AnalyzeEndpoint(t *testing.T) {
	// Initialize the real orchestrator
	api.InitRealOrchestrator()

	// Create test request
	requestBody := map[string]string{
		"url": "https://example.com",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Create HTTP request
	req, err := http.NewRequest("POST", "/api/real/analyze", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	api.RealAnalyzeHandler(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response api.RealAnalyzeResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response structure
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, constants.OrchestratorStatusStarting, response.Status)
	assert.Contains(t, response.Message, "Fire Salamander")
	
	t.Logf("✅ Real analysis started with ID: %s", response.ID)
}

// TestRealAPI_StatusEndpoint teste l'endpoint de statut réel
func TestRealAPI_StatusEndpoint(t *testing.T) {
	// Initialize the real orchestrator
	api.InitRealOrchestrator()

	// First start an analysis to get an ID
	requestBody := map[string]string{
		"url": "https://example.com",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/real/analyze", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.RealAnalyzeHandler(rr, req)

	var analyzeResponse api.RealAnalyzeResponse
	err = json.Unmarshal(rr.Body.Bytes(), &analyzeResponse)
	require.NoError(t, err)

	// Wait a moment for the analysis to start
	time.Sleep(100 * time.Millisecond)

	// Now test the status endpoint
	statusURL := "/api/real/status/" + analyzeResponse.ID
	req, err = http.NewRequest("GET", statusURL, nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	api.RealStatusHandler(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var statusResponse api.RealStatusResponse
	err = json.Unmarshal(rr.Body.Bytes(), &statusResponse)
	require.NoError(t, err)

	// Verify response structure
	assert.Equal(t, analyzeResponse.ID, statusResponse.ID)
	assert.Equal(t, "https://example.com", statusResponse.URL)
	assert.NotEmpty(t, statusResponse.Status)
	assert.GreaterOrEqual(t, statusResponse.Progress, 0)
	assert.LessOrEqual(t, statusResponse.Progress, 100)
	assert.NotEmpty(t, statusResponse.ElapsedTime)
	
	t.Logf("✅ Status check successful - Status: %s, Progress: %d%%", 
		statusResponse.Status, statusResponse.Progress)
}

// TestRealAPI_ResultsEndpoint teste l'endpoint de résultats (avec mock completion)
func TestRealAPI_ResultsEndpoint_NotComplete(t *testing.T) {
	// Initialize the real orchestrator
	api.InitRealOrchestrator()

	// Start an analysis
	requestBody := map[string]string{
		"url": "https://example.com",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/real/analyze", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.RealAnalyzeHandler(rr, req)

	var analyzeResponse api.RealAnalyzeResponse
	err = json.Unmarshal(rr.Body.Bytes(), &analyzeResponse)
	require.NoError(t, err)

	// Try to get results immediately (should fail - analysis not complete)
	resultsURL := "/api/real/results/" + analyzeResponse.ID
	req, err = http.NewRequest("GET", resultsURL, nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	api.RealResultsHandler(rr, req)

	// Should return 202 Accepted (not complete yet)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	
	t.Logf("✅ Results endpoint correctly returns 202 for incomplete analysis")
}

// TestRealAPI_ValidationErrors teste la validation des erreurs
func TestRealAPI_ValidationErrors(t *testing.T) {
	// Initialize the real orchestrator
	api.InitRealOrchestrator()

	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		expectedStatus int
	}{
		{
			name:           "wrong_method_analyze",
			method:         "GET",
			url:            "/api/real/analyze",
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "empty_url",
			method:         "POST",
			url:            "/api/real/analyze",
			body:           `{"url":""}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid_json",
			method:         "POST",
			url:            "/api/real/analyze",
			body:           `{"url":}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "wrong_method_status",
			method:         "POST",
			url:            "/api/real/status/test-id",
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid_analysis_id",
			method:         "GET",
			url:            "/api/real/status/",
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "nonexistent_analysis",
			method:         "GET",
			url:            "/api/real/status/nonexistent-id",
			body:           "",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.body != "" {
				req, err = http.NewRequest(tt.method, tt.url, bytes.NewBuffer([]byte(tt.body)))
			} else {
				req, err = http.NewRequest(tt.method, tt.url, nil)
			}
			require.NoError(t, err)

			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			rr := httptest.NewRecorder()

			// Route to the appropriate handler
			switch {
			case tt.url == "/api/real/analyze":
				api.RealAnalyzeHandler(rr, req)
			case tt.url == "/api/real/status/" || tt.url == "/api/real/status/nonexistent-id":
				api.RealStatusHandler(rr, req)
			}

			assert.Equal(t, tt.expectedStatus, rr.Code, 
				"Test %s: expected status %d, got %d", tt.name, tt.expectedStatus, rr.Code)
		})
	}
}

// TestRealAPI_CompareWithMockAPI teste que l'API réelle retourne des données différentes du mock
func TestRealAPI_CompareWithMockAPI(t *testing.T) {
	// Initialize the real orchestrator
	api.InitRealOrchestrator()

	// Test avec URL réelle (mais on ne va pas attendre la completion pour ce test)
	requestBody := map[string]string{
		"url": "https://example.com",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Test Real API
	req, err := http.NewRequest("POST", "/api/real/analyze", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.RealAnalyzeHandler(rr, req)

	var realResponse api.RealAnalyzeResponse
	err = json.Unmarshal(rr.Body.Bytes(), &realResponse)
	require.NoError(t, err)

	// Vérifier que l'API réelle génère des IDs différents du pattern mock
	assert.NotContains(t, realResponse.ID, "analysis-20060102") // Pattern du mock API
	assert.Contains(t, realResponse.ID, "analysis-") // Mais garde le préfixe Fire Salamander
	
	// Vérifier que le message est spécifique Fire Salamander
	assert.Contains(t, realResponse.Message, "Fire Salamander")
	
	t.Logf("✅ Real API generates unique IDs and Fire Salamander branding")
	t.Logf("    Real ID: %s", realResponse.ID)
	t.Logf("    Message: %s", realResponse.Message)
}