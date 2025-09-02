# Orchestrateur V2 - Design Architecture

## 🎯 Objectif
Coordonner les 4 agents feature via un système de streaming temps réel, remplaçant les données simulées de l'interface web.

## 🏗️ Architecture SOLID

### Interfaces principales
```go
type OrchestratorV2 interface {
    RegisterAgent(name string, agent agents.Agent) error
    StartAudit(ctx context.Context, request *AuditRequest) (<-chan *ProgressUpdate, error)
    GetAuditStatus(auditID string) (*AuditExecution, error)
    StreamProgress(auditID string) (<-chan *ProgressUpdate, error)
}

type AgentRegistry interface {
    Register(name string, agent agents.Agent) error
    Get(name string) (agents.Agent, bool)
    List() map[string]agents.Agent
    HealthCheckAll() map[string]error
}

type PipelineExecutor interface {
    Execute(ctx context.Context, request *AuditRequest) (<-chan *PipelineResult, error)
    GetProgress(auditID string) (*ExecutionProgress, error)
}
```

### Composants clés

#### 1. Agent Registry
- **Responsabilité** : Gestion centralisée des agents
- **Pattern** : Registry + Factory
- **Thread-safe** : Mutex pour accès concurrent

#### 2. Pipeline Executor  
- **Responsabilité** : Orchestration séquentielle/parallèle
- **Pattern** : Pipeline + Observer
- **Streaming** : Channels Go pour temps réel

#### 3. Progress Manager
- **Responsabilité** : Suivi état et progression
- **Pattern** : State Manager + Event Emitter
- **Persistance** : In-memory avec possibilité extension

## 🔄 Workflow Pipeline

### Étapes séquentielles
1. **Crawling** → Exploration du site (prérequis)
2. **Parallel Analysis** → 4 agents en parallèle :
   - Keyword Extractor
   - Technical Auditor  
   - Linking Mapper
   - Broken Links Detector
3. **Results Aggregation** → Combinaison résultats

### Streaming temps réel
```go
type ProgressUpdate struct {
    AuditID    string                 `json:"audit_id"`
    Step       string                 `json:"step"`
    Progress   float64                `json:"progress"`
    AgentName  string                 `json:"agent_name,omitempty"`
    AgentStatus string                `json:"agent_status,omitempty"`
    Data       map[string]interface{} `json:"data,omitempty"`
    Timestamp  time.Time              `json:"timestamp"`
}
```

## 📊 Intégration agents existants

### Agents à intégrer
- ✅ `keyword_extractor` (internal/agents/keyword)
- ✅ `technical_auditor` (internal/agents/technical)  
- ✅ `linking_mapper` (internal/agents/linking)
- ✅ `broken_links_detector` (internal/agents/broken)

### Agent Factory usage
```go
factory := agents.NewAgentFactory()
factory.RegisterAgent(keyword.NewKeywordExtractor())
factory.RegisterAgent(technical.NewTechnicalAuditor())
factory.RegisterAgent(linking.NewLinkingMapper())
factory.RegisterAgent(broken.NewBrokenLinksDetector())
```

## 🔧 Configuration (constants.go)

### Orchestrateur constants
```go
const (
    // Orchestrator V2 constants
    OrchestratorV2Name = "orchestrator_v2"
    DefaultMaxConcurrentAudits = 5
    DefaultStreamBufferSize = 100
    DefaultAuditTimeout = 30 * time.Minute
    
    // Pipeline steps
    PipelineStepCrawling = "crawling"
    PipelineStepAnalyzing = "analyzing"  
    PipelineStepReporting = "reporting"
    PipelineStepCompleted = "completed"
    
    // Progress thresholds
    ProgressCrawlingComplete = 30.0
    ProgressAnalysisComplete = 80.0
    ProgressReportingComplete = 95.0
    ProgressFullComplete = 100.0
)
```

## 📡 API Integration

### WebSocket endpoint
```
WS /api/v1/audits/{auditId}/stream
```

### REST endpoints
```
POST /api/v1/audits          ← Démarrer audit (existant)
GET  /api/v1/audits/{id}     ← Status audit (existant)  
GET  /api/v1/agents          ← Liste agents enregistrés
GET  /api/v1/agents/health   ← Health check tous agents
```

## 🧪 Stratégie de test TDD

### 1. Tests unitaires (85% couverture)
- `OrchestratorV2` core logic
- `AgentRegistry` thread-safety
- `PipelineExecutor` error handling
- `ProgressManager` state transitions

### 2. Tests d'intégration  
- Pipeline complet avec vrais agents
- Streaming WebSocket
- Timeout et error recovery
- Concurrent audits

### 3. Tests de performance
- Multiple audits simultanés
- Memory leaks detection
- Channel deadlocks prevention

## 🚀 Plan d'implémentation TDD

### Phase 1: Interfaces + Tests
1. Définir toutes les interfaces
2. Écrire tests qui échouent  
3. Implémenter minimal pour compiler

### Phase 2: Core Components  
1. AgentRegistry avec tests
2. PipelineExecutor avec tests
3. ProgressManager avec tests

### Phase 3: Integration
1. OrchestratorV2 orchestration
2. WebSocket streaming
3. API endpoint integration

### Phase 4: QA Gates
1. ✅ Compilation sans erreurs
2. ✅ Tests unitaires 85%+  
3. ✅ Standards respect (pas de hardcoding)
4. ✅ Architecture SOLID validée
5. ✅ Documentation à jour

## 📝 Fichiers à créer

```
internal/orchestrator/v2/
├── orchestrator.go          # Interface principale  
├── orchestrator_test.go     # Tests TDD
├── registry.go              # Agent registry
├── registry_test.go         # Tests registry
├── pipeline.go              # Pipeline executor
├── pipeline_test.go         # Tests pipeline  
├── progress.go              # Progress manager
├── progress_test.go         # Tests progress
└── types.go                 # Types partagés
```

## 🔄 Remplacement simulation MVP

L'Orchestrateur V2 remplacera :
- ✅ `simulateAuditStart()` → `orchestrator.StartAudit()`
- ✅ `simulateAuditResults()` → Stream temps réel
- ✅ `startProgressMonitoring()` → WebSocket streaming
- ✅ Données fake → Vrais agents avec résultats

---

**🎯 Cette architecture respecte :**
- ✅ **SOLID** : Interfaces, SRP, DIP
- ✅ **No hardcoding** : Constantes configurables  
- ✅ **TDD** : Tests d'abord, puis implémentation
- ✅ **Multi-casquettes** : Architecte → Dev Backend → QA
- ✅ **Streaming temps réel** : Channels + WebSocket
- ✅ **Intégration agents** : Factory pattern existant