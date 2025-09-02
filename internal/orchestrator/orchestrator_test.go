package v2

import (
	"context"
	"fmt"
	"testing"
	"time"

	"firesalamander/internal/agents"
	"firesalamander/internal/agents/broken"
	"firesalamander/internal/agents/keyword"
	"firesalamander/internal/agents/linking" 
	"firesalamander/internal/agents/technical"
	"firesalamander/internal/constants"
)

func TestNewOrchestratorV2(t *testing.T) {
	orchestrator := NewOrchestratorV2()
	
	if orchestrator == nil {
		t.Fatal("NewOrchestratorV2 should not return nil")
	}
	
	// Vérifier l'état initial
	activeAudits := orchestrator.ListActiveAudits()
	if len(activeAudits) != 0 {
		t.Errorf("Expected 0 active audits initially, got %d", len(activeAudits))
	}
}

func TestOrchestratorV2_RegisterAgent(t *testing.T) {
	orchestrator := NewOrchestratorV2()
	
	tests := []struct {
		name        string
		agentName   string
		agent       interface{}
		expectError bool
	}{
		{
			name:        "register keyword agent",
			agentName:   constants.AgentNameKeyword,
			agent:       keyword.NewKeywordExtractor(),
			expectError: false,
		},
		{
			name:        "register technical agent",
			agentName:   constants.AgentNameTechnical,
			agent:       technical.NewTechnicalAuditor(),
			expectError: false,
		},
		{
			name:        "empty name should fail",
			agentName:   "",
			agent:       keyword.NewKeywordExtractor(),
			expectError: true,
		},
		{
			name:        "nil agent should fail",
			agentName:   "nil-agent",
			agent:       nil,
			expectError: true,
		},
		{
			name:        "duplicate registration should fail",
			agentName:   constants.AgentNameKeyword, // Already registered
			agent:       keyword.NewKeywordExtractor(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.agent != nil {
				if agent, ok := tt.agent.(agents.Agent); ok {
					err = orchestrator.RegisterAgent(tt.agentName, agent)
				} else {
					err = fmt.Errorf("invalid agent type")
				}
			} else {
				err = orchestrator.RegisterAgent(tt.agentName, nil)
			}
			
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestOrchestratorV2_StartAudit(t *testing.T) {
	orchestrator := NewOrchestratorV2()
	
	// Enregistrer quelques agents
	orchestrator.RegisterAgent(constants.AgentNameKeyword, keyword.NewKeywordExtractor())
	orchestrator.RegisterAgent(constants.AgentNameTechnical, technical.NewTechnicalAuditor())
	
	request := &AuditRequest{
		AuditID:   "test-audit-123",
		SeedURL:   "https://example.com",
		MaxPages:  10,
		Options:   map[string]interface{}{"test": true},
		Timestamp: time.Now(),
	}
	
	ctx := context.Background()
	progressChan, err := orchestrator.StartAudit(ctx, request)
	
	if err != nil {
		t.Errorf("StartAudit should not return error: %v", err)
	}
	
	if progressChan == nil {
		t.Fatal("StartAudit should return a progress channel")
	}
	
	// Vérifier que l'audit est maintenant actif
	activeAudits := orchestrator.ListActiveAudits()
	if len(activeAudits) != 1 {
		t.Errorf("Expected 1 active audit, got %d", len(activeAudits))
	}
	
	if activeAudits[0].AuditID != request.AuditID {
		t.Errorf("Expected audit ID %s, got %s", request.AuditID, activeAudits[0].AuditID)
	}
}

func TestOrchestratorV2_StartAuditInvalidRequest(t *testing.T) {
	orchestrator := NewOrchestratorV2()
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := orchestrator.StartAudit(ctx, tt.request)
			
			if err == nil {
				t.Error("StartAudit should return error for invalid request")
			}
		})
	}
}

func TestOrchestratorV2_GetAuditStatus(t *testing.T) {
	orchestrator := NewOrchestratorV2()
	
	// Test avec audit inexistant
	status, err := orchestrator.GetAuditStatus("non-existing-audit")
	if err == nil {
		t.Error("GetAuditStatus should return error for non-existing audit")
	}
	if status != nil {
		t.Error("GetAuditStatus should return nil status for non-existing audit")
	}
	
	// Démarrer un audit réel
	orchestrator.RegisterAgent(constants.AgentNameKeyword, keyword.NewKeywordExtractor())
	
	request := &AuditRequest{
		AuditID:   "status-test-audit",
		SeedURL:   "https://example.com",
		MaxPages:  5,
		Timestamp: time.Now(),
	}
	
	ctx := context.Background()
	_, err = orchestrator.StartAudit(ctx, request)
	if err != nil {
		t.Fatalf("Failed to start audit: %v", err)
	}
	
	// Maintenant tester le statut
	status, err = orchestrator.GetAuditStatus(request.AuditID)
	if err != nil {
		t.Errorf("GetAuditStatus should not return error: %v", err)
	}
	
	if status == nil {
		t.Fatal("GetAuditStatus should return status")
	}
	
	if status.AuditID != request.AuditID {
		t.Errorf("Expected audit ID %s, got %s", request.AuditID, status.AuditID)
	}
}

func TestOrchestratorV2_StreamProgress(t *testing.T) {
	orchestrator := NewOrchestratorV2()
	orchestrator.RegisterAgent(constants.AgentNameKeyword, keyword.NewKeywordExtractor())
	
	request := &AuditRequest{
		AuditID:   "stream-test-audit",
		SeedURL:   "https://example.com",
		MaxPages:  5,
		Timestamp: time.Now(),
	}
	
	ctx := context.Background()
	_, err := orchestrator.StartAudit(ctx, request)
	if err != nil {
		t.Fatalf("Failed to start audit: %v", err)
	}
	
	// Test streaming
	streamChan, err := orchestrator.StreamProgress(request.AuditID)
	if err != nil {
		t.Errorf("StreamProgress should not return error: %v", err)
	}
	
	if streamChan == nil {
		t.Fatal("StreamProgress should return a channel")
	}
	
	// Test streaming pour audit inexistant
	_, err = orchestrator.StreamProgress("non-existing-audit")
	if err == nil {
		t.Error("StreamProgress should return error for non-existing audit")
	}
}

func TestOrchestratorV2_CancelAudit(t *testing.T) {
	orchestrator := NewOrchestratorV2()
	orchestrator.RegisterAgent(constants.AgentNameKeyword, keyword.NewKeywordExtractor())
	
	request := &AuditRequest{
		AuditID:   "cancel-test-audit",
		SeedURL:   "https://example.com",
		MaxPages:  5,
		Timestamp: time.Now(),
	}
	
	ctx := context.Background()
	_, err := orchestrator.StartAudit(ctx, request)
	if err != nil {
		t.Fatalf("Failed to start audit: %v", err)
	}
	
	// Vérifier que l'audit est actif
	activeAudits := orchestrator.ListActiveAudits()
	if len(activeAudits) != 1 {
		t.Fatalf("Expected 1 active audit before cancel, got %d", len(activeAudits))
	}
	
	// Annuler l'audit
	err = orchestrator.CancelAudit(request.AuditID)
	if err != nil {
		t.Errorf("CancelAudit should not return error: %v", err)
	}
	
	// Attendre un peu pour que l'annulation prenne effet
	time.Sleep(100 * time.Millisecond)
	
	// Vérifier que l'audit n'est plus actif
	activeAudits = orchestrator.ListActiveAudits()
	if len(activeAudits) != 0 {
		t.Errorf("Expected 0 active audits after cancel, got %d", len(activeAudits))
	}
	
	// Test annulation d'audit inexistant
	err = orchestrator.CancelAudit("non-existing-audit")
	if err == nil {
		t.Error("CancelAudit should return error for non-existing audit")
	}
}

func TestOrchestratorV2_GetResults(t *testing.T) {
	orchestrator := NewOrchestratorV2()
	
	// Test avec audit inexistant
	results, err := orchestrator.GetResults("non-existing-audit")
	if err == nil {
		t.Error("GetResults should return error for non-existing audit")
	}
	if results != nil {
		t.Error("GetResults should return nil for non-existing audit")
	}
}

func TestOrchestratorV2_WithAllAgents(t *testing.T) {
	orchestrator := NewOrchestratorV2()
	
	// Enregistrer tous les agents
	agentMap := map[string]agents.Agent{
		constants.AgentNameKeyword:     keyword.NewKeywordExtractor(),
		constants.AgentNameTechnical:   technical.NewTechnicalAuditor(),
		constants.AgentNameLinking:     linking.NewLinkingMapper(),
		constants.AgentNameBrokenLinks: broken.NewBrokenLinksDetector(),
	}
	
	for name, agent := range agentMap {
		err := orchestrator.RegisterAgent(name, agent)
		if err != nil {
			t.Errorf("Failed to register agent %s: %v", name, err)
		}
	}
	
	request := &AuditRequest{
		AuditID:   "full-pipeline-test",
		SeedURL:   "https://example.com",
		MaxPages:  3,
		Timestamp: time.Now(),
	}
	
	ctx := context.Background()
	progressChan, err := orchestrator.StartAudit(ctx, request)
	if err != nil {
		t.Fatalf("Failed to start audit with all agents: %v", err)
	}
	
	// Collecter les mises à jour de progression
	var updates []*ProgressUpdate
	timeout := time.After(5 * time.Second)
	
	for {
		select {
		case update, ok := <-progressChan:
			if !ok {
				// Channel fermé, audit terminé
				goto AnalyzeUpdates
			}
			updates = append(updates, update)
		case <-timeout:
			t.Fatal("Audit timed out")
		}
	}
	
AnalyzeUpdates:
	if len(updates) == 0 {
		t.Fatal("Should receive at least one progress update")
	}
	
	// Vérifier que nous avons reçu des updates pour différentes étapes
	steps := make(map[string]bool)
	for _, update := range updates {
		steps[update.Step] = true
	}
	
	expectedSteps := []string{
		constants.PipelineStepCrawling,
		constants.PipelineStepAnalyzing,
	}
	
	for _, expectedStep := range expectedSteps {
		if !steps[expectedStep] {
			t.Errorf("Missing progress update for step: %s", expectedStep)
		}
	}
}

func TestOrchestratorV2_Shutdown(t *testing.T) {
	orchestrator := NewOrchestratorV2()
	
	// Démarrer quelques audits
	orchestrator.RegisterAgent(constants.AgentNameKeyword, keyword.NewKeywordExtractor())
	
	for i := 0; i < 3; i++ {
		request := &AuditRequest{
			AuditID:   fmt.Sprintf("shutdown-test-%d", i),
			SeedURL:   "https://example.com",
			MaxPages:  5,
			Timestamp: time.Now(),
		}
		
		ctx := context.Background()
		_, err := orchestrator.StartAudit(ctx, request)
		if err != nil {
			t.Errorf("Failed to start audit %d: %v", i, err)
		}
	}
	
	// Vérifier qu'il y a des audits actifs
	activeAudits := orchestrator.ListActiveAudits()
	if len(activeAudits) != 3 {
		t.Fatalf("Expected 3 active audits before shutdown, got %d", len(activeAudits))
	}
	
	// Arrêt propre
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	err := orchestrator.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown should not return error: %v", err)
	}
	
	// Vérifier que tous les audits ont été arrêtés
	activeAudits = orchestrator.ListActiveAudits()
	if len(activeAudits) != 0 {
		t.Errorf("Expected 0 active audits after shutdown, got %d", len(activeAudits))
	}
}