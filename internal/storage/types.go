package storage

import "time"

// 🔥🦎 FIRE SALAMANDER - SPRINT 6 TYPES
// Structures pour la persistance MCP

// AnalysisState - État d'une analyse SEO (compatible avec integration.AnalysisState)
type AnalysisState struct {
	ID              string    `json:"id"`
	URL             string    `json:"url"`
	Domain          string    `json:"domain"`
	Status          string    `json:"status"`
	StartTime       time.Time `json:"start_time"`
	Duration        time.Duration `json:"duration,omitempty"`
	Score           int       `json:"score"`
	PagesFound      int       `json:"pages_found"`
	PagesAnalyzed   int       `json:"pages_analyzed"`
	CurrentWorkers  int       `json:"current_workers"`
	PagesPerSecond  float64   `json:"pages_per_second"`
	TopIssues       []string  `json:"top_issues,omitempty"`
	Recommendations []string  `json:"recommendations,omitempty"`
	GlobalGrade     string    `json:"global_grade,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// AnalysisIndex - Index de toutes les analyses (pour navigation rapide)
type AnalysisIndex struct {
	LastUpdated time.Time              `json:"last_updated"`
	TotalCount  int                    `json:"total_count"`
	Analyses    []AnalysisIndexItem    `json:"analyses"`
}

// AnalysisIndexItem - Item dans l'index
type AnalysisIndexItem struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Domain    string    `json:"domain"`
	Status    string    `json:"status"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	FilePath  string    `json:"file_path"`
}

// MetricsSnapshot - Snapshot des métriques système
type MetricsSnapshot struct {
	Timestamp           time.Time `json:"timestamp"`
	ActiveAnalyses      int       `json:"active_analyses"`      // NEVER nil
	CompletedAnalyses   int       `json:"completed_analyses"`   
	FailedAnalyses      int       `json:"failed_analyses"`
	TotalRequests       int       `json:"total_requests"`
	InfiniteLoopsDetected int     `json:"infinite_loops_detected"`
	TimeoutsOccurred    int       `json:"timeouts_occurred"`
	Goroutines          int       `json:"goroutines"`
	MemoryUsageMB       float64   `json:"memory_usage_mb"`
	HealthStatus        string    `json:"health_status"`
}

// ValidationLevel - Niveaux de validation système
type ValidationLevel int

const (
	LevelRed    ValidationLevel = 1 // Système démarre
	LevelOrange ValidationLevel = 2 // Features de base fonctionnent  
	LevelYellow ValidationLevel = 3 // Utilisable en production
	LevelGreen  ValidationLevel = 4 // Production-ready
)

func (vl ValidationLevel) String() string {
	switch vl {
	case LevelRed:
		return "🔴 RED - System Starting"
	case LevelOrange:
		return "🟠 ORANGE - Basic Features"
	case LevelYellow:
		return "🟡 YELLOW - Production Ready"
	case LevelGreen:
		return "🟢 GREEN - Excellent"
	default:
		return "❌ UNKNOWN"
	}
}

// ValidationResult - Résultat de validation
type ValidationResult struct {
	Level   ValidationLevel `json:"level"`
	Message string          `json:"message"`
	Details map[string]bool `json:"details"`
	Timestamp time.Time     `json:"timestamp"`
}