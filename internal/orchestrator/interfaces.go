package v2

import (
	"context"

	"firesalamander/internal/agents"
)

// OrchestratorV2 définit l'interface principale pour l'orchestration des agents
type OrchestratorV2 interface {
	// RegisterAgent enregistre un agent dans l'orchestrateur
	RegisterAgent(name string, agent agents.Agent) error
	
	// StartAudit démarre un audit et retourne un channel de progression
	StartAudit(ctx context.Context, request *AuditRequest) (<-chan *ProgressUpdate, error)
	
	// GetAuditStatus récupère le statut actuel d'un audit
	GetAuditStatus(auditID string) (*AuditExecution, error)
	
	// StreamProgress démarre le streaming de progression pour un audit
	StreamProgress(auditID string) (<-chan *ProgressUpdate, error)
	
	// ListActiveAudits retourne la liste des audits actifs
	ListActiveAudits() []*AuditExecution
	
	// CancelAudit annule un audit en cours
	CancelAudit(auditID string) error
	
	// GetResults récupère les résultats consolidés d'un audit terminé
	GetResults(auditID string) (*AuditResults, error)
	
	// Shutdown arrête proprement l'orchestrateur
	Shutdown(ctx context.Context) error
}

// AgentRegistry définit l'interface pour la gestion des agents
type AgentRegistry interface {
	// Register enregistre un agent dans le registry
	Register(name string, agent agents.Agent) error
	
	// Get récupère un agent par son nom
	Get(name string) (agents.Agent, bool)
	
	// List retourne tous les agents enregistrés
	List() map[string]agents.Agent
	
	// HealthCheckAll vérifie la santé de tous les agents
	HealthCheckAll() map[string]error
	
	// Unregister supprime un agent du registry
	Unregister(name string) error
	
	// Count retourne le nombre d'agents enregistrés
	Count() int
	
	// GetStats retourne les statistiques du registry
	GetStats() *RegistryStats
}

// PipelineExecutor définit l'interface pour l'exécution du pipeline d'audit
type PipelineExecutor interface {
	// Execute lance l'exécution du pipeline pour un audit
	Execute(ctx context.Context, request *AuditRequest, registry AgentRegistry) (<-chan *PipelineResult, error)
	
	// GetProgress récupère la progression actuelle d'un audit
	GetProgress(auditID string) (*ExecutionProgress, error)
	
	// Cancel annule l'exécution d'un audit
	Cancel(auditID string) error
	
	// IsRunning vérifie si un audit est en cours d'exécution
	IsRunning(auditID string) bool
	
	// GetRunningAudits retourne la liste des audits en cours
	GetRunningAudits() []string
}

// ProgressManager définit l'interface pour la gestion de la progression
type ProgressManager interface {
	// StartTracking commence le suivi de progression pour un audit
	StartTracking(auditID string) error
	
	// UpdateProgress met à jour la progression d'un audit
	UpdateProgress(auditID string, update *ProgressUpdate) error
	
	// GetProgress récupère la progression actuelle
	GetProgress(auditID string) (*ExecutionProgress, error)
	
	// Subscribe s'abonne aux mises à jour de progression
	Subscribe(auditID string) (<-chan *ProgressUpdate, error)
	
	// Unsubscribe se désabonne des mises à jour
	Unsubscribe(auditID string, subscription *StreamSubscription) error
	
	// CompleteAudit marque un audit comme terminé
	CompleteAudit(auditID string, results *AuditResults) error
	
	// FailAudit marque un audit comme échoué
	FailAudit(auditID string, err error) error
	
	// StopTracking arrête le suivi d'un audit
	StopTracking(auditID string) error
}

// ResultsAggregator définit l'interface pour l'agrégation des résultats
type ResultsAggregator interface {
	// Aggregate combine les résultats de tous les agents
	Aggregate(auditID string, agentResults map[string]*agents.AgentResult) (*AuditResults, error)
	
	// GenerateSummary génère un résumé des résultats
	GenerateSummary(results *AuditResults) (*AuditSummary, error)
	
	// ExtractRecommendations extrait les recommandations des résultats
	ExtractRecommendations(results *AuditResults) ([]string, error)
}

// AuditStorage définit l'interface pour la persistance des audits
type AuditStorage interface {
	// Store sauvegarde les données d'un audit
	Store(auditID string, data interface{}) error
	
	// Retrieve récupère les données d'un audit
	Retrieve(auditID string, target interface{}) error
	
	// Delete supprime les données d'un audit
	Delete(auditID string) error
	
	// List retourne la liste des audits stockés
	List() ([]string, error)
	
	// Exists vérifie si un audit existe dans le storage
	Exists(auditID string) bool
}