package integration

import (
	"context"
	"fmt"
	"time"
)

// ErrorHandler manages error recovery and fallback strategies
type ErrorHandler struct {
	MaxRetries      int
	RetryDelay      time.Duration
	FallbackEnabled bool
}

// ErrorType categorizes different types of errors
type ErrorType string

const (
	ErrorTypeCrawler   ErrorType = "crawler"
	ErrorTypeTechnical ErrorType = "technical"
	ErrorTypeSemantic  ErrorType = "semantic"
	ErrorTypeReport    ErrorType = "report"
	ErrorTypeTimeout   ErrorType = "timeout"
	ErrorTypeNetwork   ErrorType = "network"
)

// RecoveryAction defines what to do when an error occurs
type RecoveryAction struct {
	Type        ErrorType `json:"type"`
	Action      string    `json:"action"`
	Fallback    bool      `json:"fallback"`
	RetryCount  int       `json:"retry_count"`
	MaxRetries  int       `json:"max_retries"`
	LastError   string    `json:"last_error"`
	NextRetryAt time.Time `json:"next_retry_at"`
}

// NewErrorHandler creates a new error handler with fallback strategies
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		MaxRetries:      3,
		RetryDelay:      2 * time.Second,
		FallbackEnabled: true,
	}
}

// HandleCrawlerError implements fallback for crawler failures
func (eh *ErrorHandler) HandleCrawlerError(ctx context.Context, err error, auditID string) (*RecoveryAction, error) {
	action := &RecoveryAction{
		Type:       ErrorTypeCrawler,
		LastError:  err.Error(),
		MaxRetries: eh.MaxRetries,
	}

	// Classify error type
	if ctx.Err() == context.DeadlineExceeded {
		action.Type = ErrorTypeTimeout
		action.Action = "timeout_fallback"
		action.Fallback = true
		return action, nil
	}

	if isNetworkError(err) {
		action.Type = ErrorTypeNetwork
		action.Action = "network_retry"
		action.Fallback = false
		action.NextRetryAt = time.Now().Add(eh.RetryDelay)
		return action, nil
	}

	// Default: mark as failed but continue pipeline
	action.Action = "continue_without_crawl"
	action.Fallback = true
	return action, nil
}

// HandleSemanticError implements fallback for semantic analysis failures
func (eh *ErrorHandler) HandleSemanticError(ctx context.Context, err error, auditID string) (*RecoveryAction, error) {
	action := &RecoveryAction{
		Type:       ErrorTypeSemantic,
		LastError:  err.Error(),
		MaxRetries: eh.MaxRetries,
		Action:     "skip_semantic",
		Fallback:   true,
	}

	// Semantic analysis is optional - always fallback to basic keywords
	return action, nil
}

// HandleTechnicalError implements fallback for technical analysis failures
func (eh *ErrorHandler) HandleTechnicalError(ctx context.Context, err error, auditID string) (*RecoveryAction, error) {
	action := &RecoveryAction{
		Type:       ErrorTypeTechnical,
		LastError:  err.Error(),
		MaxRetries: eh.MaxRetries,
	}

	// Technical analysis is core - try basic validation
	action.Action = "basic_technical_analysis"
	action.Fallback = true
	return action, nil
}

// HandleReportError implements fallback for report generation failures
func (eh *ErrorHandler) HandleReportError(ctx context.Context, err error, auditID string) (*RecoveryAction, error) {
	action := &RecoveryAction{
		Type:       ErrorTypeReport,
		LastError:  err.Error(),
		MaxRetries: eh.MaxRetries,
		Action:     "generate_basic_report",
		Fallback:   true,
	}

	return action, nil
}

// ShouldRetry determines if an error warrants a retry
func (eh *ErrorHandler) ShouldRetry(action *RecoveryAction) bool {
	if action.RetryCount >= action.MaxRetries {
		return false
	}

	// Network errors can be retried
	if action.Type == ErrorTypeNetwork {
		return true
	}

	// Timeouts generally shouldn't be retried
	if action.Type == ErrorTypeTimeout {
		return false
	}

	return false
}

// ExecuteRecovery performs the recovery action
func (eh *ErrorHandler) ExecuteRecovery(ctx context.Context, action *RecoveryAction, pipeline *Pipeline, execution *AuditExecution) error {
	switch action.Action {
	case "timeout_fallback":
		return eh.timeoutFallback(execution)
	case "skip_semantic":
		return eh.skipSemantic(execution)
	case "basic_technical_analysis":
		return eh.basicTechnicalAnalysis(execution)
	case "generate_basic_report":
		return eh.generateBasicReport(execution)
	case "continue_without_crawl":
		return eh.continueWithoutCrawl(execution)
	default:
		return fmt.Errorf("unknown recovery action: %s", action.Action)
	}
}

// Fallback implementations
func (eh *ErrorHandler) timeoutFallback(execution *AuditExecution) error {
	execution.Status = "partial"
	execution.CurrentStep = "timeout_recovery"
	execution.Error = "Timeout - résultats partiels disponibles"
	return nil
}

func (eh *ErrorHandler) skipSemantic(execution *AuditExecution) error {
	// Create minimal semantic results
	execution.Results["semantic"] = map[string]interface{}{
		"audit_id":   execution.AuditID,
		"status":     "skipped",
		"keywords":   []string{},
		"topics":     []string{},
		"error":      "Service sémantique indisponible",
		"fallback":   true,
	}
	return nil
}

func (eh *ErrorHandler) basicTechnicalAnalysis(execution *AuditExecution) error {
	// Create minimal technical results
	execution.Results["technical"] = map[string]interface{}{
		"audit_id": execution.AuditID,
		"status":   "basic",
		"findings": []map[string]interface{}{
			{
				"id":       "fallback-check",
				"severity": "info",
				"message":  "Analyse technique de base effectuée",
			},
		},
		"fallback": true,
	}
	return nil
}

func (eh *ErrorHandler) generateBasicReport(execution *AuditExecution) error {
	// Create minimal report
	execution.Results["html_report"] = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><title>Rapport Audit %s</title></head>
<body>
<h1>Rapport d'audit basique - %s</h1>
<p>Audit partiellement complété avec fallbacks.</p>
<p>Audit ID: %s</p>
<p>Statut: %s</p>
</body>
</html>`, execution.AuditID, execution.AuditID, execution.AuditID, execution.Status)
	return nil
}

func (eh *ErrorHandler) continueWithoutCrawl(execution *AuditExecution) error {
	// Create empty crawl results
	execution.Results["crawl"] = map[string]interface{}{
		"audit_id":      execution.AuditID,
		"pages_crawled": 0,
		"status":        "failed",
		"error":         "Crawl impossible",
	}
	return nil
}

// Helper functions
func isNetworkError(err error) bool {
	errStr := err.Error()
	return containsAny(errStr, []string{
		"connection refused",
		"timeout",
		"network unreachable",
		"no such host",
		"certificate",
	})
}

func containsAny(str string, substrings []string) bool {
	for _, substr := range substrings {
		if contains(str, substr) {
			return true
		}
	}
	return false
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr ||
		   len(str) > len(substr) && containsRecursive(str[1:], substr)
}

func containsRecursive(str, substr string) bool {
	if len(str) < len(substr) {
		return false
	}
	if str[:len(substr)] == substr {
		return true
	}
	return containsRecursive(str[1:], substr)
}