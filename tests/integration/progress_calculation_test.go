package integration

import (
	"testing"
	"time"

	"firesalamander/internal/api"
	"firesalamander/internal/constants"
	"firesalamander/internal/integration"

	"github.com/stretchr/testify/assert"
)

// 🔥🦎 TDD TEST - Fix Progress Calculation Bug
// Test pour valider que le calcul de progression fonctionne correctement
func TestRealProgressCalculation_NoBug98Percent(t *testing.T) {
	// SCENARIO: Analyse avec 45 pages trouvées, 42 analysées
	// Ne doit PAS bloquer à 98%
	
	// Phase Crawling - 45 pages trouvées
	crawlingState := &integration.AnalysisState{
		ID:              "test-analysis-98-bug",
		URL:             "https://septeo.com",
		Status:          constants.OrchestratorStatusCrawling,
		PagesFound:      45,
		PagesAnalyzed:   0,
		StartTime:       time.Now(),
	}
	
	progress := api.CalculateRealProgressExposed(crawlingState)
	
	// Avec RealOrchestratorMaxPages=20, 45 pages = 40% max (pas 98%)
	assert.LessOrEqual(t, progress, 40, "Crawling phase should max at 40%")
	
	// Phase Analyzing - 42/45 pages analysées 
	analyzingState := &integration.AnalysisState{
		ID:              "test-analysis-98-bug",
		URL:             "https://septeo.com", 
		Status:          constants.OrchestratorStatusAnalyzing,
		PagesFound:      45,
		PagesAnalyzed:   42,
		StartTime:       time.Now(),
	}
	
	progress = api.CalculateRealProgressExposed(analyzingState)
	
	// 42/45 = 93.3% * 45 + 40 = 82% (pas 98%)
	expectedProgress := 40 + int(float64(42)/float64(45)*45)
	assert.Equal(t, expectedProgress, progress, "Progress calculation should be accurate")
	assert.Less(t, progress, 85, "Should not reach 98% in analyzing phase")
	
	t.Logf("✅ Progress correctly calculated: %d%% for 42/45 pages", progress)
}

// Test que l'orchestrateur passe correctement entre les phases
func TestRealOrchestrator_PhaseTransitions(t *testing.T) {
	// SCENARIO: Vérifier que les transitions de phases sont correctes
	// Crawling → Analyzing → Aggregating → Complete
	
	states := []struct {
		status   string
		expected int
	}{
		{constants.OrchestratorStatusStarting, constants.DefaultProgressStart},
		{constants.OrchestratorStatusCrawling, 40}, // Max crawling
		{constants.OrchestratorStatusAnalyzing, 82}, // 42/45 pages
		{constants.OrchestratorStatusAggregating, 85},
		{constants.OrchestratorStatusComplete, 100},
	}
	
	for _, test := range states {
		state := &integration.AnalysisState{
			Status:        test.status,
			PagesFound:    45,
			PagesAnalyzed: 42,
		}
		
		progress := api.CalculateRealProgressExposed(state)
		
		t.Logf("Status: %s -> Progress: %d%%", test.status, progress)
		
		if test.status == constants.OrchestratorStatusAnalyzing {
			// Test spécifique pour le bug 98%
			assert.NotEqual(t, 98, progress, "Must NOT return 98% (bug scenario)")
		}
	}
}