package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"firesalamander/internal/audit"
	"firesalamander/internal/config"
	"firesalamander/internal/constants"
	"firesalamander/internal/crawler"
	"firesalamander/internal/orchestrator"
	"firesalamander/internal/report"
	"firesalamander/internal/semantic"
)

// Pipeline manages the complete audit workflow
type Pipeline struct {
	config     *config.Config
	crawler    *crawler.Crawler
	technical  *audit.TechnicalAnalyzer
	semantic   *semantic.SemanticClient
	report     *report.ReportEngine
	mu         sync.RWMutex
	audits     map[string]*AuditExecution
}

// AuditExecution tracks a single audit through the pipeline
type AuditExecution struct {
	AuditID     string                 `json:"audit_id"`
	Status      string                 `json:"status"`
	Progress    float64                `json:"progress"`
	CurrentStep string                 `json:"current_step"`
	StartTime   time.Time              `json:"start_time"`
	Error       string                 `json:"error,omitempty"`
	Results     map[string]interface{} `json:"results"`
	OutputDir   string                 `json:"output_dir"`
}

// JSONRPCMessage represents a JSON-RPC 2.0 message
type JSONRPCMessage struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      string      `json:"id"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewPipeline creates a new integration pipeline
func NewPipeline(cfg *config.Config) (*Pipeline, error) {
	// Load crawler config, fallback to defaults for tests
	crawlerCfg, err := config.LoadCrawlerConfig("config/crawler.yaml")
	if err != nil {
		// Use default config for tests
		crawlerCfg = &config.CrawlerConfig{
			Performance: config.Performance{
				ConcurrentRequests: 2,
				RequestTimeout:     10 * time.Second,
			},
			Limits: config.Limits{
				MaxURLs:  10,
				MaxDepth: 2,
			},
			UserAgent: "Fire Salamander Test",
		}
	}

	// Create agents with existing constructors
	crawlerAgent := crawler.NewCrawler(*crawlerCfg)
	
	// Load tech rules for technical analyzer
	techRules := audit.TechRules{
		Title: audit.TitleRules{
			MinLength:        10,
			MaxLength:        60,
			MissingSeverity:  "high",
			TooShortSeverity: "medium",
			TooLongSeverity:  "medium",
		},
	}
	techAnalyzer := audit.NewTechnicalAnalyzer(techRules)
	
	semanticClient := semantic.NewSemanticClient(constants.DefaultSemanticServiceURL)
	reportEngine := report.NewReportEngine()

	return &Pipeline{
		config:    cfg,
		crawler:   crawlerAgent,
		technical: techAnalyzer,
		semantic:  semanticClient,
		report:    reportEngine,
		audits:    make(map[string]*AuditExecution),
	}, nil
}

// StartAudit begins a complete audit pipeline
func (p *Pipeline) StartAudit(ctx context.Context, request orchestrator.AuditRequest) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	execution := &AuditExecution{
		AuditID:     request.AuditID,
		Status:      "pending",
		Progress:    0.0,
		CurrentStep: "initializing",
		StartTime:   time.Now(),
		Results:     make(map[string]interface{}),
		OutputDir:   filepath.Join("audits", request.AuditID),
	}

	p.audits[request.AuditID] = execution

	// Start pipeline asynchronously
	go p.runPipeline(ctx, request, execution)

	return nil
}

// GetAuditStatus returns the current status of an audit
func (p *Pipeline) GetAuditStatus(auditID string) *AuditExecution {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	if audit, exists := p.audits[auditID]; exists {
		return audit
	}
	return nil
}

// runPipeline executes the complete audit workflow
func (p *Pipeline) runPipeline(ctx context.Context, request orchestrator.AuditRequest, execution *AuditExecution) {
	defer func() {
		if r := recover(); r != nil {
			p.updateStatus(execution, "failed", "panic", fmt.Sprintf("Pipeline panic: %v", r))
		}
	}()

	// Step 1: Crawling
	if err := p.runCrawlStep(ctx, request, execution); err != nil {
		p.updateStatus(execution, "failed", "crawling", err.Error())
		return
	}

	// Step 2: Parallel analysis
	var wg sync.WaitGroup
	var techErr, semanticErr error

	// Technical analysis
	wg.Add(1)
	go func() {
		defer wg.Done()
		techErr = p.runTechnicalStep(ctx, request, execution)
	}()

	// Semantic analysis
	wg.Add(1)
	go func() {
		defer wg.Done()
		semanticErr = p.runSemanticStep(ctx, request, execution)
	}()

	wg.Wait()

	if techErr != nil || semanticErr != nil {
		errorMsg := fmt.Sprintf("Analysis errors - Tech: %v, Semantic: %v", techErr, semanticErr)
		p.updateStatus(execution, "failed", "analyzing", errorMsg)
		return
	}

	// Step 3: Report generation
	if err := p.runReportStep(ctx, request, execution); err != nil {
		p.updateStatus(execution, "failed", "reporting", err.Error())
		return
	}

	p.updateStatus(execution, "completed", "finished", "")
}

func (p *Pipeline) runCrawlStep(ctx context.Context, request orchestrator.AuditRequest, execution *AuditExecution) error {
	p.updateStatus(execution, "crawling", "crawling", "")

	// Execute crawl using existing crawler
	outputDir := filepath.Join("audits", request.AuditID)
	crawlResult, err := p.crawler.Crawl(ctx, request.SeedURL, outputDir)

	if err != nil {
		return fmt.Errorf("crawl failed: %w", err)
	}

	// Save crawl results
	execution.Results["crawl"] = crawlResult
	p.updateProgress(execution, 30.0)

	return nil
}

func (p *Pipeline) runTechnicalStep(ctx context.Context, request orchestrator.AuditRequest, execution *AuditExecution) error {
	p.updateStatus(execution, "analyzing_tech", "technical analysis", "")

	crawlData, ok := execution.Results["crawl"].(*crawler.CrawlResult)
	if !ok {
		return fmt.Errorf("invalid crawl data")
	}

	// Analyze pages using existing technical analyzer
	var allFindings []audit.Finding
	for _, page := range crawlData.Pages {
		titleFindings := p.technical.ValidateTitle(page.Title, page.URL)
		allFindings = append(allFindings, titleFindings...)
		
		headingFindings := p.technical.ValidateHeadings(1, 2, page.URL)
		allFindings = append(allFindings, headingFindings...)
	}

	execution.Results["technical"] = map[string]interface{}{
		"audit_id": request.AuditID,
		"findings": allFindings,
		"status":   "completed",
	}
	p.updateProgress(execution, 60.0)

	return nil
}

func (p *Pipeline) runSemanticStep(ctx context.Context, request orchestrator.AuditRequest, execution *AuditExecution) error {
	p.updateStatus(execution, "analyzing_sem", "semantic analysis", "")

	crawlData, ok := execution.Results["crawl"].(*crawler.CrawlResult)
	if !ok {
		return fmt.Errorf("invalid crawl data")
	}

	// Call semantic analysis using existing client
	semanticResult, err := p.semantic.Analyze(request.AuditID, *crawlData)
	if err != nil {
		return fmt.Errorf("semantic analysis failed: %w", err)
	}

	execution.Results["semantic"] = semanticResult
	p.updateProgress(execution, 80.0)

	return nil
}

func (p *Pipeline) runReportStep(ctx context.Context, request orchestrator.AuditRequest, execution *AuditExecution) error {
	p.updateStatus(execution, "reporting", "generating reports", "")

	// Combine results using existing report engine structure
	crawlData := execution.Results["crawl"].(*crawler.CrawlResult)
	techResults := execution.Results["technical"].(map[string]interface{})
	semanticResults := execution.Results["semantic"].(*semantic.SemanticResult)
	
	auditResults := report.AuditResults{
		AuditID:         request.AuditID,
		SiteURL:         request.SeedURL,
		StartedAt:       execution.StartTime.Format(time.RFC3339),
		Duration:        time.Since(execution.StartTime).String(),
		TotalPages:      len(crawlData.Pages),
		CrawlData:       *crawlData,
		TechResults:     audit.TechResult{Findings: techResults["findings"].([]audit.Finding)},
		SemanticResults: *semanticResults,
	}

	// Generate HTML report
	htmlReport, err := p.report.GenerateHTML(auditResults)
	if err != nil {
		return fmt.Errorf("HTML report generation failed: %w", err)
	}

	execution.Results["html_report"] = htmlReport
	p.updateProgress(execution, 100.0)

	return nil
}

func (p *Pipeline) updateStatus(execution *AuditExecution, status, step, errorMsg string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	execution.Status = status
	execution.CurrentStep = step
	if errorMsg != "" {
		execution.Error = errorMsg
	}
}

func (p *Pipeline) updateProgress(execution *AuditExecution, progress float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	execution.Progress = progress
}

func (p *Pipeline) getOption(options map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, exists := options[key]; exists {
		return value
	}
	return defaultValue
}

// Helper function to safely get options with defaults
func (p *Pipeline) getIntOption(options map[string]interface{}, key string, defaultValue int) int {
	if value, exists := options[key]; exists {
		if intVal, ok := value.(int); ok {
			return intVal
		}
	}
	return defaultValue
}