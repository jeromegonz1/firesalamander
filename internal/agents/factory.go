package agents

import (
	"fmt"
)

// AgentFactory gère la création et l'enregistrement des agents
type AgentFactory struct {
	agents map[string]Agent
}

// NewAgentFactory crée une nouvelle factory d'agents
func NewAgentFactory() *AgentFactory {
	return &AgentFactory{
		agents: make(map[string]Agent),
	}
}

// RegisterAgent enregistre un agent dans la factory
func (f *AgentFactory) RegisterAgent(agent Agent) error {
	name := agent.Name()
	if name == "" {
		return fmt.Errorf("agent name cannot be empty")
	}

	if _, exists := f.agents[name]; exists {
		return fmt.Errorf("agent %s already registered", name)
	}

	f.agents[name] = agent
	return nil
}

// GetAgent récupère un agent enregistré par son nom
func (f *AgentFactory) GetAgent(name string) (Agent, error) {
	agent, exists := f.agents[name]
	if !exists {
		return nil, fmt.Errorf("agent %s not found", name)
	}
	return agent, nil
}

// ListAgents retourne la liste de tous les agents enregistrés
func (f *AgentFactory) ListAgents() map[string]Agent {
	agents := make(map[string]Agent)
	for name, agent := range f.agents {
		agents[name] = agent
	}
	return agents
}

// HealthCheckAll vérifie la santé de tous les agents
func (f *AgentFactory) HealthCheckAll() map[string]error {
	results := make(map[string]error)
	for name, agent := range f.agents {
		results[name] = agent.HealthCheck()
	}
	return results
}

// UnregisterAgent supprime un agent de la factory
func (f *AgentFactory) UnregisterAgent(name string) error {
	if _, exists := f.agents[name]; !exists {
		return fmt.Errorf("agent %s not found", name)
	}
	
	delete(f.agents, name)
	return nil
}

// Count retourne le nombre d'agents enregistrés
func (f *AgentFactory) Count() int {
	return len(f.agents)
}