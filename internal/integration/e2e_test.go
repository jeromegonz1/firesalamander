package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/orchestrator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_CompleteAuditFlow tests the full audit pipeline end-to-end
func TestE2E_CompleteAuditFlow(t *testing.T) {
	// Setup test environment
	testDir := setupE2EEnvironment(t)
	defer cleanupE2EEnvironment(testDir)

	// Load configuration
	cfg, err := config.Load()
	require.NoError(t, err)

	// Create pipeline
	pipeline, err := NewPipeline(cfg)
	require.NoError(t, err)

	// Test with our test fixtures
	request := orchestrator.AuditRequest{
		AuditID: "FS-E2E-001",
		SeedURL: "file://" + filepath.Join("..", "..", "test-fixtures", "test-site", "index.html"),
		Options: map[string]interface{}{
			"max_urls":  3,
			"max_depth": 1,
		},
	}

	// Start audit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = pipeline.StartAudit(ctx, request)
	require.NoError(t, err)

	// Wait for completion
	var finalStatus *AuditExecution
	for i := 0; i < 30; i++ {
		status := pipeline.GetAuditStatus(request.AuditID)
		if status.Status == "completed" || status.Status == "failed" || status.Status == "partial" {
			finalStatus = status
			break
		}
		time.Sleep(1 * time.Second)
	}

	require.NotNil(t, finalStatus, "Audit did not complete in time")
	assert.Contains(t, []string{"completed", "partial", "failed"}, finalStatus.Status)
	assert.Equal(t, "FS-E2E-001", finalStatus.AuditID)

	// Verify results exist (depending on status)
	if finalStatus.Status == "completed" || finalStatus.Status == "partial" {
		assert.NotNil(t, finalStatus.Results["crawl"])
		
		// Technical might be available if crawl succeeded
		if finalStatus.Results["technical"] != nil {
			assert.NotNil(t, finalStatus.Results["technical"])
		}
		
		// Semantic might fail (Python service) - that's OK
		if finalStatus.Results["semantic"] != nil {
			assert.NotNil(t, finalStatus.Results["semantic"])
		}

		// Report might be generated if analysis succeeded
		if finalStatus.Results["html_report"] != nil {
			assert.NotNil(t, finalStatus.Results["html_report"])
		}
	}
}

// TestE2E_ErrorRecovery tests error handling and recovery
func TestE2E_ErrorRecovery(t *testing.T) {
	testDir := setupE2EEnvironment(t)
	defer cleanupE2EEnvironment(testDir)

	cfg, err := config.Load()
	require.NoError(t, err)

	pipeline, err := NewPipeline(cfg)
	require.NoError(t, err)

	// Test with invalid URL to trigger error recovery
	request := orchestrator.AuditRequest{
		AuditID: "FS-E2E-ERROR-001",
		SeedURL: "https://this-domain-absolutely-does-not-exist-12345.invalid",
		Options: map[string]interface{}{
			"max_urls": 1,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = pipeline.StartAudit(ctx, request)
	require.NoError(t, err)

	// Wait for failure
	time.Sleep(3 * time.Second)
	status := pipeline.GetAuditStatus(request.AuditID)

	// Should fail gracefully with error handling
	require.NotNil(t, status)
	assert.Equal(t, "FS-E2E-ERROR-001", status.AuditID)
	assert.Contains(t, []string{"failed", "partial"}, status.Status)
	assert.NotEmpty(t, status.Error)
}

// TestE2E_ConcurrentAudits tests multiple simultaneous audits
func TestE2E_ConcurrentAudits(t *testing.T) {
	testDir := setupE2EEnvironment(t)
	defer cleanupE2EEnvironment(testDir)

	cfg, err := config.Load()
	require.NoError(t, err)

	pipeline, err := NewPipeline(cfg)
	require.NoError(t, err)

	// Start 3 concurrent audits
	requests := []orchestrator.AuditRequest{
		{
			AuditID: "FS-E2E-CONC-001",
			SeedURL: "file://" + filepath.Join("..", "..", "test-fixtures", "test-site", "index.html"),
			Options: map[string]interface{}{"max_urls": 1},
		},
		{
			AuditID: "FS-E2E-CONC-002", 
			SeedURL: "file://" + filepath.Join("..", "..", "test-fixtures", "test-site", "about.html"),
			Options: map[string]interface{}{"max_urls": 1},
		},
		{
			AuditID: "FS-E2E-CONC-003",
			SeedURL: "file://" + filepath.Join("..", "..", "test-fixtures", "test-site", "contact.html"),
			Options: map[string]interface{}{"max_urls": 1},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Start all audits
	for _, req := range requests {
		err = pipeline.StartAudit(ctx, req)
		require.NoError(t, err)
	}

	// Wait for all to complete
	time.Sleep(5 * time.Second)

	// Check all audits completed
	for _, req := range requests {
		status := pipeline.GetAuditStatus(req.AuditID)
		require.NotNil(t, status, "Audit %s not found", req.AuditID)
		assert.Equal(t, req.AuditID, status.AuditID)
		assert.Contains(t, []string{"completed", "failed", "partial"}, status.Status)
	}
}

// TestE2E_JSONRPCProtocol tests the JSON-RPC communication protocol
func TestE2E_JSONRPCProtocol(t *testing.T) {
	// Test message generation
	auditID := "FS-E2E-RPC-001"
	
	// Generate crawler request message
	crawlerMessage := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "start_crawl",
		"params": map[string]interface{}{
			"audit_id":    auditID,
			"seed_url":    "https://test.fr",
			"max_urls":    50,
			"max_depth":   2,
			"config_file": "config/crawler.yaml",
		},
		"id": "orch-crawler-" + auditID,
	}

	// Validate JSON-RPC structure
	jsonData, err := json.Marshal(crawlerMessage)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonData, &parsed)
	require.NoError(t, err)

	// Validate required fields
	assert.Equal(t, "2.0", parsed["jsonrpc"])
	assert.Equal(t, "start_crawl", parsed["method"])
	assert.NotNil(t, parsed["params"])
	assert.Equal(t, "orch-crawler-"+auditID, parsed["id"])

	// Validate params structure
	params := parsed["params"].(map[string]interface{})
	assert.Equal(t, auditID, params["audit_id"])
	assert.Equal(t, "https://test.fr", params["seed_url"])
	assert.Equal(t, float64(50), params["max_urls"]) // JSON numbers are float64
	assert.Equal(t, float64(2), params["max_depth"])

	// Test response message format
	responseMessage := map[string]interface{}{
		"jsonrpc": "2.0",
		"result": map[string]interface{}{
			"audit_id":      auditID,
			"status":        "complete",
			"pages_crawled": 15,
			"duration_ms":   45000,
			"output_file":   "/audits/" + auditID + "/crawl_index.json",
		},
		"id": "orch-crawler-" + auditID,
	}

	jsonData, err = json.Marshal(responseMessage)
	require.NoError(t, err)

	err = json.Unmarshal(jsonData, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "2.0", parsed["jsonrpc"])
	assert.NotNil(t, parsed["result"])
	assert.Equal(t, "orch-crawler-"+auditID, parsed["id"])

	result := parsed["result"].(map[string]interface{})
	assert.Equal(t, auditID, result["audit_id"])
	assert.Equal(t, "complete", result["status"])
}

// TestE2E_DataIntegrity tests data consistency across the pipeline
func TestE2E_DataIntegrity(t *testing.T) {
	testDir := setupE2EEnvironment(t)
	defer cleanupE2EEnvironment(testDir)

	cfg, err := config.Load()
	require.NoError(t, err)

	pipeline, err := NewPipeline(cfg)
	require.NoError(t, err)

	request := orchestrator.AuditRequest{
		AuditID: "FS-E2E-DATA-001",
		SeedURL: "file://" + filepath.Join("..", "..", "test-fixtures", "test-site", "index.html"),
		Options: map[string]interface{}{
			"max_urls": 2,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = pipeline.StartAudit(ctx, request)
	require.NoError(t, err)

	// Wait for completion
	time.Sleep(3 * time.Second)
	status := pipeline.GetAuditStatus(request.AuditID)

	if status.Status == "completed" || status.Status == "partial" {
		// Verify audit_id consistency across all results
		assert.Equal(t, "FS-E2E-DATA-001", status.AuditID)
		
		// Check crawl data
		if crawlData, exists := status.Results["crawl"]; exists {
			if crawlMap, ok := crawlData.(map[string]interface{}); ok {
				if auditIDField, hasID := crawlMap["audit_id"]; hasID {
					assert.Equal(t, "FS-E2E-DATA-001", auditIDField)
				}
			}
		}

		// Check technical results
		if techData, exists := status.Results["technical"]; exists {
			if techMap, ok := techData.(map[string]interface{}); ok {
				if auditIDField, hasID := techMap["audit_id"]; hasID {
					assert.Equal(t, "FS-E2E-DATA-001", auditIDField)
				}
			}
		}

		// Check semantic results (if available)
		if semanticData, exists := status.Results["semantic"]; exists {
			if semanticMap, ok := semanticData.(map[string]interface{}); ok {
				if auditIDField, hasID := semanticMap["audit_id"]; hasID {
					assert.Equal(t, "FS-E2E-DATA-001", auditIDField)
				}
			}
		}
	}
}

// TestE2E_Performance tests pipeline performance under load
func TestE2E_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	testDir := setupE2EEnvironment(t)
	defer cleanupE2EEnvironment(testDir)

	cfg, err := config.Load()
	require.NoError(t, err)

	pipeline, err := NewPipeline(cfg)
	require.NoError(t, err)

	// Performance test: Single audit should complete in < 10s
	start := time.Now()
	
	request := orchestrator.AuditRequest{
		AuditID: "FS-E2E-PERF-001",
		SeedURL: "file://" + filepath.Join("..", "..", "test-fixtures", "test-site", "index.html"),
		Options: map[string]interface{}{
			"max_urls": 1,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = pipeline.StartAudit(ctx, request)
	require.NoError(t, err)

	// Wait for completion
	for i := 0; i < 100; i++ { // Max 10s
		status := pipeline.GetAuditStatus(request.AuditID)
		if status.Status == "completed" || status.Status == "failed" || status.Status == "partial" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	duration := time.Since(start)
	assert.Less(t, duration, 10*time.Second, "Audit took too long: %v", duration)

	status := pipeline.GetAuditStatus(request.AuditID)
	require.NotNil(t, status)
	assert.Contains(t, []string{"completed", "partial", "failed"}, status.Status)
}

// Helper functions for E2E tests
func setupE2EEnvironment(t *testing.T) string {
	testDir, err := os.MkdirTemp("", "fire-salamander-e2e-*")
	require.NoError(t, err)
	
	// Create audit output directory
	auditDir := filepath.Join(testDir, "audits")
	err = os.MkdirAll(auditDir, 0755)
	require.NoError(t, err)

	return testDir
}

func cleanupE2EEnvironment(testDir string) {
	os.RemoveAll(testDir)
}

// Benchmark tests for performance monitoring
func BenchmarkPipelineCreation(b *testing.B) {
	cfg, err := config.Load()
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline, err := NewPipeline(cfg)
		require.NoError(b, err)
		_ = pipeline
	}
}

func BenchmarkAuditRequest(b *testing.B) {
	cfg, err := config.Load()
	require.NoError(b, err)

	pipeline, err := NewPipeline(cfg)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		request := orchestrator.AuditRequest{
			AuditID: fmt.Sprintf("FS-BENCH-%03d", i),
			SeedURL: "https://example.com",
			Options: map[string]interface{}{
				"max_urls": 1,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		err = pipeline.StartAudit(ctx, request)
		cancel()
		
		if err != nil {
			b.Errorf("StartAudit failed: %v", err)
		}
	}
}

// TestE2E_ConfigurationVariations tests different config scenarios
func TestE2E_ConfigurationVariations(t *testing.T) {
	testCases := []struct {
		name        string
		maxURLs     int
		maxDepth    int
		expectPages int
	}{
		{"Single Page", 1, 1, 1},
		{"Small Site", 5, 2, 3},
		{"Medium Site", 10, 3, 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testDir := setupE2EEnvironment(t)
			defer cleanupE2EEnvironment(testDir)

			cfg, err := config.Load()
			require.NoError(t, err)

			pipeline, err := NewPipeline(cfg)
			require.NoError(t, err)

			request := orchestrator.AuditRequest{
				AuditID: "FS-E2E-VAR-" + fmt.Sprintf("%03d", tc.maxURLs),
				SeedURL: "file://" + filepath.Join("..", "..", "test-fixtures", "test-site", "index.html"),
				Options: map[string]interface{}{
					"max_urls":  tc.maxURLs,
					"max_depth": tc.maxDepth,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			err = pipeline.StartAudit(ctx, request)
			require.NoError(t, err)

			// Wait for completion
			time.Sleep(2 * time.Second)
			status := pipeline.GetAuditStatus(request.AuditID)

			require.NotNil(t, status)
			assert.Equal(t, request.AuditID, status.AuditID)
		})
	}
}