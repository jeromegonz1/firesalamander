package integration

import (
	"testing"
	"time"

	"firesalamander/internal/integration"

	"github.com/stretchr/testify/assert"
)

// 🔥🦎 TDD TEST - Multi-Audits avec IDs Uniques
// Vérifier que plusieurs audits du même site ont des IDs différents
func TestOrchestrator_MultipleAuditsUniqueIDs(t *testing.T) {
	orchestrator := integration.NewOrchestrator()
	
	// SCENARIO: 5 audits simultanés du même site
	sameURL := "https://septeo.com"
	var analysisIDs []string
	
	// Lancer plusieurs analyses du même site rapidement
	for i := 0; i < 5; i++ {
		analysisID, err := orchestrator.StartAnalysis(sameURL)
		assert.NoError(t, err, "Should start analysis successfully")
		assert.NotEmpty(t, analysisID, "Analysis ID should not be empty")
		
		analysisIDs = append(analysisIDs, analysisID)
		
		// Petit délai pour éviter collision nanoseconde
		time.Sleep(1 * time.Millisecond)
		
		t.Logf("✅ Analysis %d started with ID: %s", i+1, analysisID)
	}
	
	// VALIDATION: Tous les IDs doivent être uniques
	uniqueIDs := make(map[string]bool)
	for _, id := range analysisIDs {
		assert.False(t, uniqueIDs[id], "ID %s should be unique, but found duplicate", id)
		uniqueIDs[id] = true
		
		// Vérifier le nouveau format attendu (timestamp-nanoseconde-pid)
		assert.Regexp(t, `^analysis-\d+-\d+-\d+$`, id, "ID should match expected format")
	}
	
	assert.Equal(t, 5, len(uniqueIDs), "Should have 5 unique analysis IDs")
	
	// VALIDATION: Chaque analyse doit avoir son propre état
	for _, id := range analysisIDs {
		state, err := orchestrator.GetStatus(id)
		assert.NoError(t, err, "Should get status for ID %s", id)
		assert.Equal(t, id, state.ID, "State ID should match requested ID")
		assert.Equal(t, sameURL, state.URL, "All analyses should be for same URL")
	}
	
	t.Logf("🎯 SUCCESS: %d simultaneous audits with unique IDs", len(analysisIDs))
}

// Test concurrent access to analysis states
func TestOrchestrator_ConcurrentAccess(t *testing.T) {
	orchestrator := integration.NewOrchestrator()
	
	// Démarrer une analyse
	analysisID, err := orchestrator.StartAnalysis("https://example.com")
	assert.NoError(t, err)
	
	// Test accès concurrent au même état
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(workerID int) {
			defer func() { done <- true }()
			
			// Chaque goroutine accède au même état
			state, err := orchestrator.GetStatus(analysisID)
			assert.NoError(t, err, "Worker %d should get status", workerID)
			assert.Equal(t, analysisID, state.ID, "Worker %d should get correct state", workerID)
		}(i)
	}
	
	// Attendre tous les workers
	for i := 0; i < 10; i++ {
		<-done
	}
	
	t.Logf("✅ Concurrent access test passed")
}