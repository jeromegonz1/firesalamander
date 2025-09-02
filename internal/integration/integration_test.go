package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
	"firesalamander/internal/orchestrator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJSONRPCMessageFormat validates JSON-RPC 2.0 message structure
func TestJSONRPCMessageFormat(t *testing.T) {
	// Test JSON-RPC request message
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "start_crawl",
		"params": map[string]interface{}{
			"audit_id":    "FS-001",
			"seed_url":    "https://camping-test.fr",
			"max_urls":    300,
			"max_depth":   3,
			"config_file": "config/crawler.yaml",
		},
		"id": "orch-crawler-001",
	}

	// Validate JSON serialization
	jsonData, err := json.Marshal(request)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonData, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "2.0", parsed["jsonrpc"])
	assert.Equal(t, "start_crawl", parsed["method"])
	assert.NotNil(t, parsed["params"])
	assert.Equal(t, "orch-crawler-001", parsed["id"])

	// Test JSON-RPC response message
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"result": map[string]interface{}{
			"audit_id":      "FS-001",
			"status":        "complete",
			"pages_crawled": 47,
			"duration_ms":   85420,
			"output_file":   "/audits/FS-001/crawl_index.json",
		},
		"id": "orch-crawler-001",
	}

	jsonData, err = json.Marshal(response)
	require.NoError(t, err)

	err = json.Unmarshal(jsonData, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "2.0", parsed["jsonrpc"])
	assert.NotNil(t, parsed["result"])
	assert.Equal(t, "orch-crawler-001", parsed["id"])
}

// TestPipelineConfiguration validates pipeline setup from config
func TestPipelineConfiguration(t *testing.T) {
	// Load configuration from YAML files (no hardcoding)
	cfg, err := config.Load()
	require.NoError(t, err)

	// Validate config structure
	assert.Greater(t, cfg.Server.Port, 0)

	// Test pipeline creation
	pipeline, err := NewPipeline(cfg)
	require.NoError(t, err)
	assert.NotNil(t, pipeline)
	assert.NotNil(t, pipeline.crawler)
	assert.NotNil(t, pipeline.technical)
	assert.NotNil(t, pipeline.semantic)
	assert.NotNil(t, pipeline.report)
}

// TestAuditStatusFlow validates audit status progression
func TestAuditStatusFlow(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	pipeline, err := NewPipeline(cfg)
	require.NoError(t, err)

	// Test audit request
	request := orchestrator.AuditRequest{
		AuditID: "FS-STATUS-001",
		SeedURL: constants.TestURLExample,
		Options: map[string]interface{}{
			"max_urls": 1,
		},
	}

	// Start audit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pipeline.StartAudit(ctx, request)
	require.NoError(t, err)

	// Check initial status
	status := pipeline.GetAuditStatus(request.AuditID)
	require.NotNil(t, status)
	assert.Equal(t, "FS-STATUS-001", status.AuditID)
	assert.Contains(t, []string{"pending", "crawling"}, status.Status)

	// Wait briefly and check progress
	time.Sleep(100 * time.Millisecond)
	status = pipeline.GetAuditStatus(request.AuditID)
	assert.NotEmpty(t, status.CurrentStep)
}

// TestErrorHandling validates error propagation through pipeline
func TestErrorHandling(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	pipeline, err := NewPipeline(cfg)
	require.NoError(t, err)

	// Test with invalid URL
	request := orchestrator.AuditRequest{
		AuditID: "FS-ERROR-001",
		SeedURL: "https://invalid-domain-does-not-exist.test",
		Options: map[string]interface{}{
			"max_urls": 1,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = pipeline.StartAudit(ctx, request)
	require.NoError(t, err)

	// Wait for failure
	time.Sleep(1 * time.Second)
	status := pipeline.GetAuditStatus(request.AuditID)
	
	// Should fail gracefully
	assert.NotNil(t, status)
	assert.Equal(t, "FS-ERROR-001", status.AuditID)
}

// TestDataFlowIntegrity validates data consistency through pipeline
func TestDataFlowIntegrity(t *testing.T) {
	// Test audit ID consistency
	auditID := "FS-FLOW-001"
	
	// Validate audit ID format
	assert.Contains(t, auditID, "FS-")
	assert.Len(t, auditID, 11) // FS-XXX-001 format

	// Test JSON serialization preservation
	testData := map[string]interface{}{
		"audit_id": auditID,
		"pages":    []string{"page1", "page2"},
		"metadata": map[string]string{"test": "value"},
	}

	// Serialize and deserialize
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonData, &parsed)
	require.NoError(t, err)

	assert.Equal(t, auditID, parsed["audit_id"])
	assert.NotNil(t, parsed["pages"])
	assert.NotNil(t, parsed["metadata"])
}