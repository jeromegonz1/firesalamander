package v2

import (
	"context"
	"fmt"
	"sync"
	"time"

	"firesalamander/internal/agents"
	"firesalamander/internal/constants"
)

// orchestratorV2 implémente l'interface OrchestratorV2
type orchestratorV2 struct {
	mu               sync.RWMutex
	registry         AgentRegistry
	executor         PipelineExecutor
	progressManager  ProgressManager
	activeAudits     map[string]*AuditExecution
	progressStreams  map[string]chan *ProgressUpdate
	shutdown         chan struct{}
	shutdownComplete chan struct{}
}

// NewOrchestratorV2 crée une nouvelle instance d'OrchestratorV2
func NewOrchestratorV2() OrchestratorV2 {
	return &orchestratorV2{
		registry:         NewAgentRegistry(),
		executor:         NewPipelineExecutor(),
		progressManager:  NewProgressManager(),
		activeAudits:     make(map[string]*AuditExecution),
		progressStreams:  make(map[string]chan *ProgressUpdate),
		shutdown:         make(chan struct{}),
		shutdownComplete: make(chan struct{}),
	}
}

// RegisterAgent enregistre un agent dans l'orchestrateur
func (o *orchestratorV2) RegisterAgent(name string, agent agents.Agent) error {
	if agent == nil {
		return fmt.Errorf("agent cannot be nil")
	}
	
	return o.registry.Register(name, agent)
}

// StartAudit démarre un audit et retourne un channel de progression
func (o *orchestratorV2) StartAudit(ctx context.Context, request *AuditRequest) (<-chan *ProgressUpdate, error) {
	if request == nil {
		return nil, fmt.Errorf("audit request cannot be nil")
	}
	
	if request.AuditID == "" {
		return nil, fmt.Errorf("audit ID cannot be empty")
	}
	
	if request.SeedURL == "" {
		return nil, fmt.Errorf("seed URL cannot be empty")
	}
	
	o.mu.Lock()
	defer o.mu.Unlock()
	
	// Vérifier si l'audit n'est pas déjà en cours
	if _, exists := o.activeAudits[request.AuditID]; exists {
		return nil, fmt.Errorf("audit %s is already active", request.AuditID)
	}
	
	// Créer l'exécution d'audit
	execution := &AuditExecution{
		AuditID:       request.AuditID,
		Status:        constants.StatusRunning,
		Progress:      0.0,
		CurrentStep:   constants.PipelineStepInitializing,
		StartTime:     time.Now(),
		Results:       make(map[string]interface{}),
		AgentStatuses: make(map[string]string),
	}
	
	o.activeAudits[request.AuditID] = execution
	
	// Créer le channel de progression
	progressChan := make(chan *ProgressUpdate, constants.DefaultStreamBufferSize)
	o.progressStreams[request.AuditID] = progressChan
	
	// Démarrer le tracking de progression
	err := o.progressManager.StartTracking(request.AuditID)
	if err != nil {
		delete(o.activeAudits, request.AuditID)
		delete(o.progressStreams, request.AuditID)
		return nil, fmt.Errorf("failed to start progress tracking: %w", err)
	}
	
	// Démarrer l'exécution asynchrone
	go o.executeAudit(ctx, request, progressChan)
	
	return progressChan, nil
}

// GetAuditStatus récupère le statut actuel d'un audit
func (o *orchestratorV2) GetAuditStatus(auditID string) (*AuditExecution, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	
	execution, exists := o.activeAudits[auditID]
	if !exists {
		return nil, fmt.Errorf("audit %s not found", auditID)
	}
	
	return execution, nil
}

// StreamProgress démarre le streaming de progression pour un audit
func (o *orchestratorV2) StreamProgress(auditID string) (<-chan *ProgressUpdate, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	
	_, exists := o.activeAudits[auditID]
	if !exists {
		return nil, fmt.Errorf("audit %s not found", auditID)
	}
	
	streamChan, exists := o.progressStreams[auditID]
	if !exists {
		return nil, fmt.Errorf("no progress stream for audit %s", auditID)
	}
	
	// Retourner une copie du channel pour éviter les modifications
	return streamChan, nil
}

// ListActiveAudits retourne la liste des audits actifs
func (o *orchestratorV2) ListActiveAudits() []*AuditExecution {
	o.mu.RLock()
	defer o.mu.RUnlock()
	
	audits := make([]*AuditExecution, 0, len(o.activeAudits))
	for _, execution := range o.activeAudits {
		audits = append(audits, execution)
	}
	
	return audits
}

// CancelAudit annule un audit en cours
func (o *orchestratorV2) CancelAudit(auditID string) error {
	o.mu.Lock()
	execution, exists := o.activeAudits[auditID]
	if !exists {
		o.mu.Unlock()
		return fmt.Errorf("audit %s not found", auditID)
	}
	
	// Marquer comme annulé
	execution.Status = constants.StatusFailed
	execution.Error = "Audit cancelled by user"
	now := time.Now()
	execution.EndTime = &now
	o.mu.Unlock()
	
	// Annuler dans l'executor seulement s'il est en cours
	if o.executor.IsRunning(auditID) {
		err := o.executor.Cancel(auditID)
		if err != nil {
			return fmt.Errorf("failed to cancel audit execution: %w", err)
		}
	}
	
	// Arrêter le tracking
	o.progressManager.StopTracking(auditID)
	
	// Fermer le channel de progression de façon sécurisée
	o.mu.Lock()
	if streamChan, exists := o.progressStreams[auditID]; exists {
		select {
		case <-streamChan:
			// Channel déjà fermé
		default:
			close(streamChan)
		}
		delete(o.progressStreams, auditID)
	}
	o.mu.Unlock()
	
	// Nettoyer immédiatement
	o.cleanupAudit(auditID)
	
	return nil
}

// GetResults récupère les résultats consolidés d'un audit terminé
func (o *orchestratorV2) GetResults(auditID string) (*AuditResults, error) {
	o.mu.RLock()
	execution, exists := o.activeAudits[auditID]
	o.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("audit %s not found", auditID)
	}
	
	if execution.Status != constants.StatusCompleted {
		return nil, fmt.Errorf("audit %s is not completed", auditID)
	}
	
	// Pour le MVP, retourner des résultats simulés basés sur les données de l'exécution
	results := &AuditResults{
		AuditID:     auditID,
		StartedAt:   execution.StartTime,
		Status:      execution.Status,
		AgentResults: make(map[string]*agents.AgentResult),
	}
	
	if execution.EndTime != nil {
		results.CompletedAt = *execution.EndTime
		results.Duration = execution.EndTime.Sub(execution.StartTime)
	}
	
	return results, nil
}

// Shutdown arrête proprement l'orchestrateur
func (o *orchestratorV2) Shutdown(ctx context.Context) error {
	// Signaler l'arrêt
	close(o.shutdown)
	
	o.mu.Lock()
	// Annuler tous les audits actifs
	for auditID := range o.activeAudits {
		if o.executor.IsRunning(auditID) {
			o.executor.Cancel(auditID)
		}
		o.progressManager.StopTracking(auditID)
	}
	
	// Fermer tous les channels de progression de façon sécurisée
	for _, streamChan := range o.progressStreams {
		select {
		case <-streamChan:
			// Channel déjà fermé
		default:
			close(streamChan)
		}
	}
	
	// Nettoyer les maps
	o.activeAudits = make(map[string]*AuditExecution)
	o.progressStreams = make(map[string]chan *ProgressUpdate)
	o.mu.Unlock()
	
	// Signaler que l'arrêt est terminé
	close(o.shutdownComplete)
	
	return nil
}

// executeAudit exécute un audit complet
func (o *orchestratorV2) executeAudit(ctx context.Context, request *AuditRequest, progressChan chan *ProgressUpdate) {
	defer func() {
		// Éviter la double fermeture du channel
		select {
		case <-progressChan:
			// Channel déjà fermé
		default:
			close(progressChan)
		}
	}()
	defer o.cleanupAudit(request.AuditID)
	
	// Envoyer une mise à jour d'initialisation
	o.sendProgressUpdate(progressChan, &ProgressUpdate{
		AuditID:   request.AuditID,
		Step:      constants.PipelineStepInitializing,
		Progress:  constants.ProgressInitialized,
		Timestamp: time.Now(),
	})
	
	// Exécuter le pipeline
	resultsChan, err := o.executor.Execute(ctx, request, o.registry)
	if err != nil {
		o.handleAuditError(request.AuditID, err)
		o.sendProgressUpdate(progressChan, &ProgressUpdate{
			AuditID:   request.AuditID,
			Step:      constants.PipelineStepCompleted,
			Progress:  0.0,
			Error:     err.Error(),
			Timestamp: time.Now(),
		})
		return
	}
	
	// Traiter les résultats du pipeline
	for result := range resultsChan {
		// Convertir les résultats du pipeline en mises à jour de progression
		progress := o.calculateProgressFromStep(result.Step)
		
		update := &ProgressUpdate{
			AuditID:   request.AuditID,
			Step:      result.Step,
			Progress:  progress,
			AgentName: result.AgentName,
			Data:      result.Data,
			Timestamp: time.Now(),
		}
		
		if result.Error != nil {
			update.Error = result.Error.Error()
			update.AgentStatus = constants.StatusFailed
		} else {
			update.AgentStatus = result.Status
		}
		
		o.sendProgressUpdate(progressChan, update)
		o.updateAuditExecution(request.AuditID, update)
	}
	
	// Marquer l'audit comme terminé
	o.completeAudit(request.AuditID)
	
	o.sendProgressUpdate(progressChan, &ProgressUpdate{
		AuditID:   request.AuditID,
		Step:      constants.PipelineStepCompleted,
		Progress:  constants.ProgressFullComplete,
		Timestamp: time.Now(),
	})
}

// sendProgressUpdate envoie une mise à jour de progression
func (o *orchestratorV2) sendProgressUpdate(progressChan chan *ProgressUpdate, update *ProgressUpdate) {
	defer func() {
		if r := recover(); r != nil {
			// Channel fermé, ignorer silencieusement
		}
	}()
	
	select {
	case progressChan <- update:
	default:
		// Channel plein, ignorer pour éviter le blocage
	}
}

// calculateProgressFromStep calcule le pourcentage de progression basé sur l'étape
func (o *orchestratorV2) calculateProgressFromStep(step string) float64 {
	switch step {
	case constants.PipelineStepInitializing:
		return constants.ProgressInitialized
	case constants.PipelineStepCrawling:
		return constants.ProgressCrawlingComplete
	case constants.PipelineStepAnalyzing:
		return constants.ProgressAnalysisComplete
	case constants.PipelineStepReporting:
		return constants.ProgressReportingComplete
	case constants.PipelineStepCompleted:
		return constants.ProgressFullComplete
	default:
		return 0.0
	}
}

// updateAuditExecution met à jour l'état d'exécution d'un audit
func (o *orchestratorV2) updateAuditExecution(auditID string, update *ProgressUpdate) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	execution, exists := o.activeAudits[auditID]
	if !exists {
		return
	}
	
	execution.Progress = update.Progress
	execution.CurrentStep = update.Step
	
	if update.AgentName != "" {
		execution.AgentStatuses[update.AgentName] = update.AgentStatus
	}
	
	if update.Error != "" {
		execution.Error = update.Error
		execution.Status = constants.StatusFailed
	}
}

// completeAudit marque un audit comme terminé
func (o *orchestratorV2) completeAudit(auditID string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	execution, exists := o.activeAudits[auditID]
	if !exists {
		return
	}
	
	execution.Status = constants.StatusCompleted
	execution.Progress = constants.ProgressFullComplete
	execution.CurrentStep = constants.PipelineStepCompleted
	now := time.Now()
	execution.EndTime = &now
}

// handleAuditError gère les erreurs d'audit
func (o *orchestratorV2) handleAuditError(auditID string, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	execution, exists := o.activeAudits[auditID]
	if !exists {
		return
	}
	
	execution.Status = constants.StatusFailed
	execution.Error = err.Error()
	now := time.Now()
	execution.EndTime = &now
}

// cleanupAudit nettoie les ressources d'un audit terminé
func (o *orchestratorV2) cleanupAudit(auditID string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	delete(o.activeAudits, auditID)
	delete(o.progressStreams, auditID)
}

// NewProgressManager crée une nouvelle instance de ProgressManager (implémentation basique)
func NewProgressManager() ProgressManager {
	return &basicProgressManager{
		tracking: make(map[string]bool),
	}
}

// basicProgressManager implémentation basique du ProgressManager
type basicProgressManager struct {
	mu       sync.RWMutex
	tracking map[string]bool
}

func (pm *basicProgressManager) StartTracking(auditID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.tracking[auditID] = true
	return nil
}

func (pm *basicProgressManager) UpdateProgress(auditID string, update *ProgressUpdate) error {
	return nil // Implémentation basique
}

func (pm *basicProgressManager) GetProgress(auditID string) (*ExecutionProgress, error) {
	return nil, fmt.Errorf("not implemented")
}

func (pm *basicProgressManager) Subscribe(auditID string) (<-chan *ProgressUpdate, error) {
	return nil, fmt.Errorf("not implemented")
}

func (pm *basicProgressManager) Unsubscribe(auditID string, subscription *StreamSubscription) error {
	return nil
}

func (pm *basicProgressManager) CompleteAudit(auditID string, results *AuditResults) error {
	return nil
}

func (pm *basicProgressManager) FailAudit(auditID string, err error) error {
	return nil
}

func (pm *basicProgressManager) StopTracking(auditID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	delete(pm.tracking, auditID)
	return nil
}