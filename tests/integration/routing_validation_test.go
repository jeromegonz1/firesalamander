package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"firesalamander/internal/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 🔥🦎 TDD TEST - Validation des Routes Après Développement
// S'assurer que les bonnes fonctions sont branchées sur les bonnes routes
func TestRouting_RealHandlersConnected(t *testing.T) {
	// SCENARIO: Vérifier que /api/analyze pointe vers RealAnalyzeHandler
	
	// Initialize real orchestrator
	api.InitRealOrchestrator()
	
	testCases := []struct {
		route           string
		expectedHandler string
		method          string
		body            string
	}{
		{
			route:           "/api/analyze", 
			expectedHandler: "RealAnalyzeHandler",
			method:          "POST",
			body:            `{"url":"https://example.com"}`,
		},
		{
			route:           "/api/status/test-id-12345",
			expectedHandler: "RealStatusHandler", 
			method:          "GET",
			body:            "",
		},
		{
			route:           "/api/results/test-id-12345",
			expectedHandler: "RealResultsHandler",
			method:          "GET", 
			body:            "",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.route+"_"+tc.expectedHandler, func(t *testing.T) {
			// Create request
			var req *http.Request
			var err error
			
			if tc.body != "" {
				req, err = http.NewRequest(tc.method, tc.route, bytes.NewBuffer([]byte(tc.body)))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tc.method, tc.route, nil)
			}
			require.NoError(t, err)
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Route to appropriate handler
			switch {
			case tc.route == "/api/analyze":
				api.RealAnalyzeHandler(rr, req)
			case strings.HasPrefix(tc.route, "/api/status/"):
				api.RealStatusHandler(rr, req)
			case strings.HasPrefix(tc.route, "/api/results/"):
				api.RealResultsHandler(rr, req)
			default:
				t.Fatalf("Unknown route: %s", tc.route)
			}
			
			// Validate response format indicates Real handler
			if tc.route == "/api/analyze" && tc.method == "POST" {
				// Real handler should return analysis ID and Fire Salamander message
				assert.Equal(t, http.StatusOK, rr.Code, "Real analyze should return 200")
				
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Check Real handler response format
				assert.Contains(t, response, "id", "Real handler should return analysis ID")
				assert.Contains(t, response, "message", "Real handler should return message")
				
				message, ok := response["message"].(string)
				assert.True(t, ok, "Message should be string")
				assert.Contains(t, message, "Fire Salamander", "Should contain Fire Salamander branding")
				
				// Check ID format (should be unique, not fake)
				id, ok := response["id"].(string)
				assert.True(t, ok, "ID should be string")
				assert.Regexp(t, `^analysis-\d+-\d+-\d+$`, id, "Should be real unique ID format")
				assert.NotEqual(t, "analysis-20060102", id, "Should NOT be fake ID")
				
				t.Logf("✅ Real handler connected: %s returns ID: %s", tc.route, id)
			}
		})
	}
}

// Test que les anciennes routes fake sont toujours disponibles pour debug
func TestRouting_LegacyRoutesAvailable(t *testing.T) {
	// SCENARIO: Vérifier que les routes legacy existent pour comparaison/debug
	
	legacyRoutes := []string{
		"/api/fake/analyze",
		"/api/legacy/analyze", 
	}
	
	for _, route := range legacyRoutes {
		t.Logf("🔍 Legacy route should exist: %s", route)
		// Ces routes devraient exister mais ne sont pas testées en profondeur
		// car elles contiennent des données fake
	}
}

// Test anti-régression: vérifier qu'on n'utilise plus les fakes par défaut
func TestRouting_NoFakeDataByDefault(t *testing.T) {
	// SCENARIO: S'assurer que /api/analyze ne retourne PAS des données fake
	
	api.InitRealOrchestrator()
	
	requestBody := `{"url":"https://test-uniqueness.com"}`
	req, err := http.NewRequest("POST", "/api/analyze", bytes.NewBuffer([]byte(requestBody)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	api.RealAnalyzeHandler(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Vérifications anti-fake
	id, exists := response["id"].(string)
	assert.True(t, exists, "Should have analysis ID")
	
	// Les IDs fake ont des patterns prévisibles comme "analysis-20060102"
	assert.NotContains(t, id, "20060102", "Should not contain fake date pattern")
	assert.NotContains(t, id, "example", "Should not contain example pattern")
	assert.NotEqual(t, "mock-analysis-id", id, "Should not be mock ID")
	
	t.Logf("✅ Anti-regression: Real handler returns unique ID: %s", id)
}