package v2

import (
	"time"

	"firesalamander/internal/agents"
)

// AuditRequest représente une demande d'audit pour l'orchestrateur V2
type AuditRequest struct {
	AuditID   string                 `json:"audit_id"`
	SeedURL   string                 `json:"seed_url"`
	MaxPages  int                    `json:"max_pages"`
	Options   map[string]interface{} `json:"options"`
	Timestamp time.Time              `json:"timestamp"`
}

// ProgressUpdate représente une mise à jour de progression streamée
type ProgressUpdate struct {
	AuditID     string                 `json:"audit_id"`
	Step        string                 `json:"step"`
	Progress    float64                `json:"progress"`
	AgentName   string                 `json:"agent_name,omitempty"`
	AgentStatus string                 `json:"agent_status,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// AuditExecution représente l'état d'exécution d'un audit
type AuditExecution struct {
	AuditID       string                 `json:"audit_id"`
	Status        string                 `json:"status"`
	Progress      float64                `json:"progress"`
	CurrentStep   string                 `json:"current_step"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       *time.Time             `json:"end_time,omitempty"`
	Error         string                 `json:"error,omitempty"`
	Results       map[string]interface{} `json:"results"`
	AgentStatuses map[string]string      `json:"agent_statuses"`
}

// PipelineResult représente le résultat d'exécution du pipeline
type PipelineResult struct {
	Step      string                 `json:"step"`
	AgentName string                 `json:"agent_name,omitempty"`
	Status    string                 `json:"status"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     error                  `json:"-"`
	Duration  time.Duration          `json:"duration"`
}

// ExecutionProgress représente la progression détaillée d'une exécution
type ExecutionProgress struct {
	AuditID          string            `json:"audit_id"`
	OverallProgress  float64           `json:"overall_progress"`
	CurrentStep      string            `json:"current_step"`
	AgentProgresses  map[string]float64 `json:"agent_progresses"`
	CompletedAgents  []string          `json:"completed_agents"`
	RunningAgents    []string          `json:"running_agents"`
	PendingAgents    []string          `json:"pending_agents"`
	FailedAgents     []string          `json:"failed_agents"`
}

// AgentExecution représente l'état d'exécution d'un agent
type AgentExecution struct {
	Name      string                 `json:"name"`
	Status    string                 `json:"status"`
	StartTime *time.Time             `json:"start_time,omitempty"`
	EndTime   *time.Time             `json:"end_time,omitempty"`
	Result    *agents.AgentResult    `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// AuditResults représente les résultats consolidés d'un audit complet
type AuditResults struct {
	AuditID       string                   `json:"audit_id"`
	SiteURL       string                   `json:"site_url"`
	StartedAt     time.Time                `json:"started_at"`
	CompletedAt   time.Time                `json:"completed_at"`
	Duration      time.Duration            `json:"duration"`
	Status        string                   `json:"status"`
	AgentResults  map[string]*agents.AgentResult `json:"agent_results"`
	Summary       *AuditSummary            `json:"summary"`
	Errors        []string                 `json:"errors,omitempty"`
}

// AuditSummary représente un résumé consolidé des résultats
type AuditSummary struct {
	TotalPages       int     `json:"total_pages"`
	TotalKeywords    int     `json:"total_keywords"`
	BrokenLinksCount int     `json:"broken_links_count"`
	AverageSeoScore  float64 `json:"average_seo_score"`
	IssuesCount      int     `json:"issues_count"`
	Recommendations  []string `json:"recommendations"`
}

// StreamSubscription représente un abonnement au streaming de progression
type StreamSubscription struct {
	AuditID   string
	Channel   chan *ProgressUpdate
	Active    bool
	StartTime time.Time
}

// RegistryStats représente les statistiques du registry d'agents
type RegistryStats struct {
	TotalAgents     int               `json:"total_agents"`
	HealthyAgents   int               `json:"healthy_agents"`
	UnhealthyAgents int               `json:"unhealthy_agents"`
	AgentsList      []string          `json:"agents_list"`
	HealthStatus    map[string]string `json:"health_status"`
}