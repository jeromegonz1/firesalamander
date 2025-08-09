package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"firesalamander/internal/constants"
	"firesalamander/internal/integration"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 🔥🦎 FIRE SALAMANDER SPRINT 5 - TDD INTEGRATION TESTS
// Ces tests DOIVENT échouer jusqu'à l'implémentation du RealOrchestrator

// TestRealOrchestrator_StartAnalysis teste le démarrage d'une analyse réelle
func TestRealOrchestrator_StartAnalysis(t *testing.T) {
	// Créer serveur de test avec contenu HTML réel
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
	<title>Fire Salamander Test Page - SEO Analysis Tool</title>
	<meta name="description" content="Test page for Fire Salamander real SEO analysis with comprehensive scoring">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
	<h1>Fire Salamander SEO Analysis</h1>
	<h2>Features</h2>
	<h3>Real-time Analysis</h3>
	<p>Our tool provides comprehensive SEO analysis for your website.</p>
	<img src="logo.png" alt="Fire Salamander Logo">
	<img src="feature.jpg" alt="Real-time SEO analysis feature">
	<a href="/page2">Internal Link</a>
	<a href="https://external.com">External Link</a>
</body>
</html>`))
	}))
	defer testServer.Close()

	// Test RealOrchestrator creation and analysis start
	orchestrator := integration.NewRealOrchestrator()
	require.NotNil(t, orchestrator)

	// Start analysis - should return unique ID
	analysisID, err := orchestrator.StartAnalysis(testServer.URL)
	require.NoError(t, err)
	require.NotEmpty(t, analysisID)

	// Verify initial state
	state, err := orchestrator.GetStatus(analysisID)
	require.NoError(t, err)
	assert.Equal(t, analysisID, state.ID)
	assert.Equal(t, testServer.URL, state.URL)
	assert.Contains(t, []string{constants.OrchestratorStatusStarting, constants.OrchestratorStatusCrawling}, state.Status)
	assert.False(t, state.StartTime.IsZero())
}

// TestRealOrchestrator_CompleteAnalysisFlow teste le flow complet d'analyse
func TestRealOrchestrator_CompleteAnalysisFlow(t *testing.T) {
	// Créer serveur multi-pages pour test réaliste
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		
		switch r.URL.Path {
		case "/":
			w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<title>Home - Fire Salamander</title>
	<meta name="description" content="Fire Salamander homepage with SEO optimization">
</head>
<body>
	<h1>Welcome to Fire Salamander</h1>
	<img src="logo.png" alt="Fire Salamander Logo">
	<a href="/about">About Us</a>
	<a href="/contact">Contact</a>
</body>
</html>`))
		case "/about":
			w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<title>About - Fire Salamander SEO Tool</title>
	<meta name="description" content="Learn about Fire Salamander SEO analysis capabilities">
</head>
<body>
	<h1>About Fire Salamander</h1>
	<h2>Our Mission</h2>
	<p>Professional SEO analysis for everyone.</p>
	<img src="team.jpg" alt="Fire Salamander team">
</body>
</html>`))
		case "/contact":
			w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<meta name="description" content="Contact Fire Salamander for SEO help">
</head>
<body>
	<h2>Contact Us</h2>
	<p>Get in touch for SEO analysis.</p>
	<img src="contact.png">
</body>
</html>`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer testServer.Close()

	orchestrator := integration.NewRealOrchestrator()
	analysisID, err := orchestrator.StartAnalysis(testServer.URL)
	require.NoError(t, err)

	// Attendre la completion (max 30 secondes)
	ctx, cancel := context.WithTimeout(context.Background(), constants.MaxAnalysisWaitTime)
	defer cancel()

	ticker := time.NewTicker(constants.AnalysisPollingInterval)
	defer ticker.Stop()

	var finalState *integration.AnalysisState

	for {
		select {
		case <-ctx.Done():
			t.Fatal("Analysis timeout - should complete within 30 seconds")
		case <-ticker.C:
			state, err := orchestrator.GetStatus(analysisID)
			require.NoError(t, err)

			// Log progress pour debugging
			t.Logf("Status: %s, Pages Found: %d, Pages Analyzed: %d", 
				state.Status, state.PagesFound, state.PagesAnalyzed)

			if state.Status == constants.OrchestratorStatusComplete {
				finalState = state
				goto completed
			}
			
			if state.Status == constants.OrchestratorStatusError {
				t.Fatalf("Analysis failed with error: %v", state.Error)
			}
		}
	}

completed:
	// Vérifications des résultats RÉELS
	require.NotNil(t, finalState)
	
	// Vérifier que des pages ont été crawlées et analysées
	assert.Greater(t, finalState.PagesFound, 0, "Should find at least 1 page")
	assert.Greater(t, finalState.PagesAnalyzed, 0, "Should analyze at least 1 page")
	assert.Equal(t, finalState.PagesFound, finalState.PagesAnalyzed, "All found pages should be analyzed")
	
	// Vérifier les résultats SEO réels
	assert.NotEmpty(t, finalState.Pages, "Should have page analyses")
	assert.Greater(t, finalState.GlobalScore, 0, "Should have a real global score > 0")
	assert.NotEmpty(t, finalState.GlobalGrade, "Should have a grade")
	assert.Contains(t, []string{
		constants.SEOGradeAPlus, constants.SEOGradeA, constants.SEOGradeBPlus, 
		constants.SEOGradeB, constants.SEOGradeC, constants.SEOGradeD, constants.SEOGradeF,
	}, finalState.GlobalGrade, "Grade should be valid")
	
	// Vérifier les recommandations RÉELLES
	assert.NotEmpty(t, finalState.Recommendations, "Should have real recommendations")
	
	// Vérifier qu'au moins une page a un titre (page d'accueil et about ont des titres)
	hasPageWithTitle := false
	for _, page := range finalState.Pages {
		if page.Title.Present {
			hasPageWithTitle = true
			break
		}
	}
	assert.True(t, hasPageWithTitle, "At least one page should have a title")
	
	// Vérifier les métriques de temps
	assert.True(t, finalState.Duration > 0, "Analysis should take measurable time")
	assert.True(t, finalState.Duration < constants.MaxAnalysisWaitTime, "Analysis should complete within timeout")
}

// TestRealOrchestrator_RealTimeUpdates teste les mises à jour temps réel
func TestRealOrchestrator_RealTimeUpdates(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slower server for testing updates
		time.Sleep(constants.TestServerDelay)
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Test Page</title></head>
<body>
	<h1>Test Content</h1>
	<a href="/page2">Link 2</a>
	<a href="/page3">Link 3</a>
</body>
</html>`))
	}))
	defer testServer.Close()

	orchestrator := integration.NewRealOrchestrator()
	analysisID, err := orchestrator.StartAnalysis(testServer.URL)
	require.NoError(t, err)

	// Vérifier que les métriques sont mises à jour en temps réel
	var lastPagesFound, lastPagesAnalyzed int
	updateCount := 0
	
	for i := 0; i < constants.MaxUpdateChecks; i++ {
		time.Sleep(constants.UpdateCheckInterval)
		
		state, err := orchestrator.GetStatus(analysisID)
		require.NoError(t, err)
		
		// Vérifier que les métriques progressent
		if state.PagesFound > lastPagesFound || state.PagesAnalyzed > lastPagesAnalyzed {
			updateCount++
			lastPagesFound = state.PagesFound
			lastPagesAnalyzed = state.PagesAnalyzed
			t.Logf("Update %d: Found=%d, Analyzed=%d, Workers=%d, PPS=%.2f", 
				updateCount, state.PagesFound, state.PagesAnalyzed, state.CurrentWorkers, state.PagesPerSecond)
		}
		
		if state.Status == constants.OrchestratorStatusComplete {
			break
		}
	}
	
	// Vérifier qu'on a eu des mises à jour
	assert.Greater(t, updateCount, 0, "Should have real-time updates during analysis")
}

// TestRealOrchestrator_ErrorHandling teste la gestion d'erreur
func TestRealOrchestrator_ErrorHandling(t *testing.T) {
	orchestrator := integration.NewRealOrchestrator()

	// Test avec URL invalide
	_, err := orchestrator.StartAnalysis("invalid-url")
	assert.Error(t, err)
	
	// Test avec URL inexistante
	analysisID, err := orchestrator.StartAnalysis("https://nonexistent.domain.fake")
	require.NoError(t, err) // StartAnalysis should succeed, but analysis should fail
	
	// Attendre que l'analyse échoue
	ctx, cancel := context.WithTimeout(context.Background(), constants.ErrorAnalysisTimeout)
	defer cancel()
	
	ticker := time.NewTicker(constants.ErrorCheckInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			t.Fatal("Expected analysis to fail within timeout")
		case <-ticker.C:
			state, err := orchestrator.GetStatus(analysisID)
			require.NoError(t, err)
			
			if state.Status == constants.OrchestratorStatusError {
				assert.NotEmpty(t, state.Error, "Error state should have error message")
				return
			}
		}
	}
}

// TestRealOrchestrator_NoHardcoding vérifie zéro hardcoding
func TestRealOrchestrator_NoHardcoding(t *testing.T) {
	orchestrator := integration.NewRealOrchestrator()
	require.NotNil(t, orchestrator)
	
	// Vérifier que la configuration utilise des constantes
	config := orchestrator.GetConfig()
	require.NotNil(t, config)
	
	// Tous les paramètres doivent venir des constantes
	assert.Equal(t, constants.RealOrchestratorMaxPages, config.MaxPages)
	assert.Equal(t, constants.RealOrchestratorMaxWorkers, config.MaxWorkers)
	assert.Equal(t, constants.RealOrchestratorInitialWorkers, config.InitialWorkers)
	assert.Equal(t, time.Duration(constants.RealOrchestratorAnalysisTimeout), config.Timeout)
	
	// Les seuils de scoring doivent utiliser des constantes
	assert.True(t, constants.MaxSEOScore > 0)
	assert.True(t, constants.GradeAThreshold > 0)
	assert.True(t, constants.GradeBThreshold > 0)
}

// TestRealOrchestrator_PerformanceMetrics teste les métriques de performance
func TestRealOrchestrator_PerformanceMetrics(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Performance Test</title></head>
<body><h1>Test</h1></body>
</html>`))
	}))
	defer testServer.Close()

	orchestrator := integration.NewRealOrchestrator()
	analysisID, err := orchestrator.StartAnalysis(testServer.URL)
	require.NoError(t, err)

	// Attendre la completion
	ctx, cancel := context.WithTimeout(context.Background(), constants.MaxAnalysisWaitTime)
	defer cancel()

	ticker := time.NewTicker(constants.AnalysisPollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Fatal("Analysis timeout")
		case <-ticker.C:
			state, err := orchestrator.GetStatus(analysisID)
			require.NoError(t, err)

			if state.Status == constants.OrchestratorStatusComplete {
				// Vérifier les métriques de performance
				assert.True(t, state.Duration < constants.MaxAnalysisWaitTime, "Analysis should be reasonably fast")
				assert.Greater(t, state.PagesPerSecond, float64(0), "Should have measurable pages/second rate")
				assert.Greater(t, state.CurrentWorkers, 0, "Should have used workers")
				
				// Vérifier que l'analyse est < 2 minutes comme spécifié
				assert.True(t, state.Duration < constants.MaxAcceptableAnalysisTime, 
					fmt.Sprintf("Analysis took %v, should be < %v", state.Duration, constants.MaxAcceptableAnalysisTime))
				return
			}
		}
	}
}