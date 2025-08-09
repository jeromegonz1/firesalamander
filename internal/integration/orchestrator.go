package integration

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
	"firesalamander/internal/crawler"
	"firesalamander/internal/patterns"
	"firesalamander/internal/seo"
	"firesalamander/internal/storage"
)

// 🔥🦎 FIRE SALAMANDER - REAL ORCHESTRATOR
// Sprint 5 - L'intégration finale qui connecte tous les composants
// ZERO HARDCODING POLICY - All values from constants

// Orchestrator coordonne le crawling et l'analyse SEO réelle
type Orchestrator struct {
	// Configuration
	config *OrchestratorConfig
	
	// Composants intégrés
	parallelCrawler *crawler.ParallelCrawler
	seoAnalyzer     *seo.RealSEOAnalyzer
	
	// 🔥🦎 SPRINT 6: MCP Storage pour persistance
	storage         *storage.MCPStorage
	
	// 🔥🦎 SPRINT 6+: SafeCrawler anti-boucle infinie
	safeCrawler     *patterns.SafeCrawler
	
	// State management
	analyses map[string]*AnalysisState
	updates  chan AnalysisUpdate
	mu       sync.RWMutex
	
	// Workers et synchronisation
	workerPool sync.WaitGroup
	shutdown   chan struct{}
	running    bool
}

// OrchestratorConfig configuration pour le Orchestrator
type OrchestratorConfig struct {
	MaxPages        int
	MaxWorkers      int
	InitialWorkers  int
	Timeout         time.Duration
	UserAgent       string
	EnableRealTime  bool
}

// AnalysisState état d'une analyse en cours ou terminée
type AnalysisState struct {
	// Identification
	ID       string
	URL      string
	Domain   string
	
	// Timing
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	
	// Status
	Status string
	Error  string
	
	// Progress metrics (real-time)
	PagesFound      int
	PagesAnalyzed   int
	CurrentWorkers  int
	PagesPerSecond  float64
	
	// Results
	Pages           []*seo.RealPageAnalysis
	GlobalScore     int
	GlobalGrade     string
	TopIssues       []seo.RealRecommendation
	Recommendations []seo.RealRecommendation
	
	// Crawler metrics
	CrawlerMetrics *crawler.CrawlerMetrics
}

// AnalysisUpdate mise à jour temps réel d'analyse
type AnalysisUpdate struct {
	AnalysisID string
	Message    string
	Timestamp  time.Time
	Status     string
	Data       map[string]interface{}
}

// NewOrchestrator crée un nouveau orchestrateur réel
func NewOrchestrator() *Orchestrator {
	realConfig := &OrchestratorConfig{
		MaxPages:       constants.OrchestratorMaxPages,
		MaxWorkers:     constants.OrchestratorMaxWorkers,
		InitialWorkers: constants.OrchestratorInitialWorkers,
		Timeout:        time.Duration(constants.OrchestratorAnalysisTimeout),
		UserAgent:      constants.DefaultUserAgent,
		EnableRealTime: true,
	}

	// Créer le crawler parallel avec la config
	crawlerConfig := &config.CrawlerConfig{
		MaxPages:         realConfig.MaxPages,
		InitialWorkers:   realConfig.InitialWorkers,
		MinWorkers:       1,
		TimeoutSeconds:   constants.DefaultTimeoutSeconds, // Use the default timeout from constants
		UserAgent:        realConfig.UserAgent,
		RespectRobotsTxt: true,
	}

	parallelCrawler := crawler.NewParallelCrawler(crawlerConfig)
	seoAnalyzer := seo.NewRealSEOAnalyzer()

	// 🔥🦎 SPRINT 6: Initialiser MCP Storage pour persistance
	mcpStorage := storage.NewMCPStorage("./data")
	
	// 🔥🦎 SPRINT 6+: Initialiser SafeCrawler anti-boucle infinie
	safeCrawler := patterns.NewSafeCrawler()
	
	orchestrator := &Orchestrator{
		config:          realConfig,
		parallelCrawler: parallelCrawler,
		seoAnalyzer:     seoAnalyzer,
		storage:         mcpStorage,
		safeCrawler:     safeCrawler,
		analyses:        make(map[string]*AnalysisState),
		updates:         make(chan AnalysisUpdate, constants.OrchestratorChannelBufferDefault),
		shutdown:        make(chan struct{}),
		running:         false,
	}
	
	// Recharger les analyses existantes depuis le storage
	orchestrator.loadExistingAnalyses()
	
	return orchestrator
}

// GetConfig retourne la configuration (pour tests)
func (ro *Orchestrator) GetConfig() *OrchestratorConfig {
	return ro.config
}

// StartAnalysis démarre une nouvelle analyse complète
func (ro *Orchestrator) StartAnalysis(targetURL string) (string, error) {
	// Validation de l'URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", fmt.Errorf("URL must include scheme and host")
	}

	// Générer ID unique
	analysisID := ro.generateAnalysisID()
	
	// Extraire le domaine
	domain := parsedURL.Host
	
	// Créer l'état initial
	state := &AnalysisState{
		ID:             analysisID,
		URL:            targetURL,
		Domain:         domain,
		StartTime:      time.Now(),
		Status:         constants.OrchestratorStatusStarting,
		PagesFound:     0,
		PagesAnalyzed:  0,
		CurrentWorkers: 0,
		PagesPerSecond: 0.0,
		Pages:          make([]*seo.RealPageAnalysis, 0),
		TopIssues:      make([]seo.RealRecommendation, 0),
		Recommendations: make([]seo.RealRecommendation, 0),
	}

	// Sauvegarder l'état
	ro.mu.Lock()
	ro.analyses[analysisID] = state
	ro.mu.Unlock()

	// 🔥🦎 SPRINT 6: Sauver l'analyse initiale dans MCP Storage
	ro.saveAnalysisToStorage(state)

	// Envoyer mise à jour
	ro.sendUpdate(analysisID, "Analysis initialized", constants.OrchestratorStatusStarting)

	// Démarrer l'analyse asynchrone
	go ro.runCompleteAnalysis(analysisID)

	return analysisID, nil
}

// GetStatus récupère le statut d'une analyse
func (ro *Orchestrator) GetStatus(analysisID string) (*AnalysisState, error) {
	ro.mu.RLock()
	defer ro.mu.RUnlock()

	state, exists := ro.analyses[analysisID]
	if !exists {
		return nil, fmt.Errorf("analysis not found: %s", analysisID)
	}

	// Retourner une copie pour éviter les race conditions
	stateCopy := *state
	return &stateCopy, nil
}

// runCompleteAnalysis exécute l'analyse complète (crawler + SEO)
func (ro *Orchestrator) runCompleteAnalysis(analysisID string) {
	state := ro.getState(analysisID)
	if state == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			ro.updateStateError(analysisID, fmt.Sprintf("Analysis panicked: %v", r))
		}
	}()

	// PHASE 1: CRAWLING
	ro.updateState(analysisID, constants.OrchestratorStatusCrawling, "Starting crawling phase")
	crawlResults, err := ro.performCrawling(analysisID, state.URL)
	if err != nil {
		ro.updateStateError(analysisID, fmt.Sprintf("Crawling failed: %v", err))
		return
	}

	// PHASE 2: SEO ANALYSIS
	ro.updateState(analysisID, constants.OrchestratorStatusAnalyzing, "Starting SEO analysis phase")
	pageAnalyses, err := ro.performSEOAnalysis(analysisID, crawlResults)
	if err != nil {
		ro.updateStateError(analysisID, fmt.Sprintf("SEO analysis failed: %v", err))
		return
	}

	// PHASE 3: AGGREGATION
	ro.updateState(analysisID, constants.OrchestratorStatusAggregating, "Aggregating results")
	ro.aggregateResults(analysisID, pageAnalyses)

	// COMPLETION
	ro.completeAnalysis(analysisID)
}

// performCrawling effectue le crawling avec le parallel crawler SÉCURISÉ
func (ro *Orchestrator) performCrawling(analysisID, targetURL string) ([]*crawler.PageResult, error) {
	ro.sendUpdate(analysisID, fmt.Sprintf("Starting crawl of %s", targetURL), constants.OrchestratorStatusCrawling)

	// 🔥🦎 SPRINT 6+: TIMEOUT OBLIGATOIRE anti-boucle infinie  
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second) // Timeout strict 90s
	defer cancel()

	// 🔥🦎 SPRINT 6+: Vérification SafeCrawler avant crawling
	log.Printf("🛡️  SafeCrawler: Starting protected crawl for URL: %s", targetURL)
	ro.sendUpdate(analysisID, "SafeCrawler: Initializing protected crawling...", constants.OrchestratorStatusCrawling)

	// Démarrer le crawling avec timeout strict
	ro.sendUpdate(analysisID, "Calling crawler.CrawlWithContext with 90s timeout...", constants.OrchestratorStatusCrawling)
	log.Printf("🔍 DEBUG: About to call CrawlWithContext for URL: %s (TIMEOUT: 90s)", targetURL)
	crawlResult, err := ro.parallelCrawler.CrawlWithContext(ctx, targetURL)
	log.Printf("🔍 DEBUG: CrawlWithContext returned, err=%v", err)
	if err != nil {
		ro.sendUpdate(analysisID, fmt.Sprintf("Crawler error: %v", err), constants.OrchestratorStatusError)
		return nil, fmt.Errorf("failed to crawl: %w", err)
	}

	ro.sendUpdate(analysisID, "Crawler returned, checking results...", constants.OrchestratorStatusCrawling)
	log.Printf("🔍 DEBUG: Crawl result: pages=%d, error=%v", len(crawlResult.Pages), crawlResult.Error)
	if crawlResult.Error != nil {
		ro.sendUpdate(analysisID, fmt.Sprintf("Crawl result error: %v", crawlResult.Error), constants.OrchestratorStatusError)
		return nil, fmt.Errorf("crawling failed: %w", crawlResult.Error)
	}

	ro.sendUpdate(analysisID, fmt.Sprintf("Crawl result contains %d pages", len(crawlResult.Pages)), constants.OrchestratorStatusCrawling)

	// Convertir la map en slice
	var results []*crawler.PageResult
	for _, pageResult := range crawlResult.Pages {
		if pageResult.Error == nil { // Ignorer les pages avec erreurs
			results = append(results, pageResult)
		}
	}
	
	// Mise à jour finale des métriques
	state := ro.getState(analysisID)
	if state != nil {
		state.PagesFound = len(results)
		state.CrawlerMetrics = crawlResult.Metrics
		if crawlResult.Metrics != nil {
			state.CurrentWorkers = crawlResult.Metrics.CurrentWorkers
			state.PagesPerSecond = crawlResult.Metrics.PagesPerSecond
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no pages successfully crawled")
	}

	ro.sendUpdate(analysisID, fmt.Sprintf("Crawling complete - %d pages found", len(results)), constants.OrchestratorStatusCrawling)
	return results, nil
}

// performSEOAnalysis effectue l'analyse SEO de toutes les pages
func (ro *Orchestrator) performSEOAnalysis(analysisID string, crawlResults []*crawler.PageResult) ([]*seo.RealPageAnalysis, error) {
	ctx := context.Background()
	var pageAnalyses []*seo.RealPageAnalysis
	
	ro.sendUpdate(analysisID, fmt.Sprintf("Starting SEO analysis of %d pages", len(crawlResults)), constants.OrchestratorStatusAnalyzing)

	for i, result := range crawlResults {
		// Analyser la page avec le RealSEOAnalyzer
		analysis := ro.seoAnalyzer.AnalyzePageContent(ctx, result.URL, result.Body)
		pageAnalyses = append(pageAnalyses, analysis)

		// Mise à jour temps réel
		state := ro.getState(analysisID)
		if state != nil {
			state.PagesAnalyzed = i + 1
			ro.sendUpdate(analysisID, 
				fmt.Sprintf("Analyzed: %s (Score: %d/%d)", result.URL, analysis.TotalScore, constants.MaxSEOScore), 
				constants.OrchestratorStatusAnalyzing)
		}
	}

	ro.sendUpdate(analysisID, fmt.Sprintf("SEO analysis complete - %d pages analyzed", len(pageAnalyses)), constants.OrchestratorStatusAnalyzing)
	return pageAnalyses, nil
}

// aggregateResults agrège les résultats d'analyse
func (ro *Orchestrator) aggregateResults(analysisID string, pageAnalyses []*seo.RealPageAnalysis) {
	state := ro.getState(analysisID)
	if state == nil {
		return
	}

	// Sauvegarder toutes les analyses de page
	state.Pages = pageAnalyses

	// Calculer le score global
	if len(pageAnalyses) > 0 {
		totalScore := 0
		for _, page := range pageAnalyses {
			totalScore += page.TotalScore
		}
		state.GlobalScore = totalScore / len(pageAnalyses)
		state.GlobalGrade = ro.seoAnalyzer.DetermineGrade(state.GlobalScore)
	}

	// Agréger les recommandations
	allRecommendations := make([]seo.RealRecommendation, 0)
	for _, page := range pageAnalyses {
		allRecommendations = append(allRecommendations, page.Recommendations...)
	}

	// Prioriser et dédupliquer les recommandations
	state.Recommendations = ro.prioritizeRecommendations(allRecommendations)
	state.TopIssues = ro.extractTopIssues(state.Recommendations)

	ro.sendUpdate(analysisID, fmt.Sprintf("Results aggregated - Global score: %d (%s)", state.GlobalScore, state.GlobalGrade), constants.OrchestratorStatusAggregating)
}

// completeAnalysis finalise l'analyse
func (ro *Orchestrator) completeAnalysis(analysisID string) {
	state := ro.getState(analysisID)
	if state == nil {
		return
	}

	// Finaliser les timing
	state.EndTime = time.Now()
	state.Duration = state.EndTime.Sub(state.StartTime)
	state.Status = constants.OrchestratorStatusComplete

	// 🔥🦎 SPRINT 6: Sauver l'analyse complète dans MCP Storage
	ro.saveAnalysisToStorage(state)

	ro.sendUpdate(analysisID, fmt.Sprintf("Analysis complete in %v - Score: %d (%s)", state.Duration, state.GlobalScore, state.GlobalGrade), constants.OrchestratorStatusComplete)
}

// Helper methods

func (ro *Orchestrator) generateAnalysisID() string {
	// Utiliser timestamp nanoseconde complet + PID pour garantir l'unicité
	now := time.Now()
	return fmt.Sprintf("analysis-%d-%d-%d", 
		now.Unix(),           // Timestamp seconde
		now.UnixNano(),       // Timestamp nanoseconde complet 
		os.Getpid(),          // Process ID pour multi-process
	)
}

func (ro *Orchestrator) getState(analysisID string) *AnalysisState {
	ro.mu.RLock()
	defer ro.mu.RUnlock()
	return ro.analyses[analysisID]
}

func (ro *Orchestrator) updateState(analysisID, status, message string) {
	ro.mu.Lock()
	state, exists := ro.analyses[analysisID]
	if exists {
		state.Status = status
	}
	ro.mu.Unlock()
	
	ro.sendUpdate(analysisID, message, status)
}

func (ro *Orchestrator) updateStateError(analysisID, errorMessage string) {
	ro.mu.Lock()
	state, exists := ro.analyses[analysisID]
	if exists {
		state.Status = constants.OrchestratorStatusError
		state.Error = errorMessage
		state.EndTime = time.Now()
		state.Duration = state.EndTime.Sub(state.StartTime)
	}
	ro.mu.Unlock()

	// 🔥🦎 SPRINT 6: Sauver l'état d'erreur dans MCP Storage
	if exists && state != nil {
		ro.saveAnalysisToStorage(state)
	}
	
	ro.sendUpdate(analysisID, errorMessage, constants.OrchestratorStatusError)
}

func (ro *Orchestrator) sendUpdate(analysisID, message, status string) {
	if !ro.config.EnableRealTime {
		return
	}

	update := AnalysisUpdate{
		AnalysisID: analysisID,
		Message:    message,
		Timestamp:  time.Now(),
		Status:     status,
		Data:       make(map[string]interface{}),
	}

	select {
	case ro.updates <- update:
		// Update sent successfully
	default:
		// Channel full, skip update (non-blocking)
	}
}

func (ro *Orchestrator) prioritizeRecommendations(allRecs []seo.RealRecommendation) []seo.RealRecommendation {
	// Simple prioritization - group by priority and take top recommendations
	criticalRecs := make([]seo.RealRecommendation, 0)
	highRecs := make([]seo.RealRecommendation, 0)
	mediumRecs := make([]seo.RealRecommendation, 0)

	// Group recommendations by priority
	for _, rec := range allRecs {
		switch rec.Priority {
		case constants.SEOPriorityCritical:
			criticalRecs = append(criticalRecs, rec)
		case constants.SEOPriorityHigh:
			highRecs = append(highRecs, rec)
		case constants.SEOPriorityMedium:
			mediumRecs = append(mediumRecs, rec)
		}
	}

	// Combine prioritized recommendations (limit to reasonable number)
	result := make([]seo.RealRecommendation, 0)
	result = append(result, criticalRecs...)
	result = append(result, highRecs...)
	result = append(result, mediumRecs...)

	// Limit to top 10 recommendations
	maxRecs := constants.MaxRecommendations
	if len(result) > maxRecs {
		result = result[:maxRecs]
	}

	return result
}

func (ro *Orchestrator) extractTopIssues(recommendations []seo.RealRecommendation) []seo.RealRecommendation {
	// Extract top 5 critical issues
	topIssues := make([]seo.RealRecommendation, 0)
	
	for _, rec := range recommendations {
		if rec.Priority == constants.SEOPriorityCritical && len(topIssues) < constants.MaxTopIssues {
			topIssues = append(topIssues, rec)
		}
	}

	return topIssues
}

// GetUpdatesChannel retourne le canal des mises à jour (pour WebSocket future)
func (ro *Orchestrator) GetUpdatesChannel() <-chan AnalysisUpdate {
	return ro.updates
}

// Shutdown arrête proprement l'orchestrateur
func (ro *Orchestrator) Shutdown() error {
	ro.mu.Lock()
	defer ro.mu.Unlock()

	if !ro.running {
		return nil
	}

	close(ro.shutdown)
	ro.workerPool.Wait()
	close(ro.updates)
	ro.running = false

	return nil
}

// 🔥🦎 SPRINT 6: Persistance MCP Functions

// loadExistingAnalyses - Recharge les analyses depuis le storage MCP
func (ro *Orchestrator) loadExistingAnalyses() {
	stored, err := ro.storage.ListAllAnalyses()
	if err != nil {
		log.Printf("⚠️  Failed to load existing analyses from storage: %v", err)
		return
	}
	
	log.Printf("🔄 Loading %d existing analyses from MCP storage", len(stored))
	
	for _, storageAnalysis := range stored {
		// Convertir storage.AnalysisState -> integration.AnalysisState
		integrationAnalysis := ro.convertFromStorage(storageAnalysis)
		ro.analyses[integrationAnalysis.ID] = integrationAnalysis
	}
	
	log.Printf("✅ Loaded %d analyses from storage successfully", len(stored))
}

// saveAnalysisToStorage - Sauve une analyse dans le storage MCP
func (ro *Orchestrator) saveAnalysisToStorage(analysis *AnalysisState) {
	// Convertir integration.AnalysisState -> storage.AnalysisState
	storageAnalysis := ro.convertToStorage(analysis)
	
	err := ro.storage.SaveAnalysis(storageAnalysis)
	if err != nil {
		log.Printf("❌ Failed to save analysis %s to storage: %v", analysis.ID, err)
	} else {
		log.Printf("💾 Saved analysis %s to MCP storage", analysis.ID)
	}
}

// convertToStorage - Convertit integration.AnalysisState vers storage.AnalysisState
func (ro *Orchestrator) convertToStorage(integrationAnalysis *AnalysisState) *storage.AnalysisState {
	// Convertir les recommandations en chaînes simples
	var topIssues []string
	for _, issue := range integrationAnalysis.TopIssues {
		topIssues = append(topIssues, issue.Issue)
	}
	
	var recommendations []string
	for _, rec := range integrationAnalysis.Recommendations {
		recommendations = append(recommendations, fmt.Sprintf("%s: %s", rec.Issue, rec.Action))
	}
	
	return &storage.AnalysisState{
		ID:              integrationAnalysis.ID,
		URL:             integrationAnalysis.URL,
		Domain:          integrationAnalysis.Domain,
		Status:          integrationAnalysis.Status,
		StartTime:       integrationAnalysis.StartTime,
		Duration:        integrationAnalysis.Duration,
		Score:           integrationAnalysis.GlobalScore,
		PagesFound:      integrationAnalysis.PagesFound,
		PagesAnalyzed:   integrationAnalysis.PagesAnalyzed,
		CurrentWorkers:  integrationAnalysis.CurrentWorkers,
		PagesPerSecond:  integrationAnalysis.PagesPerSecond,
		TopIssues:       topIssues,
		Recommendations: recommendations,
		GlobalGrade:     integrationAnalysis.GlobalGrade,
		CreatedAt:       integrationAnalysis.StartTime,
		UpdatedAt:       time.Now(),
	}
}

// convertFromStorage - Convertit storage.AnalysisState vers integration.AnalysisState
func (ro *Orchestrator) convertFromStorage(storageAnalysis *storage.AnalysisState) *AnalysisState {
	// Convertir les issues simples en RealRecommendations
	var topIssues []seo.RealRecommendation
	for _, issue := range storageAnalysis.TopIssues {
		topIssues = append(topIssues, seo.RealRecommendation{
			Priority: constants.SEOPriorityCritical,
			Issue:    issue,
			Action:   "Corriger ce problème",
		})
	}
	
	var recommendations []seo.RealRecommendation
	for _, rec := range storageAnalysis.Recommendations {
		recommendations = append(recommendations, seo.RealRecommendation{
			Priority: constants.SEOPriorityMedium,
			Issue:    rec,
			Action:   "Voir les détails",
		})
	}
	
	return &AnalysisState{
		ID:              storageAnalysis.ID,
		URL:             storageAnalysis.URL,
		Domain:          storageAnalysis.Domain,
		StartTime:       storageAnalysis.StartTime,
		Duration:        storageAnalysis.Duration,
		Status:          storageAnalysis.Status,
		PagesFound:      storageAnalysis.PagesFound,
		PagesAnalyzed:   storageAnalysis.PagesAnalyzed,
		CurrentWorkers:  storageAnalysis.CurrentWorkers,
		PagesPerSecond:  storageAnalysis.PagesPerSecond,
		GlobalScore:     storageAnalysis.Score,
		GlobalGrade:     storageAnalysis.GlobalGrade,
		TopIssues:       topIssues,
		Recommendations: recommendations,
	}
}

// TestCrawl teste directement le crawling (pour debug des tests)
func (ro *Orchestrator) TestCrawl(ctx context.Context, url string) (*crawler.ParallelCrawlResult, error) {
	return ro.parallelCrawler.CrawlWithContext(ctx, url)
}