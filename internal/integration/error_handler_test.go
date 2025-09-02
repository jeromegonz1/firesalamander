package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandler_HandleCrawlerError(t *testing.T) {
	eh := NewErrorHandler()
	
	// Test network error
	networkErr := errors.New("connection refused")
	action, err := eh.HandleCrawlerError(context.Background(), networkErr, "FS-TEST-001")
	require.NoError(t, err)
	
	assert.Equal(t, ErrorTypeNetwork, action.Type)
	assert.Equal(t, "network_retry", action.Action)
	assert.False(t, action.Fallback)
	assert.Contains(t, action.LastError, "connection refused")

	// Test timeout error
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond) // Force timeout
	
	timeoutErr := errors.New("context deadline exceeded")
	action, err = eh.HandleCrawlerError(ctx, timeoutErr, "FS-TEST-002")
	require.NoError(t, err)
	
	assert.Equal(t, ErrorTypeTimeout, action.Type)
	assert.Equal(t, "timeout_fallback", action.Action)
	assert.True(t, action.Fallback)
}

func TestErrorHandler_HandleSemanticError(t *testing.T) {
	eh := NewErrorHandler()
	
	semanticErr := errors.New("Python service unavailable")
	action, err := eh.HandleSemanticError(context.Background(), semanticErr, "FS-TEST-003")
	require.NoError(t, err)
	
	assert.Equal(t, ErrorTypeSemantic, action.Type)
	assert.Equal(t, "skip_semantic", action.Action)
	assert.True(t, action.Fallback)
	assert.Contains(t, action.LastError, "Python service")
}

func TestErrorHandler_HandleTechnicalError(t *testing.T) {
	eh := NewErrorHandler()
	
	techErr := errors.New("rule validation failed")
	action, err := eh.HandleTechnicalError(context.Background(), techErr, "FS-TEST-004")
	require.NoError(t, err)
	
	assert.Equal(t, ErrorTypeTechnical, action.Type)
	assert.Equal(t, "basic_technical_analysis", action.Action)
	assert.True(t, action.Fallback)
}

func TestErrorHandler_ShouldRetry(t *testing.T) {
	eh := NewErrorHandler()
	
	// Network error should retry
	networkAction := &RecoveryAction{
		Type:       ErrorTypeNetwork,
		RetryCount: 1,
		MaxRetries: 3,
	}
	assert.True(t, eh.ShouldRetry(networkAction))
	
	// Max retries reached - no retry
	maxRetriesAction := &RecoveryAction{
		Type:       ErrorTypeNetwork,
		RetryCount: 3,
		MaxRetries: 3,
	}
	assert.False(t, eh.ShouldRetry(maxRetriesAction))
	
	// Timeout error - no retry
	timeoutAction := &RecoveryAction{
		Type:       ErrorTypeTimeout,
		RetryCount: 0,
		MaxRetries: 3,
	}
	assert.False(t, eh.ShouldRetry(timeoutAction))
}

func TestErrorHandler_ExecuteRecovery(t *testing.T) {
	eh := NewErrorHandler()
	
	execution := &AuditExecution{
		AuditID: "FS-RECOVERY-001",
		Status:  "failed",
		Results: make(map[string]interface{}),
	}
	
	// Test timeout fallback
	timeoutAction := &RecoveryAction{
		Type:   ErrorTypeTimeout,
		Action: "timeout_fallback",
	}
	
	err := eh.ExecuteRecovery(context.Background(), timeoutAction, nil, execution)
	require.NoError(t, err)
	
	assert.Equal(t, "partial", execution.Status)
	assert.Equal(t, "timeout_recovery", execution.CurrentStep)
	assert.Contains(t, execution.Error, "Timeout")

	// Test skip semantic
	execution.Status = "failed" // Reset
	skipAction := &RecoveryAction{
		Type:   ErrorTypeSemantic,
		Action: "skip_semantic",
	}
	
	err = eh.ExecuteRecovery(context.Background(), skipAction, nil, execution)
	require.NoError(t, err)
	
	semanticResult, exists := execution.Results["semantic"]
	require.True(t, exists)
	
	semanticMap := semanticResult.(map[string]interface{})
	assert.Equal(t, "FS-RECOVERY-001", semanticMap["audit_id"])
	assert.Equal(t, "skipped", semanticMap["status"])
	assert.Equal(t, true, semanticMap["fallback"])
}

func TestIsNetworkError(t *testing.T) {
	testCases := []struct {
		err      error
		expected bool
	}{
		{errors.New("connection refused"), true},
		{errors.New("timeout occurred"), true},
		{errors.New("no such host"), true},
		{errors.New("certificate error"), true},
		{errors.New("network unreachable"), true},
		{errors.New("validation failed"), false},
		{errors.New("parsing error"), false},
	}
	
	for _, tc := range testCases {
		result := isNetworkError(tc.err)
		assert.Equal(t, tc.expected, result, "Error: %s", tc.err.Error())
	}
}

func TestAuditIDFormat(t *testing.T) {
	testCases := []struct {
		auditID string
		valid   bool
	}{
		{"FS-PROD-001", true},
		{"FS-TEST-042", true},
		{"FS-DEV-123", true},
		{"FS-PROD-999", true},
		{"INVALID-001", false},
		{"FS-001", false},
		{"FS-PROD-", false},
		{"FS-PROD-1000", false}, // Too high
	}
	
	for _, tc := range testCases {
		result := isValidAuditID(tc.auditID)
		assert.Equal(t, tc.valid, result, "Audit ID: %s", tc.auditID)
	}
}

// Helper function to validate audit ID format
func isValidAuditID(auditID string) bool {
	if len(auditID) < 10 || len(auditID) > 11 {
		return false
	}
	if auditID[:3] != "FS-" {
		return false
	}
	
	// Find the last dash
	lastDash := -1
	for i := len(auditID) - 1; i >= 3; i-- {
		if auditID[i] == '-' {
			lastDash = i
			break
		}
	}
	
	if lastDash == -1 {
		return false
	}
	
	env := auditID[3:lastDash]
	if env != "PROD" && env != "TEST" && env != "DEV" {
		return false
	}
	
	// Check if sequence after last dash are digits
	seq := auditID[lastDash+1:]
	if len(seq) != 3 {
		return false
	}
	for _, char := range seq {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}