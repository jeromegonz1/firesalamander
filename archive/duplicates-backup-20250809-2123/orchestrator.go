package integration

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
	"firesalamander/internal/agents/crawler"
	"firesalamander/internal/agents/semantic"
	"firesalamander/internal/seo"
)

// Orchestrator coordonne toutes les analyses Fire Salamander
type Orchestrator struct {
	config *config.Config
	
	// Composants principaux
	crawler         *crawler.Crawler
	semanticAnalyzer *semantic.TestSemanticAnalyzer // Version simplifiée pour l'intégration
	seoAnalyzer     *seo.SEOAnalyzer
	storage         *StorageManager
	
	// Gestion des tâches
	taskQueue       chan *AnalysisTask
	resultsChan     chan *UnifiedAnalysisResult
	workers         int
	workerPool      sync.WaitGroup
	
	// Statistiques
	stats           *AnalysisStats
	statsMutex      sync.RWMutex
	
	// État
	isRunning       bool
	shutdownChan    chan struct{}
	mutex           sync.RWMutex
}

// AnalysisTask tâche d'analyse à effectuer
type AnalysisTask struct {
	ID          string                 `json:"id"`
	URL         string                 `json:"url"`
	Type        AnalysisType           `json:"type"`
	Options     AnalysisOptions        `json:"options"`
	Priority    TaskPriority           `json:"priority"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Status      TaskStatus             `json:"status"`
	Error       string                 `json:"error,omitempty"`
	
	// Canal pour retourner le résultat
	ResultChan  chan *UnifiedAnalysisResult `json:"-"`
}

// UnifiedAnalysisResult résultat unifié de toutes les analyses
type UnifiedAnalysisResult struct {
	// Métadonnées de base
	TaskID              string                          `json:"task_id"`
	URL                 string                          `json:"url"`
	Domain              string                          `json:"domain"`
	AnalyzedAt          time.Time                       `json:"analyzed_at"`
	ProcessingTime      time.Duration                   `json:"processing_time"`
	
	// Résultats des différents modules
	CrawlerResult       *crawler.CrawlResult            `json:"crawler_result,omitempty"`
	SemanticAnalysis    *semantic.AnalysisResult        `json:"semantic_analysis,omitempty"`
	SEOAnalysis         *seo.SEOAnalysisResult          `json:"seo_analysis,omitempty"`
	
	// Analyse unifiée et recommandations
	UnifiedMetrics      UnifiedMetrics                  `json:"unified_metrics"`
	CrossModuleInsights []CrossModuleInsight            `json:"cross_module_insights"`
	PriorityActions     []PriorityAction                `json:"priority_actions"`
	
	// Scoring global
	OverallScore        float64                         `json:"overall_score"`
	CategoryScores      map[string]float64              `json:"category_scores"`
	
	// Statut et erreurs
	Status              AnalysisStatus                  `json:"status"`
	Errors              []AnalysisError                 `json:"errors,omitempty"`
	Warnings            []string                        `json:"warnings,omitempty"`
}

// AnalysisType type d'analyse à effectuer
type AnalysisType string

const (
	AnalysisTypeFull     AnalysisType = "full"     // Analyse complète
	AnalysisTypeSemantic AnalysisType = "semantic" // Analyse sémantique uniquement
	AnalysisTypeSEO      AnalysisType = "seo"      // Analyse SEO uniquement
	AnalysisTypeQuick    AnalysisType = "quick"    // Analyse rapide (SEO + base)
)

// AnalysisOptions options d'analyse
type AnalysisOptions struct {
	IncludeCrawling     bool          `json:"include_crawling"`
	MaxDepth            int           `json:"max_depth"`
	FollowRedirects     bool          `json:"follow_redirects"`
	AnalyzeImages       bool          `json:"analyze_images"`
	AnalyzePerformance  bool          `json:"analyze_performance"`
	UseAIEnrichment     bool          `json:"use_ai_enrichment"`
	Timeout             time.Duration `json:"timeout"`
}

// TaskPriority priorité de la tâche
type TaskPriority int

const (
	PriorityLow    TaskPriority = 1
	PriorityNormal TaskPriority = 2
	PriorityHigh   TaskPriority = 3
	PriorityUrgent TaskPriority = 4
)

// TaskStatus statut de la tâche
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusRunning    TaskStatus = "running"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// AnalysisStatus statut de l'analyse
type AnalysisStatus string

const (
	AnalysisStatusSuccess AnalysisStatus = "success"
	AnalysisStatusPartial AnalysisStatus = "partial"
	AnalysisStatusFailed  AnalysisStatus = "failed"
)

// UnifiedMetrics métriques unifiées cross-modules
type UnifiedMetrics struct {
	ContentQualityScore    float64 `json:"content_quality_score"`
	TechnicalHealthScore   float64 `json:"technical_health_score"`
	SEOReadinessScore      float64 `json:"seo_readiness_score"`
	UserExperienceScore    float64 `json:"user_experience_score"`
	MobileFriendlinessScore float64 `json:"mobile_friendliness_score"`
	PerformanceScore       float64 `json:"performance_score"`
}

// CrossModuleInsight insight basé sur plusieurs modules
type CrossModuleInsight struct {
	Type        string   `json:"type"`
	Severity    string   `json:"severity"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Evidence    []string `json:"evidence"`
	Modules     []string `json:"modules"`
	Impact      string   `json:"impact"`
}

// PriorityAction action prioritaire recommandée
type PriorityAction struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Priority    string        `json:"priority"`
	Impact      string        `json:"impact"`
	Effort      string        `json:"effort"`
	Module      string        `json:"module"`
	EstimatedTime string      `json:"estimated_time"`
	Dependencies  []string    `json:"dependencies"`
}

// AnalysisError erreur d'analyse
type AnalysisError struct {
	Module      string    `json:"module"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	Recoverable bool      `json:"recoverable"`
}

// AnalysisStats statistiques d'analyse
type AnalysisStats struct {
	TotalTasks       int64         `json:"total_tasks"`
	CompletedTasks   int64         `json:"completed_tasks"`
	FailedTasks      int64         `json:"failed_tasks"`
	AverageTime      time.Duration `json:"average_time"`
	LastAnalysis     time.Time     `json:"last_analysis"`
	ActiveTasks      int64         `json:"active_tasks"`
}

// NewOrchestrator crée un nouvel orchestrateur
func NewOrchestrator(cfg *config.Config) (*Orchestrator, error) {
	// Créer la configuration du crawler
	crawlerConfig := &crawler.Config{
		UserAgent:     cfg.Crawler.UserAgent,
		Workers:       cfg.Crawler.Workers,
		RateLimit:     cfg.Crawler.RateLimit,
		MaxDepth:      3,
		MaxPages:      100,
		Timeout:       constants.ClientTimeout,
		RetryAttempts: 3,
		RetryDelay:    constants.DefaultRetryDelay,
		RespectRobots: true,
		EnableCache:   true,
	}

	// Initialiser les composants
	crawlerInstance, err := crawler.New(crawlerConfig)
	if err != nil {
		return nil, fmt.Errorf("erreur initialisation crawler: %w", err)
	}

	semanticAnalyzer := semantic.NewTestSemanticAnalyzer()
	seoAnalyzer := seo.NewSEOAnalyzer()
	
	// Initialiser le storage manager
	storage, err := NewStorageManager("fire_salamander.db")
	if err != nil {
		return nil, fmt.Errorf("erreur initialisation storage: %w", err)
	}

	orchestrator := &Orchestrator{
		config:           cfg,
		crawler:          crawlerInstance,
		semanticAnalyzer: semanticAnalyzer,
		seoAnalyzer:      seoAnalyzer,
		storage:          storage,
		taskQueue:        make(chan *AnalysisTask, 100),
		resultsChan:      make(chan *UnifiedAnalysisResult, 100),
		workers:          cfg.Crawler.Workers,
		stats:            &AnalysisStats{},
		shutdownChan:     make(chan struct{}),
	}

	return orchestrator, nil
}

// Start démarre l'orchestrateur
func (o *Orchestrator) Start(ctx context.Context) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if o.isRunning {
		return fmt.Errorf("orchestrateur déjà en cours d'exécution")
	}

	log.Printf("Démarrage de l'orchestrateur Fire Salamander avec %d workers", o.workers)

	// Démarrer les workers
	for i := 0; i < o.workers; i++ {
		o.workerPool.Add(1)
		go o.worker(ctx, i)
	}

	o.isRunning = true
	log.Printf("Orchestrateur Fire Salamander démarré avec succès")

	return nil
}

// Stop arrête l'orchestrateur
func (o *Orchestrator) Stop() error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if !o.isRunning {
		return nil
	}

	log.Printf("Arrêt de l'orchestrateur Fire Salamander")

	close(o.shutdownChan)
	close(o.taskQueue)

	// Attendre la fin des workers
	o.workerPool.Wait()

	o.isRunning = false
	log.Printf("Orchestrateur Fire Salamander arrêté")

	return nil
}

// AnalyzeURL lance une analyse complète d'une URL
func (o *Orchestrator) AnalyzeURL(ctx context.Context, targetURL string, analysisType AnalysisType, options AnalysisOptions) (*UnifiedAnalysisResult, error) {
	taskID := fmt.Sprintf("task_%d", time.Now().UnixNano())

	task := &AnalysisTask{
		ID:         taskID,
		URL:        targetURL,
		Type:       analysisType,
		Options:    options,
		Priority:   PriorityNormal,
		CreatedAt:  time.Now(),
		Status:     TaskStatusPending,
		ResultChan: make(chan *UnifiedAnalysisResult, 1),
	}

	// Envoyer la tâche dans la queue
	select {
	case o.taskQueue <- task:
		log.Printf("Tâche %s ajoutée à la queue pour %s", taskID, targetURL)
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Attendre le résultat
	select {
	case result := <-task.ResultChan:
		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// worker traite les tâches d'analyse
func (o *Orchestrator) worker(ctx context.Context, workerID int) {
	defer o.workerPool.Done()

	log.Printf("Worker %d démarré", workerID)

	for {
		select {
		case task, ok := <-o.taskQueue:
			if !ok {
				log.Printf("Worker %d: queue fermée, arrêt", workerID)
				return
			}

			log.Printf("Worker %d: traitement tâche %s pour %s", workerID, task.ID, task.URL)
			result := o.processTask(ctx, task)
			
			// Envoyer le résultat
			select {
			case task.ResultChan <- result:
			default:
				log.Printf("Worker %d: impossible d'envoyer le résultat pour %s", workerID, task.ID)
			}

		case <-o.shutdownChan:
			log.Printf("Worker %d: signal d'arrêt reçu", workerID)
			return

		case <-ctx.Done():
			log.Printf("Worker %d: contexte annulé", workerID)
			return
		}
	}
}

// processTask traite une tâche d'analyse
func (o *Orchestrator) processTask(ctx context.Context, task *AnalysisTask) *UnifiedAnalysisResult {
	startTime := time.Now()
	task.Status = TaskStatusRunning
	task.StartedAt = &startTime

	result := &UnifiedAnalysisResult{
		TaskID:         task.ID,
		URL:            task.URL,
		AnalyzedAt:     startTime,
		CategoryScores: make(map[string]float64),
		Errors:         []AnalysisError{},
		Warnings:       []string{},
	}

	// Extraire le domaine
	if domain, err := o.extractDomain(task.URL); err == nil {
		result.Domain = domain
	}

	var wg sync.WaitGroup
	var crawlResult *crawler.CrawlResult
	var semanticResult *semantic.AnalysisResult
	var seoResult *seo.SEOAnalysisResult

	// 1. Crawling (si demandé)
	if task.Options.IncludeCrawling && (task.Type == AnalysisTypeFull || task.Type == AnalysisTypeQuick) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if cr, err := o.performCrawling(ctx, task); err != nil {
				result.Errors = append(result.Errors, AnalysisError{
					Module:      constants.OrchestratorAgentNameCrawler,
					Type:        constants.OrchestratorErrorTypeCrawling,
					Message:     err.Error(),
					Timestamp:   time.Now(),
					Recoverable: true,
				})
			} else {
				crawlResult = cr
			}
		}()
	}

	// 2. Analyse sémantique
	if task.Type == AnalysisTypeFull || task.Type == AnalysisTypeSemantic {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if sr, err := o.performSemanticAnalysis(ctx, task); err != nil {
				result.Errors = append(result.Errors, AnalysisError{
					Module:      constants.OrchestratorAnalysisTypeSemantic,
					Type:        constants.OrchestratorErrorTypeSemantic,
					Message:     err.Error(),
					Timestamp:   time.Now(),
					Recoverable: true,
				})
			} else {
				semanticResult = sr
			}
		}()
	}

	// 3. Analyse SEO
	if task.Type == AnalysisTypeFull || task.Type == AnalysisTypeSEO || task.Type == AnalysisTypeQuick {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if seoRes, err := o.performSEOAnalysis(ctx, task); err != nil {
				result.Errors = append(result.Errors, AnalysisError{
					Module:      constants.OrchestratorAgentNameSEO,
					Type:        constants.OrchestratorErrorTypeSEO,
					Message:     err.Error(),
					Timestamp:   time.Now(),
					Recoverable: true,
				})
			} else {
				seoResult = seoRes
			}
		}()
	}

	// Attendre toutes les analyses
	wg.Wait()

	// Assigner les résultats
	result.CrawlerResult = crawlResult
	result.SemanticAnalysis = semanticResult
	result.SEOAnalysis = seoResult

	// 4. Analyse unifiée et cross-module
	o.performUnifiedAnalysis(result)

	// 5. Finaliser
	result.ProcessingTime = time.Since(startTime)
	task.Status = TaskStatusCompleted
	task.CompletedAt = &result.AnalyzedAt

	// Déterminer le statut final
	if len(result.Errors) == 0 {
		result.Status = AnalysisStatusSuccess
	} else if crawlResult != nil || semanticResult != nil || seoResult != nil {
		result.Status = AnalysisStatusPartial
		result.Warnings = append(result.Warnings, fmt.Sprintf("Analyse partielle - %d erreurs", len(result.Errors)))
	} else {
		result.Status = AnalysisStatusFailed
		task.Status = TaskStatusFailed
	}

	// Mettre à jour les statistiques
	o.updateStats(result)

	log.Printf("Tâche %s terminée - Statut: %s, Durée: %v, Score: %.1f", 
		task.ID, result.Status, result.ProcessingTime, result.OverallScore)

	return result
}

// performCrawling effectue le crawling
func (o *Orchestrator) performCrawling(ctx context.Context, task *AnalysisTask) (*crawler.CrawlResult, error) {
	log.Printf("Début crawling pour %s", task.URL)
	
	// Configuration du crawler basée sur les options
	maxDepth := task.Options.MaxDepth
	if maxDepth == 0 {
		maxDepth = 2 // Valeur par défaut
	}

	// Utiliser CrawlPage pour une seule page
	result, err := o.crawler.CrawlPage(ctx, task.URL)
	if err != nil {
		return nil, fmt.Errorf("erreur crawling page: %w", err)
	}

	log.Printf("Crawling terminé - Page: %s, Status: %d", task.URL, result.StatusCode)
	return result, nil
}

// performSemanticAnalysis effectue l'analyse sémantique
func (o *Orchestrator) performSemanticAnalysis(ctx context.Context, task *AnalysisTask) (*semantic.AnalysisResult, error) {
	log.Printf("Début analyse sémantique pour %s", task.URL)

	// Pour l'intégration, on utilise du contenu HTML simulé
	// Dans une vraie implémentation, on récupérerait le HTML via HTTP
	htmlContent := o.fetchHTMLContent(ctx, task.URL)
	if htmlContent == "" {
		return nil, fmt.Errorf("impossible de récupérer le contenu HTML")
	}

	result, err := o.semanticAnalyzer.AnalyzePage(ctx, task.URL, htmlContent)
	if err != nil {
		return nil, fmt.Errorf("erreur analyse sémantique: %w", err)
	}

	log.Printf("Analyse sémantique terminée - Score: %.1f", result.SEOScore.Overall)
	return result, nil
}

// performSEOAnalysis effectue l'analyse SEO
func (o *Orchestrator) performSEOAnalysis(ctx context.Context, task *AnalysisTask) (*seo.SEOAnalysisResult, error) {
	log.Printf("Début analyse SEO pour %s", task.URL)

	result, err := o.seoAnalyzer.AnalyzePage(ctx, task.URL)
	if err != nil {
		return nil, fmt.Errorf("erreur analyse SEO: %w", err)
	}

	log.Printf("Analyse SEO terminée - Score: %.1f", result.OverallScore)
	return result, nil
}

// performUnifiedAnalysis effectue l'analyse unifiée
func (o *Orchestrator) performUnifiedAnalysis(result *UnifiedAnalysisResult) {
	// 1. Calculer les métriques unifiées
	result.UnifiedMetrics = o.calculateUnifiedMetrics(result)

	// 2. Générer les insights cross-module
	result.CrossModuleInsights = o.generateCrossModuleInsights(result)

	// 3. Identifier les actions prioritaires
	result.PriorityActions = o.identifyPriorityActions(result)

	// 4. Calculer le score global unifié
	result.OverallScore = o.calculateUnifiedScore(result)

	// 5. Calculer les scores par catégorie
	result.CategoryScores = o.calculateCategoryScores(result)
}

// Fonctions utilitaires et de calcul (implémentations simplifiées)

func (o *Orchestrator) extractDomain(targetURL string) (string, error) {
	domain := o.extractDomainSimple(targetURL)
	if domain == "" {
		return "", fmt.Errorf("URL invalide")
	}
	return domain, nil
}

func (o *Orchestrator) extractDomainSimple(targetURL string) string {
	// Implémentation simple - dans la vraie vie, utiliser net/url
	if len(targetURL) > constants.HTTPSPrefixLength && targetURL[:constants.HTTPSPrefixLength] == constants.HTTPSPrefix {
		domain := targetURL[constants.HTTPSPrefixLength:]
		if idx := strings.Index(domain, "/"); idx != -1 {
			domain = domain[:idx]
		}
		return domain
	}
	if len(targetURL) > constants.HTTPPrefixLength && targetURL[:constants.HTTPPrefixLength] == constants.HTTPPrefix {
		domain := targetURL[constants.HTTPPrefixLength:]
		if idx := strings.Index(domain, "/"); idx != -1 {
			domain = domain[:idx]
		}
		return domain
	}
	return ""
}

func (o *Orchestrator) fetchHTMLContent(ctx context.Context, targetURL string) string {
	// Contenu HTML simulé pour les tests
	return `<!DOCTYPE html>
<html lang="fr">
<head>
	<title>Page de test Fire Salamander</title>
	<meta name="description" content="Page de test pour l'orchestrateur Fire Salamander avec analyse complète SEO et sémantique.">
</head>
<body>
	<h1>Test Fire Salamander</h1>
	<p>Contenu de test pour valider l'intégration des modules d'analyse.</p>
</body>
</html>`
}

func (o *Orchestrator) calculateUnifiedMetrics(result *UnifiedAnalysisResult) UnifiedMetrics {
	metrics := UnifiedMetrics{}

	// Score de qualité du contenu (basé sur l'analyse sémantique)
	if result.SemanticAnalysis != nil {
		metrics.ContentQualityScore = result.SemanticAnalysis.SEOScore.Overall
	}

	// Score de santé technique (basé sur l'analyse SEO technique)
	if result.SEOAnalysis != nil {
		if techScore, exists := result.SEOAnalysis.CategoryScores[constants.OrchestratorCategoryTechnical]; exists {
			metrics.TechnicalHealthScore = techScore * 100
		}
		if perfScore, exists := result.SEOAnalysis.CategoryScores[constants.OrchestratorAgentNamePerformance]; exists {
			metrics.PerformanceScore = perfScore * 100
		}
		metrics.SEOReadinessScore = result.SEOAnalysis.OverallScore
	}

	// Score d'expérience utilisateur (moyenne pondérée)
	metrics.UserExperienceScore = (metrics.PerformanceScore*0.4 + metrics.ContentQualityScore*0.3 + metrics.TechnicalHealthScore*0.3)

	// Score mobile (basé sur l'audit technique SEO)
	if result.SEOAnalysis != nil {
		metrics.MobileFriendlinessScore = result.SEOAnalysis.TechnicalAudit.Mobile.MobileScore * 100
	}

	return metrics
}

func (o *Orchestrator) generateCrossModuleInsights(result *UnifiedAnalysisResult) []CrossModuleInsight {
	var insights []CrossModuleInsight

	// Insight 1: Cohérence titre/contenu
	if result.SemanticAnalysis != nil && result.SEOAnalysis != nil {
		if result.SEOAnalysis.TagAnalysis.Title.Present && len(result.SemanticAnalysis.LocalAnalysis.Keywords) > 0 {
			insights = append(insights, CrossModuleInsight{
				Type:        constants.OrchestratorInsightContentSEOAlignment,
				Severity:    constants.OrchestratorStatusInfo,
				Title:       "Alignement contenu-SEO détecté",
				Description: "Le titre de la page est cohérent avec les mots-clés identifiés dans le contenu",
				Evidence:    []string{"Titre présent", "Mots-clés identifiés"},
				Modules:     []string{constants.OrchestratorAnalysisTypeSemantic, constants.OrchestratorAgentNameSEO},
				Impact:      constants.OrchestratorImpactPositive,
			})
		}
	}

	// Insight 2: Performance vs Contenu
	if result.SEOAnalysis != nil {
		perfScore := result.UnifiedMetrics.PerformanceScore
		contentScore := result.UnifiedMetrics.ContentQualityScore
		
		if perfScore < 50 && contentScore > 70 {
			insights = append(insights, CrossModuleInsight{
				Type:        constants.OrchestratorInsightPerformanceContentMismatch,
				Severity:    constants.OrchestratorStatusWarning,
				Title:       "Décalage performance-contenu",
				Description: "Bon contenu mais performances techniques faibles",
				Evidence:    []string{fmt.Sprintf("Performance: %.1f%%", perfScore), fmt.Sprintf("Contenu: %.1f%%", contentScore)},
				Modules:     []string{constants.OrchestratorAnalysisTypeSemantic, constants.OrchestratorAgentNameSEO},
				Impact:      constants.OrchestratorImpactNegative,
			})
		}
	}

	return insights
}

func (o *Orchestrator) identifyPriorityActions(result *UnifiedAnalysisResult) []PriorityAction {
	var actions []PriorityAction

	// Actions basées sur l'analyse SEO
	if result.SEOAnalysis != nil {
		for _, rec := range result.SEOAnalysis.Recommendations {
			if len(actions) >= 10 { // Limiter à 10 actions max
				break
			}

			priority := "medium"
			switch rec.Priority {
			case seo.PriorityCritical:
				priority = "critical"
			case seo.PriorityHigh:
				priority = "high"
			case seo.PriorityLow:
				priority = "low"
			}

			actions = append(actions, PriorityAction{
				ID:          rec.ID,
				Title:       rec.Title,
				Description: rec.Description,
				Priority:    priority,
				Impact:      string(rec.Impact),
				Effort:      string(rec.Effort),
				Module:      constants.OrchestratorAgentNameSEO,
				EstimatedTime: constants.OrchestratorTimeVariable,
			})
		}
	}

	return actions
}

func (o *Orchestrator) calculateUnifiedScore(result *UnifiedAnalysisResult) float64 {
	var totalScore float64
	var components int

	// Score sémantique (30%)
	if result.SemanticAnalysis != nil {
		totalScore += result.SemanticAnalysis.SEOScore.Overall * 0.3
		components++
	}

	// Score SEO (70%)
	if result.SEOAnalysis != nil {
		totalScore += result.SEOAnalysis.OverallScore * 0.7
		components++
	}

	if components == 0 {
		return 0
	}

	return totalScore
}

func (o *Orchestrator) calculateCategoryScores(result *UnifiedAnalysisResult) map[string]float64 {
	scores := make(map[string]float64)

	// Reprendre les scores SEO si disponibles
	if result.SEOAnalysis != nil {
		for category, score := range result.SEOAnalysis.CategoryScores {
			scores[category] = score * 100 // Convertir en pourcentage
		}
	}

	// Ajouter les métriques unifiées
	scores[constants.OrchestratorCategoryContentQuality] = result.UnifiedMetrics.ContentQualityScore
	scores[constants.OrchestratorCategoryUserExperience] = result.UnifiedMetrics.UserExperienceScore
	scores[constants.OrchestratorCategoryMobileFriendliness] = result.UnifiedMetrics.MobileFriendlinessScore

	return scores
}

func (o *Orchestrator) updateStats(result *UnifiedAnalysisResult) {
	o.statsMutex.Lock()
	defer o.statsMutex.Unlock()

	o.stats.TotalTasks++
	
	if result.Status == AnalysisStatusSuccess || result.Status == AnalysisStatusPartial {
		o.stats.CompletedTasks++
	} else {
		o.stats.FailedTasks++
	}

	// Calculer le temps moyen
	if o.stats.CompletedTasks > 0 {
		o.stats.AverageTime = time.Duration(
			(int64(o.stats.AverageTime)*o.stats.CompletedTasks + int64(result.ProcessingTime)) / (o.stats.CompletedTasks + 1),
		)
	} else {
		o.stats.AverageTime = result.ProcessingTime
	}

	o.stats.LastAnalysis = result.AnalyzedAt
	
	// Sauvegarder l'analyse dans le storage
	if o.storage != nil {
		err := o.storage.SaveAnalysis(result)
		if err != nil {
			log.Printf("Erreur sauvegarde analyse %s: %v", result.TaskID, err)
		}
	}
}

// GetStats retourne les statistiques actuelles
func (o *Orchestrator) GetStats() *AnalysisStats {
	o.statsMutex.RLock()
	defer o.statsMutex.RUnlock()

	// Retourner une copie
	stats := *o.stats
	return &stats
}

// GetRecentAnalyses retourne la liste des analyses récentes
func (o *Orchestrator) GetRecentAnalyses() []map[string]interface{} {
	if o.storage == nil {
		return []map[string]interface{}{}
	}
	
	analyses, err := o.storage.GetRecentAnalyses(20)
	if err != nil {
		log.Printf("Erreur récupération analyses récentes: %v", err)
		return []map[string]interface{}{}
	}
	
	// Convertir en format JSON générique pour l'API
	result := make([]map[string]interface{}, len(analyses))
	for i, analysis := range analyses {
		result[i] = map[string]interface{}{
			"id":               analysis.ID,
			"task_id":          analysis.TaskID,
			"url":              analysis.URL,
			"domain":           analysis.Domain,
			"analysis_type":    analysis.AnalysisType,
			constants.OrchestratorJSONFieldStatus:           analysis.Status,
			"overall_score":    analysis.OverallScore,
			"created_at":       analysis.CreatedAt,
			"processing_time":  analysis.ProcessingTime,
		}
	}
	
	return result
}

// GetAnalysisDetails retourne les détails complets d'une analyse spécifique
func (o *Orchestrator) GetAnalysisDetails(analysisID int64) (interface{}, error) {
	if o.storage == nil {
		return nil, fmt.Errorf("storage non disponible")
	}
	
	// Récupérer l'analyse complète avec les données JSON
	analysis, err := o.storage.GetAnalysisById(analysisID)
	if err != nil {
		return nil, fmt.Errorf("analyse non trouvée: %w", err)
	}
	
	return analysis, nil
}