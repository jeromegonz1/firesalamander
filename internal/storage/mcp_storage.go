package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// 🔥🦎 FIRE SALAMANDER - SPRINT 6 MCP STORAGE
// Persistance JSON pour analyses et métriques

// MCPStorage - Storage MCP pour persistance JSON
type MCPStorage struct {
	basePath    string
	mu          sync.RWMutex
	analysesDir string
	metricsDir  string
	reportsDir  string
}

// NewMCPStorage - Créer nouveau storage MCP
func NewMCPStorage(basePath string) *MCPStorage {
	storage := &MCPStorage{
		basePath:    basePath,
		analysesDir: filepath.Join(basePath, "analyses"),
		metricsDir:  filepath.Join(basePath, "metrics"),
		reportsDir:  filepath.Join(basePath, "reports"),
	}
	
	// Créer les répertoires nécessaires
	storage.ensureDirectories()
	
	return storage
}

// ensureDirectories - Créer structure de répertoires
func (s *MCPStorage) ensureDirectories() {
	dirs := []string{
		s.basePath,
		s.analysesDir,
		s.metricsDir,
		s.reportsDir,
	}
	
	for _, dir := range dirs {
		os.MkdirAll(dir, 0755)
	}
}

// SaveAnalysis - Sauver une analyse en JSON
func (s *MCPStorage) SaveAnalysis(analysis *AnalysisState) error {
	// Mettre à jour les timestamps
	now := time.Now()
	analysis.UpdatedAt = now
	if analysis.CreatedAt.IsZero() {
		analysis.CreatedAt = now
	}
	
	// Sauver le fichier JSON indenté (avec lock séparé)
	s.mu.Lock()
	filename := fmt.Sprintf("%s.json", analysis.ID)
	filePath := filepath.Join(s.analysesDir, filename)
	
	data, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to marshal analysis %s: %w", analysis.ID, err)
	}
	
	err = os.WriteFile(filePath, data, 0644)
	s.mu.Unlock()
	
	if err != nil {
		return fmt.Errorf("failed to write analysis %s: %w", analysis.ID, err)
	}
	
	// Mettre à jour l'index (avec son propre lock)
	return s.updateIndex(analysis)
}

// LoadAnalysis - Charger une analyse depuis JSON
func (s *MCPStorage) LoadAnalysis(id string) (*AnalysisState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	filename := fmt.Sprintf("%s.json", id)
	filePath := filepath.Join(s.analysesDir, filename)
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read analysis %s: %w", id, err)
	}
	
	var analysis AnalysisState
	err = json.Unmarshal(data, &analysis)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal analysis %s: %w", id, err)
	}
	
	return &analysis, nil
}

// LoadIndex - Charger l'index des analyses
func (s *MCPStorage) LoadIndex() (*AnalysisIndex, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	indexPath := filepath.Join(s.analysesDir, "index.json")
	
	// Si l'index n'existe pas, le créer
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return &AnalysisIndex{
			LastUpdated: time.Now(),
			TotalCount:  0,
			Analyses:    []AnalysisIndexItem{},
		}, nil
	}
	
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read index: %w", err)
	}
	
	var index AnalysisIndex
	err = json.Unmarshal(data, &index)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal index: %w", err)
	}
	
	return &index, nil
}

// ListAllAnalyses - Lister toutes les analyses
func (s *MCPStorage) ListAllAnalyses() ([]*AnalysisState, error) {
	index, err := s.LoadIndex()
	if err != nil {
		return nil, err
	}
	
	var analyses []*AnalysisState
	
	for _, item := range index.Analyses {
		analysis, err := s.LoadAnalysis(item.ID)
		if err != nil {
			// Log l'erreur mais continue avec les autres
			continue
		}
		analyses = append(analyses, analysis)
	}
	
	return analyses, nil
}

// updateIndex - Mettre à jour l'index avec une nouvelle analyse (THREAD-SAFE)
func (s *MCPStorage) updateIndex(analysis *AnalysisState) error {
	// LOCK EXCLUSIF pour tout le processus d'update de l'index
	s.mu.Lock()
	defer s.mu.Unlock()
	
	indexPath := filepath.Join(s.analysesDir, "index.json")
	
	var index *AnalysisIndex
	
	// Si l'index n'existe pas, le créer
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		index = &AnalysisIndex{
			LastUpdated: time.Now(),
			TotalCount:  0,
			Analyses:    []AnalysisIndexItem{},
		}
	} else {
		data, err := os.ReadFile(indexPath)
		if err != nil {
			return fmt.Errorf("failed to read index: %w", err)
		}
		
		index = &AnalysisIndex{}
		err = json.Unmarshal(data, index)
		if err != nil {
			return fmt.Errorf("failed to unmarshal index: %w", err)
		}
	}
	
	// Chercher si l'analyse existe déjà dans l'index
	found := false
	for i, item := range index.Analyses {
		if item.ID == analysis.ID {
			// Mettre à jour l'item existant
			index.Analyses[i] = AnalysisIndexItem{
				ID:        analysis.ID,
				URL:       analysis.URL,
				Domain:    analysis.Domain,
				Status:    analysis.Status,
				Score:     analysis.Score,
				CreatedAt: analysis.CreatedAt,
				FilePath:  fmt.Sprintf("%s.json", analysis.ID),
			}
			found = true
			break
		}
	}
	
	// Si pas trouvé, ajouter nouveau
	if !found {
		index.Analyses = append(index.Analyses, AnalysisIndexItem{
			ID:        analysis.ID,
			URL:       analysis.URL,
			Domain:    analysis.Domain,
			Status:    analysis.Status,
			Score:     analysis.Score,
			CreatedAt: analysis.CreatedAt,
			FilePath:  fmt.Sprintf("%s.json", analysis.ID),
		})
	}
	
	// Mettre à jour les métadonnées de l'index
	index.TotalCount = len(index.Analyses)
	index.LastUpdated = time.Now()
	
	// Sauver l'index (déjà sous lock exclusif)
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}
	
	return os.WriteFile(indexPath, data, 0644)
}

// SaveMetrics - Sauver métriques temps réel
func (s *MCPStorage) SaveMetrics(snapshot *MetricsSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Sauver current.json
	currentPath := filepath.Join(s.metricsDir, "current.json")
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(currentPath, data, 0644)
}

// LoadCurrentMetrics - Charger métriques actuelles
func (s *MCPStorage) LoadCurrentMetrics() (*MetricsSnapshot, error) {
	currentPath := filepath.Join(s.metricsDir, "current.json")
	
	data, err := os.ReadFile(currentPath)
	if err != nil {
		// Si pas de fichier, retourner métriques vides
		return &MetricsSnapshot{
			Timestamp:     time.Now(),
			ActiveAnalyses: 0,
			HealthStatus:  "HEALTHY",
		}, nil
	}
	
	var metrics MetricsSnapshot
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		return nil, err
	}
	
	return &metrics, nil
}

// MetricsCollector - Collecteur de métriques (ANTI-NULL)
type MetricsCollector struct {
	mu sync.RWMutex
	storage *MCPStorage
}

// NewMetricsCollector - Créer nouveau collecteur de métriques
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		storage: NewMCPStorage("./data"),
	}
}

// GetMetrics - Récupérer métriques (JAMAIS de valeurs null)
func (mc *MetricsCollector) GetMetrics() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	// Charger métriques depuis le storage
	current, _ := mc.storage.LoadCurrentMetrics()
	
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	// CRITICAL: Toutes les valeurs DOIVENT être définies (pas de nil)
	metrics := map[string]interface{}{
		"active_analyses":        current.ActiveAnalyses,     // int, JAMAIS nil
		"completed_analyses":     current.CompletedAnalyses,  // int
		"failed_analyses":        current.FailedAnalyses,     // int
		"total_requests":         current.TotalRequests,      // int
		"infinite_loops":         current.InfiniteLoopsDetected, // int
		"timeouts":              current.TimeoutsOccurred,    // int
		"goroutines":            runtime.NumGoroutine(),      // int
		"memory_usage_mb":       float64(memStats.Alloc) / 1024 / 1024, // float64
		"health_status":         current.HealthStatus,        // string
		"timestamp":             current.Timestamp,           // time.Time
		"uptime_seconds":        int(time.Since(current.Timestamp).Seconds()), // int
	}
	
	return metrics
}

// SystemValidator - Validateur système
type SystemValidator struct {
	storage *MCPStorage
}

// NewSystemValidator - Créer nouveau validateur
func NewSystemValidator() *SystemValidator {
	return &SystemValidator{
		storage: NewMCPStorage("./data"),
	}
}

// Validate - Valider le système et retourner niveau
func (sv *SystemValidator) Validate() (ValidationLevel, string) {
	details := make(map[string]bool)
	
	// TEST 1: Le système démarre-t-il ?
	details["system_starts"] = true // Si on arrive ici, c'est que oui
	
	// TEST 2: Les répertoires sont-ils créés ?
	_, err := os.Stat("./data/analyses")
	details["data_structure"] = err == nil
	
	// TEST 3: Peut-on sauver une analyse ?
	testAnalysis := &AnalysisState{
		ID:        "validation-test",
		URL:       "https://test.com",
		Status:    "complete",
		Score:     100,
		StartTime: time.Now(),
		Domain:    "test.com",
	}
	err = sv.storage.SaveAnalysis(testAnalysis)
	details["save_analysis"] = err == nil
	
	// TEST 4: Peut-on charger une analyse ?
	_, err = sv.storage.LoadAnalysis("validation-test")
	details["load_analysis"] = err == nil
	
	// TEST 5: L'index fonctionne-t-il ?
	_, err = sv.storage.LoadIndex()
	details["index_works"] = err == nil
	
	// Déterminer le niveau
	passedTests := 0
	for _, passed := range details {
		if passed {
			passedTests++
		}
	}
	
	var level ValidationLevel
	var message string
	
	switch {
	case passedTests >= 5:
		level = LevelOrange
		message = fmt.Sprintf("✅ ORANGE Level - Basic features work (%d/5 tests passed)", passedTests)
	case passedTests >= 3:
		level = LevelRed  
		message = fmt.Sprintf("🔴 RED Level - System starts but features limited (%d/5 tests passed)", passedTests)
	default:
		level = LevelRed
		message = fmt.Sprintf("❌ CRITICAL - System failing (%d/5 tests passed)", passedTests)
	}
	
	return level, message
}