package api

import (
	"sync"
	"time"
)

// AnalyzeRequest - Structure pour la requête d'analyse
type AnalyzeRequest struct {
	URL string `json:"url"`
}

// AnalyzeResponse - Structure pour la réponse d'analyse
type AnalyzeResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// StatusResponse - Structure pour la réponse de status
type StatusResponse struct {
	ID             string `json:"id"`
	URL            string `json:"url"`
	Status         string `json:"status"` // "analyzing", "complete"
	Progress       int    `json:"progress"`
	PagesFound     int    `json:"pages_found"`
	PagesAnalyzed  int    `json:"pages_analyzed"`
	IssuesFound    int    `json:"issues_found"`
	EstimatedTime  string `json:"estimated_time"`
}

// ResultsResponse - Structure pour les résultats finaux
type ResultsResponse struct {
	Score      int            `json:"score"`
	PagesCount int            `json:"pages_count"`
	Issues     []ResultIssue  `json:"issues"`
	Warnings   []ResultWarning `json:"warnings"`
	Analysis   AnalysisResult `json:"analysis"`
}

// ResultIssue - Problème détecté
type ResultIssue struct {
	Title       string   `json:"title"`
	Count       int      `json:"count"`
	Description string   `json:"description"`
	Pages       []string `json:"pages"`
	Solution    string   `json:"solution"`
}

// ResultWarning - Avertissement détecté
type ResultWarning struct {
	Title       string `json:"title"`
	Count       int    `json:"count"`
	Description string `json:"description"`
}

// AnalysisResult - Résultat détaillé d'analyse
type AnalysisResult struct {
	Domain         string              `json:"domain"`
	Date           string              `json:"date"`
	Score          int                 `json:"score"`
	PagesAnalyzed  int                 `json:"pages_analyzed"`
	AnalysisTime   string              `json:"analysis_time"`
	CriticalIssues int                 `json:"critical_issues"`
	Warnings       int                 `json:"warnings"`
	AISuggestions  []AISuggestion      `json:"ai_suggestions"`
}

// AISuggestion - Suggestion IA
type AISuggestion struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
}

// AnalysisState - État interne d'une analyse
type AnalysisState struct {
	ID             string    `json:"id"`
	URL            string    `json:"url"`
	Status         string    `json:"status"`
	Progress       int       `json:"progress"`
	PagesFound     int       `json:"pages_found"`
	PagesAnalyzed  int       `json:"pages_analyzed"`
	IssuesFound    int       `json:"issues_found"`
	EstimatedTime  string    `json:"estimated_time"`
	StartTime      time.Time `json:"start_time"`
	Results        *ResultsResponse `json:"results,omitempty"`
}

// AnalysisStore - Store thread-safe pour les analyses
type AnalysisStore struct {
	mu       sync.RWMutex
	analyses map[string]*AnalysisState
}

// NewAnalysisStore - Créer un nouveau store
func NewAnalysisStore() *AnalysisStore {
	return &AnalysisStore{
		analyses: make(map[string]*AnalysisState),
	}
}

// Store - Instance globale du store
var Store = NewAnalysisStore()

// Set - Stocker une analyse
func (s *AnalysisStore) Set(id string, analysis *AnalysisState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.analyses[id] = analysis
}

// Get - Récupérer une analyse
func (s *AnalysisStore) Get(id string) (*AnalysisState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	analysis, exists := s.analyses[id]
	return analysis, exists
}

// Update - Mettre à jour une analyse
func (s *AnalysisStore) Update(id string, updater func(*AnalysisState)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if analysis, exists := s.analyses[id]; exists {
		updater(analysis)
	}
}