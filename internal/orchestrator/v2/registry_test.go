package v2

import (
	"context"
	"fmt"
	"testing"

	"firesalamander/internal/agents"
	"firesalamander/internal/agents/keyword"
	"firesalamander/internal/agents/technical"
	"firesalamander/internal/constants"
)

// MockAgent pour les tests
type MockAgent struct {
	name        string
	shouldFail  bool
	processData interface{}
}

func (m *MockAgent) Name() string {
	return m.name
}

func (m *MockAgent) Process(ctx context.Context, data interface{}) (*agents.AgentResult, error) {
	m.processData = data
	if m.shouldFail {
		return nil, fmt.Errorf("mock agent error")
	}
	
	return &agents.AgentResult{
		AgentName: m.name,
		Status:    constants.StatusCompleted,
		Data:      map[string]interface{}{"test": "data"},
		Duration:  100,
	}, nil
}

func (m *MockAgent) HealthCheck() error {
	if m.shouldFail {
		return fmt.Errorf("mock agent unhealthy")
	}
	return nil
}

func TestNewAgentRegistry(t *testing.T) {
	registry := NewAgentRegistry()
	
	if registry == nil {
		t.Fatal("NewAgentRegistry should not return nil")
	}
	
	if registry.Count() != 0 {
		t.Errorf("Expected empty registry, got count: %d", registry.Count())
	}
}

func TestAgentRegistry_Register(t *testing.T) {
	registry := NewAgentRegistry()
	agent := &MockAgent{name: "test-agent"}
	
	tests := []struct {
		name        string
		agentName   string
		agent       agents.Agent
		expectError bool
	}{
		{
			name:        "valid agent registration",
			agentName:   "test-agent",
			agent:       agent,
			expectError: false,
		},
		{
			name:        "empty name should fail",
			agentName:   "",
			agent:       agent,
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
			agentName:   "test-agent", // Same as first test
			agent:       agent,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registry.Register(tt.agentName, tt.agent)
			
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestAgentRegistry_Get(t *testing.T) {
	registry := NewAgentRegistry()
	agent := &MockAgent{name: "test-agent"}
	
	// Register an agent first
	err := registry.Register("test-agent", agent)
	if err != nil {
		t.Fatalf("Failed to register agent: %v", err)
	}
	
	tests := []struct {
		name        string
		agentName   string
		expectFound bool
	}{
		{
			name:        "existing agent",
			agentName:   "test-agent",
			expectFound: true,
		},
		{
			name:        "non-existing agent",
			agentName:   "missing-agent",
			expectFound: false,
		},
		{
			name:        "empty name",
			agentName:   "",
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrievedAgent, found := registry.Get(tt.agentName)
			
			if found != tt.expectFound {
				t.Errorf("Expected found=%v, got found=%v", tt.expectFound, found)
			}
			
			if tt.expectFound && retrievedAgent == nil {
				t.Error("Expected agent but got nil")
			}
			
			if tt.expectFound && retrievedAgent != agent {
				t.Error("Retrieved agent is not the same as registered")
			}
		})
	}
}

func TestAgentRegistry_List(t *testing.T) {
	registry := NewAgentRegistry()
	
	// Initially empty
	agents := registry.List()
	if len(agents) != 0 {
		t.Errorf("Expected empty list, got %d agents", len(agents))
	}
	
	// Add some agents
	agent1 := &MockAgent{name: "agent1"}
	agent2 := &MockAgent{name: "agent2"}
	
	registry.Register("agent1", agent1)
	registry.Register("agent2", agent2)
	
	agents = registry.List()
	if len(agents) != 2 {
		t.Errorf("Expected 2 agents, got %d", len(agents))
	}
	
	if _, exists := agents["agent1"]; !exists {
		t.Error("agent1 not found in list")
	}
	
	if _, exists := agents["agent2"]; !exists {
		t.Error("agent2 not found in list")
	}
}

func TestAgentRegistry_HealthCheckAll(t *testing.T) {
	registry := NewAgentRegistry()
	
	healthyAgent := &MockAgent{name: "healthy", shouldFail: false}
	unhealthyAgent := &MockAgent{name: "unhealthy", shouldFail: true}
	
	registry.Register("healthy", healthyAgent)
	registry.Register("unhealthy", unhealthyAgent)
	
	results := registry.HealthCheckAll()
	
	if len(results) != 2 {
		t.Errorf("Expected 2 health check results, got %d", len(results))
	}
	
	if results["healthy"] != nil {
		t.Errorf("Expected healthy agent to be healthy, got error: %v", results["healthy"])
	}
	
	if results["unhealthy"] == nil {
		t.Error("Expected unhealthy agent to have error, got nil")
	}
}

func TestAgentRegistry_Unregister(t *testing.T) {
	registry := NewAgentRegistry()
	agent := &MockAgent{name: "test-agent"}
	
	// Register first
	registry.Register("test-agent", agent)
	
	if registry.Count() != 1 {
		t.Fatalf("Expected 1 agent after registration, got %d", registry.Count())
	}
	
	// Unregister existing agent
	err := registry.Unregister("test-agent")
	if err != nil {
		t.Errorf("Failed to unregister existing agent: %v", err)
	}
	
	if registry.Count() != 0 {
		t.Errorf("Expected 0 agents after unregistration, got %d", registry.Count())
	}
	
	// Try to unregister non-existing agent
	err = registry.Unregister("missing-agent")
	if err == nil {
		t.Error("Expected error when unregistering non-existing agent")
	}
}

func TestAgentRegistry_Count(t *testing.T) {
	registry := NewAgentRegistry()
	
	if registry.Count() != 0 {
		t.Errorf("Expected 0 agents initially, got %d", registry.Count())
	}
	
	agent1 := &MockAgent{name: "agent1"}
	agent2 := &MockAgent{name: "agent2"}
	
	registry.Register("agent1", agent1)
	if registry.Count() != 1 {
		t.Errorf("Expected 1 agent after first registration, got %d", registry.Count())
	}
	
	registry.Register("agent2", agent2)
	if registry.Count() != 2 {
		t.Errorf("Expected 2 agents after second registration, got %d", registry.Count())
	}
	
	registry.Unregister("agent1")
	if registry.Count() != 1 {
		t.Errorf("Expected 1 agent after unregistration, got %d", registry.Count())
	}
}

func TestAgentRegistry_GetStats(t *testing.T) {
	registry := NewAgentRegistry()
	
	healthyAgent := &MockAgent{name: "healthy", shouldFail: false}
	unhealthyAgent := &MockAgent{name: "unhealthy", shouldFail: true}
	
	registry.Register("healthy", healthyAgent)
	registry.Register("unhealthy", unhealthyAgent)
	
	stats := registry.GetStats()
	
	if stats == nil {
		t.Fatal("GetStats should not return nil")
	}
	
	if stats.TotalAgents != 2 {
		t.Errorf("Expected 2 total agents, got %d", stats.TotalAgents)
	}
	
	if stats.HealthyAgents != 1 {
		t.Errorf("Expected 1 healthy agent, got %d", stats.HealthyAgents)
	}
	
	if stats.UnhealthyAgents != 1 {
		t.Errorf("Expected 1 unhealthy agent, got %d", stats.UnhealthyAgents)
	}
	
	if len(stats.AgentsList) != 2 {
		t.Errorf("Expected 2 agents in list, got %d", len(stats.AgentsList))
	}
	
	if len(stats.HealthStatus) != 2 {
		t.Errorf("Expected 2 health statuses, got %d", len(stats.HealthStatus))
	}
}

func TestAgentRegistry_ThreadSafety(t *testing.T) {
	registry := NewAgentRegistry()
	
	// Test concurrent access
	done := make(chan bool)
	
	// Goroutine 1: Register agents
	go func() {
		for i := 0; i < 10; i++ {
			agent := &MockAgent{name: fmt.Sprintf("agent-%d", i)}
			registry.Register(fmt.Sprintf("agent-%d", i), agent)
		}
		done <- true
	}()
	
	// Goroutine 2: List agents
	go func() {
		for i := 0; i < 10; i++ {
			registry.List()
		}
		done <- true
	}()
	
	// Goroutine 3: Health check
	go func() {
		for i := 0; i < 10; i++ {
			registry.HealthCheckAll()
		}
		done <- true
	}()
	
	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
	
	// Verify final state
	if registry.Count() != 10 {
		t.Errorf("Expected 10 agents after concurrent operations, got %d", registry.Count())
	}
}

func TestAgentRegistry_WithRealAgents(t *testing.T) {
	registry := NewAgentRegistry()
	
	// Test with real agents from our codebase
	keywordAgent := keyword.NewKeywordExtractor()
	technicalAgent := technical.NewTechnicalAuditor()
	
	err := registry.Register(constants.AgentNameKeyword, keywordAgent)
	if err != nil {
		t.Errorf("Failed to register keyword agent: %v", err)
	}
	
	err = registry.Register(constants.AgentNameTechnical, technicalAgent)
	if err != nil {
		t.Errorf("Failed to register technical agent: %v", err)
	}
	
	if registry.Count() != 2 {
		t.Errorf("Expected 2 real agents, got %d", registry.Count())
	}
	
	// Test health checks with real agents
	healthResults := registry.HealthCheckAll()
	
	for name, err := range healthResults {
		if err != nil {
			t.Errorf("Real agent %s failed health check: %v", name, err)
		}
	}
	
	stats := registry.GetStats()
	if stats.HealthyAgents != 2 {
		t.Errorf("Expected 2 healthy real agents, got %d", stats.HealthyAgents)
	}
}