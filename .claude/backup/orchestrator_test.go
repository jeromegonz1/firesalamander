package orchestrator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrchestrator(t *testing.T) {
	// Skip test if config file doesn't exist
	orch, err := NewOrchestrator()
	if err != nil {
		t.Skipf("Skipping test - config file not found: %v", err)
		return
	}

	assert.NotNil(t, orch)
	assert.NotNil(t, orch.Status)
}

func TestStartAudit(t *testing.T) {
	orch := &Orchestrator{
		Status: make(map[string]*AuditStatus),
	}

	req := AuditRequest{
		SeedURL: "https://example.com",
		AuditID: "test_audit_123",
		Options: make(map[string]interface{}),
	}

	err := orch.StartAudit(req)
	require.NoError(t, err)

	// Check status was created
	status, err := orch.GetStatus("test_audit_123")
	require.NoError(t, err)
	assert.Equal(t, "test_audit_123", status.ID)
	assert.Equal(t, "pending", status.Status)
}

func TestGetStatus(t *testing.T) {
	orch := &Orchestrator{
		Status: map[string]*AuditStatus{
			"existing_audit": {
				ID:        "existing_audit",
				Status:    "completed",
				Progress:  100,
				StartedAt: time.Now(),
			},
		},
	}

	// Test existing audit
	status, err := orch.GetStatus("existing_audit")
	require.NoError(t, err)
	assert.Equal(t, "existing_audit", status.ID)
	assert.Equal(t, "completed", status.Status)
	assert.Equal(t, 100, status.Progress)

	// Test non-existing audit
	_, err = orch.GetStatus("non_existing")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDuplicateAudit(t *testing.T) {
	orch := &Orchestrator{
		Status: map[string]*AuditStatus{
			"duplicate_test": {
				ID:     "duplicate_test",
				Status: "pending",
			},
		},
	}

	req := AuditRequest{
		SeedURL: "https://example.com",
		AuditID: "duplicate_test",
		Options: make(map[string]interface{}),
	}

	err := orch.StartAudit(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}