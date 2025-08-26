package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// StorageManager gestionnaire de persistance des résultats
type StorageManager struct {
	db     *sql.DB
	dbPath string
}

// StoredAnalysis analyse stockée en base
type StoredAnalysis struct {
	ID              int64                   `json:"id"`
	TaskID          string                  `json:"task_id"`
	URL             string                  `json:"url"`
	Domain          string                  `json:"domain"`
	AnalysisType    string                  `json:"analysis_type"`
	Status          string                  `json:"status"`
	OverallScore    float64                 `json:"overall_score"`
	ResultData      string                  `json:"result_data"` // JSON sérialisé
	CreatedAt       time.Time               `json:"created_at"`
	ProcessingTime  int64                   `json:"processing_time"` // en millisecondes
}

// AnalysisHistory historique d'analyses pour une URL
type AnalysisHistory struct {
	URL       string            `json:"url"`
	Domain    string            `json:"domain"`
	Analyses  []StoredAnalysis  `json:"analyses"`
	Trend     ScoreTrend        `json:"trend"`
	LastScore float64           `json:"last_score"`
	FirstScore float64          `json:"first_score"`
	Count     int               `json:"count"`
}

// ScoreTrend tendance des scores
type ScoreTrend struct {
	Direction string  `json:"direction"` // up, down, stable
	Change    float64 `json:"change"`    // variation en points
	Period    string  `json:"period"`    // période de comparaison
}

// NewStorageManager crée un nouveau gestionnaire de stockage
func NewStorageManager(dbPath string) (*StorageManager, error) {
	storage := &StorageManager{
		dbPath: dbPath,
	}

	if err := storage.initDatabase(); err != nil {
		return nil, fmt.Errorf("erreur initialisation base de données: %w", err)
	}

	return storage, nil
}

// initDatabase initialise la base de données
func (sm *StorageManager) initDatabase() error {
	var err error
	sm.db, err = sql.Open("sqlite3", sm.dbPath)
	if err != nil {
		return fmt.Errorf("erreur ouverture base: %w", err)
	}

	// Tester la connexion
	if err := sm.db.Ping(); err != nil {
		return fmt.Errorf("erreur ping base: %w", err)
	}

	// Créer les tables
	if err := sm.createTables(); err != nil {
		return fmt.Errorf("erreur création tables: %w", err)
	}

	log.Printf("Base de données initialisée: %s", sm.dbPath)
	return nil
}

// createTables crée les tables nécessaires
func (sm *StorageManager) createTables() error {
	// Table des analyses
	analysisTable := `
	CREATE TABLE IF NOT EXISTS analyses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT NOT NULL,
		url TEXT NOT NULL,
		domain TEXT NOT NULL,
		analysis_type TEXT NOT NULL,
		status TEXT NOT NULL,
		overall_score REAL NOT NULL,
		result_data TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		processing_time INTEGER NOT NULL,
		UNIQUE(task_id)
	);
	`

	if _, err := sm.db.Exec(analysisTable); err != nil {
		return fmt.Errorf("erreur création table analyses: %w", err)
	}

	// Index pour les requêtes fréquentes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_analyses_url ON analyses(url);",
		"CREATE INDEX IF NOT EXISTS idx_analyses_domain ON analyses(domain);",
		"CREATE INDEX IF NOT EXISTS idx_analyses_created_at ON analyses(created_at);",
		"CREATE INDEX IF NOT EXISTS idx_analyses_score ON analyses(overall_score);",
	}

	for _, indexSQL := range indexes {
		if _, err := sm.db.Exec(indexSQL); err != nil {
			return fmt.Errorf("erreur création index: %w", err)
		}
	}

	return nil
}

// SaveAnalysis sauvegarde un résultat d'analyse
func (sm *StorageManager) SaveAnalysis(result *UnifiedAnalysisResult) error {
	// Sérialiser les données en JSON
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("erreur sérialisation JSON: %w", err)
	}

	// Insérer en base
	query := `
	INSERT INTO analyses (
		task_id, url, domain, analysis_type, status, 
		overall_score, result_data, created_at, processing_time
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = sm.db.Exec(query,
		result.TaskID,
		result.URL,
		result.Domain,
		"full", // Type par défaut
		string(result.Status),
		result.OverallScore,
		string(resultJSON),
		result.AnalyzedAt,
		result.ProcessingTime.Milliseconds(),
	)

	if err != nil {
		return fmt.Errorf("erreur insertion base: %w", err)
	}

	log.Printf("Analyse sauvegardée: %s (Score: %.1f)", result.URL, result.OverallScore)
	return nil
}

// GetAnalysis récupère une analyse par son ID de tâche
func (sm *StorageManager) GetAnalysis(taskID string) (*UnifiedAnalysisResult, error) {
	query := `
	SELECT result_data FROM analyses WHERE task_id = ?
	`

	var resultData string
	err := sm.db.QueryRow(query, taskID).Scan(&resultData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("analyse non trouvée: %s", taskID)
		}
		return nil, fmt.Errorf("erreur récupération: %w", err)
	}

	// Désérialiser
	var result UnifiedAnalysisResult
	if err := json.Unmarshal([]byte(resultData), &result); err != nil {
		return nil, fmt.Errorf("erreur désérialisation: %w", err)
	}

	return &result, nil
}

// GetAnalysisHistory récupère l'historique d'analyses pour une URL
func (sm *StorageManager) GetAnalysisHistory(url string, limit int) (*AnalysisHistory, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
	SELECT id, task_id, url, domain, analysis_type, status, 
		   overall_score, created_at, processing_time
	FROM analyses 
	WHERE url = ? 
	ORDER BY created_at DESC 
	LIMIT ?
	`

	rows, err := sm.db.Query(query, url, limit)
	if err != nil {
		return nil, fmt.Errorf("erreur requête historique: %w", err)
	}
	defer rows.Close()

	var analyses []StoredAnalysis
	var domain string

	for rows.Next() {
		var analysis StoredAnalysis
		err := rows.Scan(
			&analysis.ID,
			&analysis.TaskID,
			&analysis.URL,
			&analysis.Domain,
			&analysis.AnalysisType,
			&analysis.Status,
			&analysis.OverallScore,
			&analysis.CreatedAt,
			&analysis.ProcessingTime,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lecture ligne: %w", err)
		}

		analyses = append(analyses, analysis)
		if domain == "" {
			domain = analysis.Domain
		}
	}

	if len(analyses) == 0 {
		return nil, fmt.Errorf("aucune analyse trouvée pour: %s", url)
	}

	// Calculer la tendance
	trend := sm.calculateTrend(analyses)

	history := &AnalysisHistory{
		URL:        url,
		Domain:     domain,
		Analyses:   analyses,
		Trend:      trend,
		LastScore:  analyses[0].OverallScore,
		FirstScore: analyses[len(analyses)-1].OverallScore,
		Count:      len(analyses),
	}

	return history, nil
}

// GetDomainAnalyses récupère les analyses pour un domaine
func (sm *StorageManager) GetDomainAnalyses(domain string, limit int) ([]StoredAnalysis, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `
	SELECT id, task_id, url, domain, analysis_type, status,
		   overall_score, created_at, processing_time
	FROM analyses 
	WHERE domain = ? 
	ORDER BY created_at DESC 
	LIMIT ?
	`

	rows, err := sm.db.Query(query, domain, limit)
	if err != nil {
		return nil, fmt.Errorf("erreur requête domaine: %w", err)
	}
	defer rows.Close()

	var analyses []StoredAnalysis
	for rows.Next() {
		var analysis StoredAnalysis
		err := rows.Scan(
			&analysis.ID,
			&analysis.TaskID,
			&analysis.URL,
			&analysis.Domain,
			&analysis.AnalysisType,
			&analysis.Status,
			&analysis.OverallScore,
			&analysis.CreatedAt,
			&analysis.ProcessingTime,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lecture ligne domaine: %w", err)
		}

		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// GetTopAnalyses récupère les analyses avec les meilleurs scores
func (sm *StorageManager) GetTopAnalyses(limit int) ([]StoredAnalysis, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
	SELECT id, task_id, url, domain, analysis_type, status,
		   overall_score, created_at, processing_time
	FROM analyses 
	ORDER BY overall_score DESC 
	LIMIT ?
	`

	rows, err := sm.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("erreur requête top analyses: %w", err)
	}
	defer rows.Close()

	var analyses []StoredAnalysis
	for rows.Next() {
		var analysis StoredAnalysis
		err := rows.Scan(
			&analysis.ID,
			&analysis.TaskID,
			&analysis.URL,
			&analysis.Domain,
			&analysis.AnalysisType,
			&analysis.Status,
			&analysis.OverallScore,
			&analysis.CreatedAt,
			&analysis.ProcessingTime,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lecture top ligne: %w", err)
		}

		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// GetRecentAnalyses récupère les analyses récentes
func (sm *StorageManager) GetRecentAnalyses(limit int) ([]StoredAnalysis, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `
	SELECT id, task_id, url, domain, analysis_type, status,
		   overall_score, created_at, processing_time
	FROM analyses 
	ORDER BY created_at DESC 
	LIMIT ?
	`

	rows, err := sm.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("erreur requête analyses récentes: %w", err)
	}
	defer rows.Close()

	var analyses []StoredAnalysis
	for rows.Next() {
		var analysis StoredAnalysis
		err := rows.Scan(
			&analysis.ID,
			&analysis.TaskID,
			&analysis.URL,
			&analysis.Domain,
			&analysis.AnalysisType,
			&analysis.Status,
			&analysis.OverallScore,
			&analysis.CreatedAt,
			&analysis.ProcessingTime,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lecture récente ligne: %w", err)
		}

		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// GetStorageStats récupère les statistiques de stockage
func (sm *StorageManager) GetStorageStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Nombre total d'analyses
	var totalAnalyses int
	err := sm.db.QueryRow("SELECT COUNT(*) FROM analyses").Scan(&totalAnalyses)
	if err != nil {
		return nil, fmt.Errorf("erreur comptage analyses: %w", err)
	}
	stats["total_analyses"] = totalAnalyses

	// Nombre de domaines uniques
	var uniqueDomains int
	err = sm.db.QueryRow("SELECT COUNT(DISTINCT domain) FROM analyses").Scan(&uniqueDomains)
	if err != nil {
		return nil, fmt.Errorf("erreur comptage domaines: %w", err)
	}
	stats["unique_domains"] = uniqueDomains

	// Score moyen
	var avgScore sql.NullFloat64
	err = sm.db.QueryRow("SELECT AVG(overall_score) FROM analyses").Scan(&avgScore)
	if err != nil {
		return nil, fmt.Errorf("erreur calcul score moyen: %w", err)
	}
	if avgScore.Valid {
		stats["average_score"] = avgScore.Float64
	} else {
		stats["average_score"] = 0.0
	}

	// Première et dernière analyse (avec conversion string)
	var firstAnalysisStr, lastAnalysisStr sql.NullString
	err = sm.db.QueryRow("SELECT MIN(created_at), MAX(created_at) FROM analyses").Scan(&firstAnalysisStr, &lastAnalysisStr)
	if err != nil {
		return nil, fmt.Errorf("erreur dates analyses: %w", err)
	}
	
	if firstAnalysisStr.Valid {
		if firstTime, err := time.Parse(time.RFC3339, firstAnalysisStr.String); err == nil {
			stats["first_analysis"] = firstTime
		}
	}
	if lastAnalysisStr.Valid {
		if lastTime, err := time.Parse(time.RFC3339, lastAnalysisStr.String); err == nil {
			stats["last_analysis"] = lastTime
		}
	}

	// Analyses par statut
	statusQuery := `
	SELECT status, COUNT(*) 
	FROM analyses 
	GROUP BY status
	`
	rows, err := sm.db.Query(statusQuery)
	if err != nil {
		return nil, fmt.Errorf("erreur requête statuts: %w", err)
	}
	defer rows.Close()

	statusCounts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("erreur lecture statut: %w", err)
		}
		statusCounts[status] = count
	}
	stats["status_counts"] = statusCounts

	return stats, nil
}

// CleanupOldAnalyses supprime les anciennes analyses
func (sm *StorageManager) CleanupOldAnalyses(olderThanDays int) (int, error) {
	if olderThanDays <= 0 {
		olderThanDays = 90 // Par défaut, supprimer après 90 jours
	}

	cutoffDate := time.Now().AddDate(0, 0, -olderThanDays)

	result, err := sm.db.Exec("DELETE FROM analyses WHERE created_at < ?", cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("erreur suppression anciennes analyses: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("erreur récupération lignes supprimées: %w", err)
	}

	log.Printf("Nettoyage terminé: %d analyses supprimées (plus de %d jours)", rowsAffected, olderThanDays)
	return int(rowsAffected), nil
}

// Close ferme la connexion à la base de données
func (sm *StorageManager) Close() error {
	if sm.db != nil {
		return sm.db.Close()
	}
	return nil
}

// calculateTrend calcule la tendance des scores
func (sm *StorageManager) calculateTrend(analyses []StoredAnalysis) ScoreTrend {
	if len(analyses) < 2 {
		return ScoreTrend{
			Direction: "stable",
			Change:    0,
			Period:    "insufficient_data",
		}
	}

	// Comparer les deux analyses les plus récentes
	latest := analyses[0].OverallScore
	previous := analyses[1].OverallScore
	change := latest - previous

	direction := "stable"
	if change > 1.0 {
		direction = "up"
	} else if change < -1.0 {
		direction = "down"
	}

	return ScoreTrend{
		Direction: direction,
		Change:    change,
		Period:    "last_analysis",
	}
}

// GetAnalysisById récupère une analyse spécifique avec ses données complètes
func (sm *StorageManager) GetAnalysisById(analysisID int64) (*StoredAnalysis, error) {
	query := `
	SELECT id, task_id, url, domain, analysis_type, status, overall_score, 
	       result_data, created_at, processing_time
	FROM analyses 
	WHERE id = ?`

	var analysis StoredAnalysis
	err := sm.db.QueryRow(query, analysisID).Scan(
		&analysis.ID,
		&analysis.TaskID,
		&analysis.URL,
		&analysis.Domain,
		&analysis.AnalysisType,
		&analysis.Status,
		&analysis.OverallScore,
		&analysis.ResultData,
		&analysis.CreatedAt,
		&analysis.ProcessingTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("analyse avec ID %d non trouvée", analysisID)
		}
		return nil, fmt.Errorf("erreur récupération analyse: %w", err)
	}

	return &analysis, nil
}