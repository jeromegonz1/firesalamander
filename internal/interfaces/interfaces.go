package interfaces

import (
	"context"
	"firesalamander/internal/audit"
	"firesalamander/internal/crawler"
	"firesalamander/internal/semantic"
)

// PageCrawler defines the interface for crawling web pages
type PageCrawler interface {
	Crawl(ctx context.Context, seedURL string, outputDir string) (*crawler.CrawlResult, error)
}

// TechnicalAnalyzer defines the interface for technical analysis
type TechnicalAnalyzer interface {
	Analyze(crawlResult *crawler.CrawlResult, auditID string) (*audit.TechResult, error)
}

// SemanticAnalyzer defines the interface for semantic analysis  
type SemanticAnalyzer interface {
	Analyze(auditID string, crawlData crawler.CrawlResult) (*semantic.SemanticResult, error)
}

// ReportGenerator defines the interface for generating reports
type ReportGenerator interface {
	Generate(results AuditResults, format string) (string, error)
}

// AuditResults represents the combined results of all analyses
type AuditResults struct {
	CrawlData       *crawler.CrawlResult
	TechResults     *audit.TechResult  
	SemanticResults *semantic.SemanticResult
}

// Orchestrator defines the interface for audit orchestration
type Orchestrator interface {
	StartAudit(request AuditRequest) error
	GetStatus(auditID string) (interface{}, error)
}

// AuditRequest represents a request for audit
type AuditRequest struct {
	SeedURL string
	AuditID string
	Options map[string]interface{}
}