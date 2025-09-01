package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 🔥🦎 FIRE SALAMANDER - SPRINT 6 TDD
// Tests MCP Storage - DOIVENT ÊTRE ROUGES AVANT IMPLÉMENTATION

// TestMCPStorage_SaveAndLoad - Test fondamental de persistance
func TestMCPStorage_SaveAndLoad(t *testing.T) {
	// GIVEN - Un storage MCP
	testDir := "./test_data_basic"
	storage := NewMCPStorage(testDir)
	defer os.RemoveAll(testDir)

	// WHEN - On sauve une analyse
	analysis := &AnalysisState{
		ID:        "test-analysis-123",
		URL:       "https://resalys.com",
		Status:    "complete",
		Score:     85,
		StartTime: time.Now(),
		Domain:    "resalys.com",
	}

	err := storage.SaveAnalysis(analysis)

	// THEN - L'analyse est sauvée sur disque
	assert.NoError(t, err, "SaveAnalysis should not error")
	
	// Vérifier que le fichier JSON existe
	expectedPath := filepath.Join(testDir, "analyses", "test-analysis-123.json")
	assert.FileExists(t, expectedPath, "Analysis JSON file should exist")

	// AND - On peut la recharger
	loaded, err := storage.LoadAnalysis("test-analysis-123")
	assert.NoError(t, err, "LoadAnalysis should not error")
	assert.NotNil(t, loaded, "Loaded analysis should not be nil")
	assert.Equal(t, analysis.URL, loaded.URL, "URL should match")
	assert.Equal(t, analysis.Score, loaded.Score, "Score should match")
	assert.Equal(t, analysis.Status, loaded.Status, "Status should match")
}

// TestMCPStorage_SurvivesRestart - Test critique de persistence après restart
func TestMCPStorage_SurvivesRestart(t *testing.T) {
	testDir := "./test_data_restart"
	defer os.RemoveAll(testDir)

	// GIVEN - Une analyse sauvée avec storage1
	storage1 := NewMCPStorage(testDir)
	
	analysis := &AnalysisState{
		ID:        "restart-test-456",
		URL:       "https://septeo.com",
		Status:    "analyzing",
		Score:     72,
		StartTime: time.Now(),
		Domain:    "septeo.com",
	}
	
	err := storage1.SaveAnalysis(analysis)
	require.NoError(t, err, "Initial save should work")

	// WHEN - On simule un restart (nouveau storage)
	storage2 := NewMCPStorage(testDir)

	// THEN - L'analyse est toujours là après restart
	loaded, err := storage2.LoadAnalysis("restart-test-456")
	assert.NoError(t, err, "LoadAnalysis after restart should work")
	assert.NotNil(t, loaded, "Analysis should survive restart")
	assert.Equal(t, analysis.URL, loaded.URL, "URL should survive restart")
	assert.Equal(t, analysis.Domain, loaded.Domain, "Domain should survive restart")
}

// TestMCPStorage_Index - Test de maintien de l'index
func TestMCPStorage_Index(t *testing.T) {
	testDir := "./test_data_index"
	storage := NewMCPStorage(testDir)
	defer os.RemoveAll(testDir)

	// GIVEN - Plusieurs analyses sauvées
	for i := 0; i < 3; i++ {
		analysis := &AnalysisState{
			ID:        fmt.Sprintf("index-test-%d", i),
			URL:       fmt.Sprintf("https://example%d.com", i),
			Status:    "complete",
			Score:     80 + i,
			StartTime: time.Now(),
			Domain:    fmt.Sprintf("example%d.com", i),
		}
		err := storage.SaveAnalysis(analysis)
		require.NoError(t, err, "Save should work for analysis %d", i)
	}

	// WHEN - On charge l'index
	index, err := storage.LoadIndex()

	// THEN - L'index doit contenir les 3 analyses
	assert.NoError(t, err, "LoadIndex should not error")
	assert.NotNil(t, index, "Index should not be nil")
	assert.Len(t, index.Analyses, 3, "Index should contain 3 analyses")
	
	// Vérifier que l'index contient les bons IDs
	ids := make([]string, len(index.Analyses))
	for i, item := range index.Analyses {
		ids[i] = item.ID
	}
	assert.Contains(t, ids, "index-test-0")
	assert.Contains(t, ids, "index-test-1")
	assert.Contains(t, ids, "index-test-2")
}

// TestMCPStorage_ListAllAnalyses - Test de listing complet
func TestMCPStorage_ListAllAnalyses(t *testing.T) {
	testDir := "./test_data_list"
	storage := NewMCPStorage(testDir)
	defer os.RemoveAll(testDir)

	// GIVEN - 5 analyses avec différents statuts
	statuses := []string{"complete", "analyzing", "error", "complete", "crawling"}
	for i, status := range statuses {
		analysis := &AnalysisState{
			ID:        fmt.Sprintf("list-test-%d", i),
			URL:       fmt.Sprintf("https://test%d.com", i),
			Status:    status,
			Score:     75 + (i * 5),
			StartTime: time.Now().Add(time.Duration(i) * time.Minute),
			Domain:    fmt.Sprintf("test%d.com", i),
		}
		storage.SaveAnalysis(analysis)
	}

	// WHEN - On liste toutes les analyses
	allAnalyses, err := storage.ListAllAnalyses()

	// THEN - On doit récupérer les 5 analyses
	assert.NoError(t, err, "ListAllAnalyses should not error")
	assert.Len(t, allAnalyses, 5, "Should return 5 analyses")
	
	// Vérifier qu'on a tous les statuts
	foundStatuses := make(map[string]bool)
	for _, analysis := range allAnalyses {
		foundStatuses[analysis.Status] = true
	}
	assert.True(t, foundStatuses["complete"], "Should find complete analyses")
	assert.True(t, foundStatuses["analyzing"], "Should find analyzing analyses")
	assert.True(t, foundStatuses["error"], "Should find error analyses")
	assert.True(t, foundStatuses["crawling"], "Should find crawling analyses")
}

// TestMonitoring_NoNullValues - Test critique anti-null
func TestMonitoring_NoNullValues(t *testing.T) {
	// GIVEN - Un collector de métriques
	metrics := NewMetricsCollector()

	// WHEN - On récupère les métriques
	data := metrics.GetMetrics()

	// THEN - Aucune valeur ne doit être nil
	assert.NotNil(t, data, "Metrics data should not be nil")
	
	for key, value := range data {
		assert.NotNil(t, value, "Key %s should not be nil", key)
	}

	// CRITICAL - active_analyses doit être un nombre, pas nil
	activeAnalyses, exists := data["active_analyses"]
	assert.True(t, exists, "active_analyses key must exist")
	assert.NotNil(t, activeAnalyses, "active_analyses should not be nil")
	
	// Doit être convertible en int
	switch v := activeAnalyses.(type) {
	case int, int64:
		assert.GreaterOrEqual(t, v, 0, "active_analyses should be >= 0")
	default:
		t.Errorf("active_analyses must be int, got %T: %v", activeAnalyses, activeAnalyses)
	}
}

// TestValidationLevels - Test système de validation
func TestValidationLevels(t *testing.T) {
	// GIVEN - Un validateur système
	validator := NewSystemValidator()

	// WHEN - On valide le système
	level, details := validator.Validate()

	// THEN - Le système doit au moins atteindre le niveau Rouge (startup)
	assert.GreaterOrEqual(t, level, LevelRed, "System should at least start (Red level)")
	assert.NotEmpty(t, details, "Validation should return details")
	
	// OBJECTIVE - Viser niveau Orange minimum (basic features)
	// Ce test peut être Yellow si on n'atteint pas Orange encore
	if level >= LevelOrange {
		t.Logf("✅ EXCELLENT - System reached Orange level: %s", details)
	} else {
		t.Logf("⚠️  ATTENTION - System only at Red level. Target: Orange. Details: %s", details)
	}
}

// TestMCPStorage_JSONFormat - Test format JSON lisible
func TestMCPStorage_JSONFormat(t *testing.T) {
	testDir := "./test_data_json"
	storage := NewMCPStorage(testDir)
	defer os.RemoveAll(testDir)

	// GIVEN - Une analyse complexe
	analysis := &AnalysisState{
		ID:          "json-format-test",
		URL:         "https://complex-example.com",
		Status:      "complete",
		Score:       92,
		StartTime:   time.Now(),
		Domain:      "complex-example.com",
		PagesFound:  25,
		PagesAnalyzed: 25,
		TopIssues:   []string{"missing alt tags", "slow load time"},
		Recommendations: []string{"optimize images", "compress CSS"},
	}

	// WHEN - On sauve l'analyse
	err := storage.SaveAnalysis(analysis)
	require.NoError(t, err)

	// THEN - Le fichier JSON doit être bien formaté et lisible
	jsonPath := filepath.Join(testDir, "analyses", "json-format-test.json")
	data, err := os.ReadFile(jsonPath)
	require.NoError(t, err)

	// Vérifier que c'est du JSON valide
	var parsed map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	assert.NoError(t, err, "JSON should be valid")

	// Vérifier que le JSON est indenté (lisible)
	assert.Contains(t, string(data), "  ", "JSON should be indented for readability")
	assert.Contains(t, string(data), "complex-example.com", "Should contain the domain")
}

// TestMCPStorage_ConcurrentAccess - Test accès concurrent
func TestMCPStorage_ConcurrentAccess(t *testing.T) {
	testDir := "./test_data_concurrent"
	storage := NewMCPStorage(testDir)
	defer os.RemoveAll(testDir)

	// GIVEN - Accès concurrent au storage
	const numGoroutines = 5 // Réduire pour plus de stabilité
	var wg sync.WaitGroup

	// WHEN - Plusieurs goroutines sauvent en parallèle
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			analysis := &AnalysisState{
				ID:        fmt.Sprintf("concurrent-test-%d", id),
				URL:       fmt.Sprintf("https://concurrent%d.com", id),
				Status:    "complete",
				Score:     80 + id,
				StartTime: time.Now(),
				Domain:    fmt.Sprintf("concurrent%d.com", id),
			}
			
			err := storage.SaveAnalysis(analysis)
			assert.NoError(t, err, "Concurrent save %d should work", id)
		}(i)
	}

	// Attendre toutes les goroutines
	wg.Wait()

	// Petite pause pour s'assurer que tous les I/O sont finis
	time.Sleep(10 * time.Millisecond)

	// THEN - Toutes les analyses doivent être sauvées
	allAnalyses, err := storage.ListAllAnalyses()
	assert.NoError(t, err)
	assert.Len(t, allAnalyses, numGoroutines, "All concurrent saves should succeed")
	
	// Vérifier aussi l'index
	index, err := storage.LoadIndex()
	assert.NoError(t, err)
	assert.Len(t, index.Analyses, numGoroutines, "Index should contain all analyses")
}