package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"firesalamander/internal/constants"
	"firesalamander/internal/integration"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test de debug pour comprendre pourquoi l'orchestrateur ne fonctionne pas
func TestOrchestrator_DebugCrawling(t *testing.T) {
	// Serveur de test ultra simple
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Server received request: %s %s", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Debug Test</title></head>
<body><h1>Debug Test Page</h1><p>Simple content.</p></body>
</html>`))
	}))
	defer testServer.Close()

	t.Logf("Test server URL: %s", testServer.URL)

	orchestrator := integration.NewOrchestrator()
	require.NotNil(t, orchestrator)

	// Start analysis
	analysisID, err := orchestrator.StartAnalysis(testServer.URL)
	require.NoError(t, err)
	t.Logf("Started analysis with ID: %s", analysisID)

	// Poll with debug info AND listen to updates
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second) // Poll every second
	defer ticker.Stop()

	// Listen to updates in a separate goroutine
	go func() {
		updatesChan := orchestrator.GetUpdatesChannel()
		for update := range updatesChan {
			if update.AnalysisID == analysisID {
				t.Logf("UPDATE: %s - %s", update.Status, update.Message)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			t.Fatal("Analysis timeout after 60 seconds")
		case <-ticker.C:
			state, err := orchestrator.GetStatus(analysisID)
			require.NoError(t, err)

			t.Logf("DEBUG: Status=%s, PagesFound=%d, PagesAnalyzed=%d, Workers=%d, PPS=%.2f",
				state.Status, state.PagesFound, state.PagesAnalyzed, state.CurrentWorkers, state.PagesPerSecond)

			if state.Status == constants.OrchestratorStatusComplete {
				t.Logf("Analysis completed successfully!")
				t.Logf("Final Score: %d, Grade: %s", state.GlobalScore, state.GlobalGrade)
				t.Logf("Pages Analyzed: %d", len(state.Pages))
				return
			}

			if state.Status == constants.OrchestratorStatusError {
				t.Fatalf("Analysis failed with error: %s", state.Error)
			}
		}
	}
}

// Test le crawler directement pour vérifier qu'il fonctionne
func TestDirectCrawlerUsage(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Direct Crawler Test</title></head>
<body><h1>Test Page</h1></body>
</html>`))
	}))
	defer testServer.Close()

	orchestrator := integration.NewOrchestrator()
	
	// Test le crawling directement
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Logf("Testing direct crawling of: %s", testServer.URL)
	
	// Use the test method
	crawlResult, err := orchestrator.TestCrawl(ctx, testServer.URL)
	require.NoError(t, err)
	require.NotNil(t, crawlResult)
	
	t.Logf("Direct crawl success! Pages found: %d", len(crawlResult.Pages))
	assert.Greater(t, len(crawlResult.Pages), 0)
}