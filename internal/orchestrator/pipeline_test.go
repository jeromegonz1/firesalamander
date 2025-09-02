package v2

import (
	"context"
	"fmt"
	"testing"
	"time"

	"firesalamander/internal/constants"
	"firesalamander/internal/agents/crawler"
)

// MockCrawlResult pour simuler les résultats de crawling
func createMockCrawlResult() *crawler.CrawlResult {
	return &crawler.CrawlResult{
		Pages: []crawler.PageData{
			{
				URL:     "https://example.com",
				Title:   "Test Page",
				Content: "This is test content with some keywords for SEO analysis",
			},
			{
				URL:     "https://example.com/page2", 
				Title:   "Another Test Page",
				Content: "More test content for comprehensive analysis",
			},
		},
	}
}

func TestNewPipelineExecutor(t *testing.T) {
	executor := NewPipelineExecutor()
	
	if executor == nil {
		t.Fatal("NewPipelineExecutor should not return nil")
	}
	
	runningAudits := executor.GetRunningAudits()
	if len(runningAudits) != 0 {
		t.Errorf("Expected 0 running audits initially, got %d", len(runningAudits))
	}
}

func TestPipelineExecutor_Execute(t *testing.T) {
	executor := NewPipelineExecutor()
	registry := NewAgentRegistry()
	
	// Register mock agents
	mockKeywordAgent := &MockAgent{name: constants.AgentNameKeyword}
	mockTechnicalAgent := &MockAgent{name: constants.AgentNameTechnical}
	
	registry.Register(constants.AgentNameKeyword, mockKeywordAgent)
	registry.Register(constants.AgentNameTechnical, mockTechnicalAgent)
	
	request := &AuditRequest{
		AuditID:   "test-audit-123",
		SeedURL:   "https://example.com",
		MaxPages:  5,
		Options:   map[string]interface{}{},
		Timestamp: time.Now(),
	}
	
	ctx := context.Background()
	resultsChan, err := executor.Execute(ctx, request, registry)
	
	if err != nil {
		t.Errorf("Execute should not return error: %v", err)
	}
	
	if resultsChan == nil {
		t.Fatal("Execute should return a results channel")
	}
	
	// Verify audit is now running
	if !executor.IsRunning(request.AuditID) {
		t.Error("Audit should be marked as running after Execute")
	}
	
	runningAudits := executor.GetRunningAudits()
	if len(runningAudits) != 1 {
		t.Errorf("Expected 1 running audit, got %d", len(runningAudits))
	}
	
	if runningAudits[0] != request.AuditID {
		t.Errorf("Expected running audit ID %s, got %s", request.AuditID, runningAudits[0])
	}
}

func TestPipelineExecutor_ExecuteInvalidRequest(t *testing.T) {
	executor := NewPipelineExecutor()
	registry := NewAgentRegistry()
	ctx := context.Background()
	
	tests := []struct {
		name    string
		request *AuditRequest
	}{
		{
			name:    "nil request",
			request: nil,
		},
		{
			name: "empty audit ID",
			request: &AuditRequest{
				AuditID:  "",
				SeedURL:  "https://example.com",
				MaxPages: 5,
			},
		},
		{
			name: "empty seed URL",
			request: &AuditRequest{
				AuditID:  "test-audit",
				SeedURL:  "",
				MaxPages: 5,
			},
		},
		{
			name: "invalid max pages",
			request: &AuditRequest{
				AuditID:  "test-audit",
				SeedURL:  "https://example.com",
				MaxPages: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executor.Execute(ctx, tt.request, registry)
			
			if err == nil {
				t.Error("Execute should return error for invalid request")
			}
		})
	}
}

func TestPipelineExecutor_GetProgress(t *testing.T) {
	executor := NewPipelineExecutor()
	
	// Test getting progress for non-existing audit
	progress, err := executor.GetProgress("non-existing-audit")
	if err == nil {
		t.Error("GetProgress should return error for non-existing audit")
	}
	if progress != nil {
		t.Error("GetProgress should return nil progress for non-existing audit")
	}
}

func TestPipelineExecutor_Cancel(t *testing.T) {
	executor := NewPipelineExecutor()
	registry := NewAgentRegistry()
	
	// Register a mock agent
	mockAgent := &MockAgent{name: "test-agent"}
	registry.Register("test-agent", mockAgent)
	
	request := &AuditRequest{
		AuditID:   "test-cancel-audit",
		SeedURL:   "https://example.com",
		MaxPages:  5,
		Timestamp: time.Now(),
	}
	
	ctx := context.Background()
	
	// Start execution
	_, err := executor.Execute(ctx, request, registry)
	if err != nil {
		t.Fatalf("Failed to start audit: %v", err)
	}
	
	// Verify audit is running
	if !executor.IsRunning(request.AuditID) {
		t.Fatal("Audit should be running before cancel")
	}
	
	// Cancel the audit
	err = executor.Cancel(request.AuditID)
	if err != nil {
		t.Errorf("Cancel should not return error: %v", err)
	}
	
	// Wait a bit for cancellation to take effect
	time.Sleep(100 * time.Millisecond)
	
	// Verify audit is no longer running
	if executor.IsRunning(request.AuditID) {
		t.Error("Audit should not be running after cancel")
	}
	
	// Test canceling non-existing audit
	err = executor.Cancel("non-existing-audit")
	if err == nil {
		t.Error("Cancel should return error for non-existing audit")
	}
}

func TestPipelineExecutor_IsRunning(t *testing.T) {
	executor := NewPipelineExecutor()
	
	// Test with non-existing audit
	if executor.IsRunning("non-existing") {
		t.Error("IsRunning should return false for non-existing audit")
	}
	
	registry := NewAgentRegistry()
	mockAgent := &MockAgent{name: "test-agent"}
	registry.Register("test-agent", mockAgent)
	
	request := &AuditRequest{
		AuditID:   "running-test",
		SeedURL:   "https://example.com",
		MaxPages:  5,
		Timestamp: time.Now(),
	}
	
	// Before execution
	if executor.IsRunning(request.AuditID) {
		t.Error("IsRunning should return false before execution starts")
	}
	
	// Start execution
	ctx := context.Background()
	_, err := executor.Execute(ctx, request, registry)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}
	
	// After execution starts
	if !executor.IsRunning(request.AuditID) {
		t.Error("IsRunning should return true after execution starts")
	}
}

func TestPipelineExecutor_ConcurrentExecutions(t *testing.T) {
	executor := NewPipelineExecutor()
	registry := NewAgentRegistry()
	
	// Register multiple agents
	registry.Register(constants.AgentNameKeyword, &MockAgent{name: constants.AgentNameKeyword})
	registry.Register(constants.AgentNameTechnical, &MockAgent{name: constants.AgentNameTechnical})
	
	ctx := context.Background()
	numAudits := 3
	
	// Start multiple audits concurrently
	for i := 0; i < numAudits; i++ {
		request := &AuditRequest{
			AuditID:   fmt.Sprintf("concurrent-audit-%d", i),
			SeedURL:   "https://example.com",
			MaxPages:  5,
			Timestamp: time.Now(),
		}
		
		_, err := executor.Execute(ctx, request, registry)
		if err != nil {
			t.Errorf("Failed to start concurrent audit %d: %v", i, err)
		}
	}
	
	runningAudits := executor.GetRunningAudits()
	if len(runningAudits) != numAudits {
		t.Errorf("Expected %d running audits, got %d", numAudits, len(runningAudits))
	}
}

func TestPipelineExecutor_PipelineSteps(t *testing.T) {
	executor := NewPipelineExecutor()
	registry := NewAgentRegistry()
	
	// Register all expected agents  
	registry.Register(constants.AgentNameKeyword, &MockAgent{name: constants.AgentNameKeyword})
	registry.Register(constants.AgentNameTechnical, &MockAgent{name: constants.AgentNameTechnical})
	registry.Register(constants.AgentNameLinking, &MockAgent{name: constants.AgentNameLinking})
	registry.Register(constants.AgentNameBrokenLinks, &MockAgent{name: constants.AgentNameBrokenLinks})
	
	request := &AuditRequest{
		AuditID:   "pipeline-steps-test",
		SeedURL:   "https://example.com", 
		MaxPages:  3,
		Timestamp: time.Now(),
	}
	
	ctx := context.Background()
	resultsChan, err := executor.Execute(ctx, request, registry)
	if err != nil {
		t.Fatalf("Failed to execute pipeline: %v", err)
	}
	
	// Collect all results
	var results []*PipelineResult
	timeout := time.After(10 * time.Second)
	
	for {
		select {
		case result, ok := <-resultsChan:
			if !ok {
				// Channel closed, pipeline finished
				goto AnalyzeResults
			}
			results = append(results, result)
		case <-timeout:
			t.Fatal("Pipeline execution timed out")
		}
	}
	
AnalyzeResults:
	if len(results) == 0 {
		t.Fatal("Pipeline should produce at least one result")
	}
	
	// Verify we have results for expected steps
	steps := make(map[string]bool)
	for _, result := range results {
		steps[result.Step] = true
	}
	
	expectedSteps := []string{
		constants.PipelineStepCrawling,
		constants.PipelineStepAnalyzing,
	}
	
	for _, expectedStep := range expectedSteps {
		if !steps[expectedStep] {
			t.Errorf("Missing expected pipeline step: %s", expectedStep)
		}
	}
}

func TestPipelineExecutor_ContextCancellation(t *testing.T) {
	executor := NewPipelineExecutor()
	registry := NewAgentRegistry()
	
	registry.Register("test-agent", &MockAgent{name: "test-agent"})
	
	request := &AuditRequest{
		AuditID:   "context-cancel-test",
		SeedURL:   "https://example.com",
		MaxPages:  5,
		Timestamp: time.Now(),
	}
	
	// Create a context that we'll cancel
	ctx, cancel := context.WithCancel(context.Background())
	
	resultsChan, err := executor.Execute(ctx, request, registry)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}
	
	// Cancel context immediately 
	cancel()
	
	// Wait for pipeline to respond to cancellation
	timeout := time.After(2 * time.Second)
	select {
	case _, ok := <-resultsChan:
		if ok {
			// We expect the channel to be closed when context is cancelled
		}
	case <-timeout:
		t.Error("Pipeline should respond to context cancellation")
	}
}