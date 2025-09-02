package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"firesalamander/internal/audit"
	"firesalamander/internal/config"
	"firesalamander/internal/agents/crawler"
)

type Orchestrator struct {
	CrawlerConfig *config.CrawlerConfig
	Status        map[string]*AuditStatus
	mutex         sync.RWMutex
}

type AuditStatus struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"` // pending, crawling, analyzing, reporting, completed, failed
	Progress    int       `json:"progress"`
	CurrentStep string    `json:"current_step"`
	Error       string    `json:"error,omitempty"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

type AuditRequest struct {
	SeedURL   string                 `json:"seed_url"`
	AuditID   string                 `json:"audit_id"`
	Options   map[string]interface{} `json:"options"`
}

type AuditResult struct {
	AuditID    string      `json:"audit_id"`
	Status     string      `json:"status"`
	CrawlData  interface{} `json:"crawl_data,omitempty"`
	TechData   interface{} `json:"tech_data,omitempty"`
	ReportPath string      `json:"report_path,omitempty"`
	Error      string      `json:"error,omitempty"`
}

func NewOrchestrator() (*Orchestrator, error) {
	crawlerConfig, err := config.LoadCrawlerConfig("config/crawler.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load crawler config: %w", err)
	}

	return &Orchestrator{
		CrawlerConfig: crawlerConfig,
		Status:        make(map[string]*AuditStatus),
	}, nil
}

func (o *Orchestrator) StartAudit(req AuditRequest) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Check if audit already exists
	if _, exists := o.Status[req.AuditID]; exists {
		return fmt.Errorf("audit %s already exists", req.AuditID)
	}

	// Create audit status
	o.Status[req.AuditID] = &AuditStatus{
		ID:          req.AuditID,
		Status:      "pending",
		Progress:    0,
		CurrentStep: "initializing",
		StartedAt:   time.Now(),
	}

	// Start audit in background
	go o.runAudit(req)

	return nil
}

func (o *Orchestrator) GetStatus(auditID string) (*AuditStatus, error) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	status, exists := o.Status[auditID]
	if !exists {
		return nil, fmt.Errorf("audit %s not found", auditID)
	}

	return status, nil
}

func (o *Orchestrator) runAudit(req AuditRequest) {
	auditID := req.AuditID
	
	// Update status helper
	updateStatus := func(status, step string, progress int, err error) {
		o.mutex.Lock()
		defer o.mutex.Unlock()
		
		if auditStatus, exists := o.Status[auditID]; exists {
			auditStatus.Status = status
			auditStatus.CurrentStep = step
			auditStatus.Progress = progress
			if err != nil {
				auditStatus.Error = err.Error()
				auditStatus.Status = "failed"
			}
			if status == "completed" {
				auditStatus.CompletedAt = time.Now()
			}
		}
	}

	// Create audit directory
	auditDir := filepath.Join("audits", auditID)
	if err := os.MkdirAll(auditDir, 0755); err != nil {
		updateStatus("failed", "setup", 0, err)
		return
	}

	// Step 1: Crawling
	updateStatus("crawling", "crawler", 10, nil)
	log.Printf("Starting crawl for audit %s", auditID)

	crawlerInstance := crawler.NewCrawler(*o.CrawlerConfig)
	ctx := context.Background()
	
	crawlResult, err := crawlerInstance.Crawl(ctx, req.SeedURL, auditDir)
	if err != nil {
		updateStatus("failed", "crawler", 10, err)
		return
	}

	log.Printf("Crawl completed: %d pages found", len(crawlResult.Pages))

	// Step 2: Technical Analysis
	updateStatus("analyzing", "technical_audit", 50, nil)
	log.Printf("Starting technical analysis for audit %s", auditID)

	// Load tech rules
	techRules := audit.TechRules{} // TODO: Load from config
	analyzer := audit.NewTechnicalAnalyzer(techRules)

	techResult, err := analyzer.Analyze(crawlResult, auditID)
	if err != nil {
		updateStatus("failed", "technical_audit", 50, err)
		return
	}

	// Save tech results
	techFile := filepath.Join(auditDir, "tech_result.json")
	if err := saveJSON(techResult, techFile); err != nil {
		updateStatus("failed", "saving_tech", 60, err)
		return
	}

	log.Printf("Technical analysis completed")

	// Step 3: Report Generation (simplified for now)
	updateStatus("reporting", "generate_report", 80, nil)
	
	// Create simple report
	report := map[string]interface{}{
		"audit_id":     auditID,
		"pages_count":  len(crawlResult.Pages),
		"tech_score":   techResult.Scores,
		"findings":     len(techResult.Findings),
	}

	reportFile := filepath.Join(auditDir, "report.json")
	if err := saveJSON(report, reportFile); err != nil {
		updateStatus("failed", "saving_report", 90, err)
		return
	}

	// Complete
	updateStatus("completed", "finished", 100, nil)
	log.Printf("Audit %s completed successfully", auditID)
}

func saveJSON(data interface{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}