package v2

import (
	"fmt"
	"sync"

	"firesalamander/internal/agents"
)

// agentRegistry implémente l'interface AgentRegistry avec thread-safety
type agentRegistry struct {
	mu     sync.RWMutex
	agents map[string]agents.Agent
}

// NewAgentRegistry crée une nouvelle instance d'AgentRegistry
func NewAgentRegistry() AgentRegistry {
	return &agentRegistry{
		agents: make(map[string]agents.Agent),
	}
}

// Register enregistre un agent dans le registry
func (r *agentRegistry) Register(name string, agent agents.Agent) error {
	if name == "" {
		return fmt.Errorf("agent name cannot be empty")
	}
	
	if agent == nil {
		return fmt.Errorf("agent cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[name]; exists {
		return fmt.Errorf("agent %s already registered", name)
	}

	r.agents[name] = agent
	return nil
}

// Get récupère un agent par son nom
func (r *agentRegistry) Get(name string) (agents.Agent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, exists := r.agents[name]
	return agent, exists
}

// List retourne tous les agents enregistrés
func (r *agentRegistry) List() map[string]agents.Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Créer une copie pour éviter les modifications concurrentes
	result := make(map[string]agents.Agent)
	for name, agent := range r.agents {
		result[name] = agent
	}
	return result
}

// HealthCheckAll vérifie la santé de tous les agents
func (r *agentRegistry) HealthCheckAll() map[string]error {
	r.mu.RLock()
	agents := make(map[string]agents.Agent)
	for name, agent := range r.agents {
		agents[name] = agent
	}
	r.mu.RUnlock()

	results := make(map[string]error)
	for name, agent := range agents {
		results[name] = agent.HealthCheck()
	}
	return results
}

// Unregister supprime un agent du registry
func (r *agentRegistry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[name]; !exists {
		return fmt.Errorf("agent %s not found", name)
	}

	delete(r.agents, name)
	return nil
}

// Count retourne le nombre d'agents enregistrés
func (r *agentRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return len(r.agents)
}

// GetStats retourne les statistiques du registry
func (r *agentRegistry) GetStats() *RegistryStats {
	r.mu.RLock()
	agents := make(map[string]agents.Agent)
	for name, agent := range r.agents {
		agents[name] = agent
	}
	r.mu.RUnlock()

	stats := &RegistryStats{
		TotalAgents:     len(agents),
		HealthyAgents:   0,
		UnhealthyAgents: 0,
		AgentsList:      make([]string, 0, len(agents)),
		HealthStatus:    make(map[string]string),
	}

	for name, agent := range agents {
		stats.AgentsList = append(stats.AgentsList, name)
		
		if err := agent.HealthCheck(); err != nil {
			stats.UnhealthyAgents++
			stats.HealthStatus[name] = "unhealthy"
		} else {
			stats.HealthyAgents++
			stats.HealthStatus[name] = "healthy"
		}
	}

	return stats
}