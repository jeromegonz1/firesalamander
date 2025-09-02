package v2

import (
	"context"
	"fmt"
	"sync"
	"time"

	"firesalamander/internal/constants"
	"firesalamander/internal/agents/crawler"
)

// pipelineExecutor implémente l'interface PipelineExecutor
type pipelineExecutor struct {
	mu             sync.RWMutex
	runningAudits  map[string]*auditExecution
	progressMap    map[string]*ExecutionProgress
}

// auditExecution représente l'état d'une exécution d'audit
type auditExecution struct {
	request   *AuditRequest
	ctx       context.Context
	cancel    context.CancelFunc
	startTime time.Time
}

// NewPipelineExecutor crée une nouvelle instance de PipelineExecutor
func NewPipelineExecutor() PipelineExecutor {
	return &pipelineExecutor{
		runningAudits: make(map[string]*auditExecution),
		progressMap:   make(map[string]*ExecutionProgress),
	}
}

// Execute lance l'exécution du pipeline pour un audit
func (p *pipelineExecutor) Execute(ctx context.Context, request *AuditRequest, registry AgentRegistry) (<-chan *PipelineResult, error) {
	// Validation des paramètres
	if request == nil {
		return nil, fmt.Errorf("audit request cannot be nil")
	}
	
	if request.AuditID == "" {
		return nil, fmt.Errorf("audit ID cannot be empty")
	}
	
	if request.SeedURL == "" {
		return nil, fmt.Errorf("seed URL cannot be empty")
	}
	
	if request.MaxPages <= 0 {
		return nil, fmt.Errorf("max pages must be greater than 0")
	}
	
	if registry == nil {
		return nil, fmt.Errorf("agent registry cannot be nil")
	}
	
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Vérifier si l'audit n'est pas déjà en cours
	if _, exists := p.runningAudits[request.AuditID]; exists {
		return nil, fmt.Errorf("audit %s is already running", request.AuditID)
	}
	
	// Créer un contexte cancellable
	auditCtx, cancel := context.WithCancel(ctx)
	
	execution := &auditExecution{
		request:   request,
		ctx:       auditCtx,
		cancel:    cancel,
		startTime: time.Now(),
	}
	
	p.runningAudits[request.AuditID] = execution
	
	// Initialiser la progression
	progress := &ExecutionProgress{
		AuditID:         request.AuditID,
		OverallProgress: 0.0,
		CurrentStep:     constants.PipelineStepInitializing,
		AgentProgresses: make(map[string]float64),
		CompletedAgents: []string{},
		RunningAgents:   []string{},
		PendingAgents:   []string{},
		FailedAgents:    []string{},
	}
	p.progressMap[request.AuditID] = progress
	
	// Créer le channel de résultats
	resultsChan := make(chan *PipelineResult, constants.DefaultStreamBufferSize)
	
	// Démarrer l'exécution asynchrone
	go p.executePipeline(auditCtx, execution, registry, resultsChan)
	
	return resultsChan, nil
}

// GetProgress récupère la progression actuelle d'un audit
func (p *pipelineExecutor) GetProgress(auditID string) (*ExecutionProgress, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	progress, exists := p.progressMap[auditID]
	if !exists {
		return nil, fmt.Errorf("audit %s not found", auditID)
	}
	
	return progress, nil
}

// Cancel annule l'exécution d'un audit
func (p *pipelineExecutor) Cancel(auditID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	execution, exists := p.runningAudits[auditID]
	if !exists {
		return fmt.Errorf("audit %s not found", auditID)
	}
	
	execution.cancel()
	return nil
}

// IsRunning vérifie si un audit est en cours d'exécution
func (p *pipelineExecutor) IsRunning(auditID string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	_, exists := p.runningAudits[auditID]
	return exists
}

// GetRunningAudits retourne la liste des audits en cours
func (p *pipelineExecutor) GetRunningAudits() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	audits := make([]string, 0, len(p.runningAudits))
	for auditID := range p.runningAudits {
		audits = append(audits, auditID)
	}
	
	return audits
}

// executePipeline exécute le pipeline complet
func (p *pipelineExecutor) executePipeline(ctx context.Context, execution *auditExecution, registry AgentRegistry, resultsChan chan<- *PipelineResult) {
	defer close(resultsChan)
	defer p.cleanupExecution(execution.request.AuditID)
	
	request := execution.request
	
	// Étape 1: Crawling
	crawlResult, err := p.executeCrawlStep(ctx, request, resultsChan)
	if err != nil {
		p.sendResult(resultsChan, &PipelineResult{
			Step:   constants.PipelineStepCrawling,
			Status: constants.StatusFailed,
			Error:  err,
		})
		return
	}
	
	// Étape 2: Analyse avec les agents
	err = p.executeAnalysisStep(ctx, request, crawlResult, registry, resultsChan)
	if err != nil {
		p.sendResult(resultsChan, &PipelineResult{
			Step:   constants.PipelineStepAnalyzing,
			Status: constants.StatusFailed,
			Error:  err,
		})
		return
	}
	
	// Étape finale: Completion
	p.sendResult(resultsChan, &PipelineResult{
		Step:   constants.PipelineStepCompleted,
		Status: constants.StatusCompleted,
		Data:   map[string]interface{}{"message": "Pipeline completed successfully"},
	})
}

// executeCrawlStep exécute l'étape de crawling
func (p *pipelineExecutor) executeCrawlStep(ctx context.Context, request *AuditRequest, resultsChan chan<- *PipelineResult) (*crawler.CrawlResult, error) {
	startTime := time.Now()
	
	// Pour le moment, simuler le crawling (sera intégré avec le vrai crawler plus tard)
	time.Sleep(100 * time.Millisecond) // Simulation du temps de crawling
	
	// Vérifier l'annulation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	
	// Créer un résultat simulé
	crawlResult := &crawler.CrawlResult{
		Pages: []crawler.PageData{
			{
				URL:     request.SeedURL,
				Title:   "Test Page Title",
				Content: "Test content for SEO analysis with various keywords and technical elements",
			},
		},
	}
	
	p.sendResult(resultsChan, &PipelineResult{
		Step:     constants.PipelineStepCrawling,
		Status:   constants.StatusCompleted,
		Data:     map[string]interface{}{"pages_crawled": len(crawlResult.Pages)},
		Duration: time.Since(startTime),
	})
	
	return crawlResult, nil
}

// executeAnalysisStep exécute l'étape d'analyse avec les agents
func (p *pipelineExecutor) executeAnalysisStep(ctx context.Context, request *AuditRequest, crawlResult *crawler.CrawlResult, registry AgentRegistry, resultsChan chan<- *PipelineResult) error {
	agents := registry.List()
	
	if len(agents) == 0 {
		return fmt.Errorf("no agents available for analysis")
	}
	
	// Exécuter les agents en parallèle
	var wg sync.WaitGroup
	errors := make(chan error, len(agents))
	
	for name, agent := range agents {
		wg.Add(1)
		go func(agentName string, ag interface{}) {
			defer wg.Done()
			
			startTime := time.Now()
			
			// Vérifier l'annulation
			select {
			case <-ctx.Done():
				errors <- ctx.Err()
				return
			default:
			}
			
			// Pour le MVP, simuler l'exécution des agents
			time.Sleep(50 * time.Millisecond)
			
			p.sendResult(resultsChan, &PipelineResult{
				Step:      constants.PipelineStepAnalyzing,
				AgentName: agentName,
				Status:    constants.StatusCompleted,
				Data:      map[string]interface{}{"agent": agentName, "status": "simulated"},
				Duration:  time.Since(startTime),
			})
			
			errors <- nil
		}(name, agent)
	}
	
	wg.Wait()
	close(errors)
	
	// Vérifier les erreurs
	for err := range errors {
		if err != nil {
			return err
		}
	}
	
	return nil
}

// sendResult envoie un résultat dans le channel
func (p *pipelineExecutor) sendResult(resultsChan chan<- *PipelineResult, result *PipelineResult) {
	select {
	case resultsChan <- result:
	default:
		// Channel plein, ignorer pour éviter le blocage
	}
}

// cleanupExecution nettoie les données d'une exécution terminée
func (p *pipelineExecutor) cleanupExecution(auditID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	delete(p.runningAudits, auditID)
	// Garder la progression pour consultation ultérieure
}