# Sprint 2 - Architecture Feature Agents

## Objectif : Micro-services spécialisés
Transformation de modules monolithiques vers agents spécialisés autonomes.

## Structure Agents

### 1. AG-F01: Keyword Extractor (5 pts)
**Localisation :** `internal/agents/keyword/`
**Mission :** Extraction et classification de mots-clés SEO
**Interface :** KeywordAgent
```go
type KeywordAgent interface {
    ExtractKeywords(content string) (*KeywordResult, error)
    AnalyzeDensity(keywords []string, content string) (*DensityReport, error)
}
```

### 2. AG-F03: Technical Auditor (8 pts)
**Localisation :** `internal/agents/technical/`
**Mission :** Audit technique HTML/CSS/Performance
**Interface :** TechnicalAgent
```go
type TechnicalAgent interface {
    AuditPage(page *PageData) (*TechnicalReport, error)
    ValidateStructure(html string) (*StructureResult, error)
}
```

### 3. AG-F05: Linking Mapper (5 pts)
**Localisation :** `internal/agents/linking/`
**Mission :** Cartographie et analyse des liens internes/externes
**Interface :** LinkingAgent
```go
type LinkingAgent interface {
    MapLinks(crawlResult *CrawlResult) (*LinkMap, error)
    AnalyzeLinkStructure(links []Link) (*LinkAnalysis, error)
}
```

### 4. AG-F04: Broken Links Detector (5 pts)
**Localisation :** `internal/agents/broken/`
**Mission :** Détection et validation des liens brisés
**Interface :** BrokenLinksAgent
```go
type BrokenLinksAgent interface {
    CheckLinks(urls []string) (*BrokenLinksReport, error)
    ValidateLink(url string) (*LinkStatus, error)
}
```

## Orchestrateur V2

### INTEG-001: Orchestrateur V2 (8 pts)
**Localisation :** `internal/orchestrator/v2/`
**Mission :** Coordination des agents avec système d'enregistrement
**Fonctionnalités :**
- Enregistrement dynamique d'agents
- Orchestration pipeline parallèle
- Streaming temps réel des résultats
- Gestion d'erreurs distribuées

```go
type OrchestratorV2 interface {
    RegisterAgent(name string, agent Agent) error
    ExecutePipeline(ctx context.Context, request *AuditRequest) <-chan *ProgressUpdate
    GetAgentStatus(agentName string) (*AgentStatus, error)
}
```

## Interface Web MVP

### UI-001: Interface Web MVP (5 pts)
**Localisation :** `web/`
**Mission :** Interface web pour démarrer audits et voir progrès
**Technologies :** HTML/CSS/JS vanilla
**Fonctionnalités :**
- Formulaire de soumission audit
- Affichage temps réel du progrès
- Visualisation des résultats par agent
- Export des rapports

## Points d'Intégration
- Chaque agent expose une interface commune `Agent`
- Communication via channels Go pour streaming
- Orchestrateur V2 gère le cycle de vie complet
- Interface web communique via WebSocket avec l'orchestrateur

## Priorisation
1. **Architecture fondation** (Architecte)
2. **Agents Backend** (Dev Backend) - Parallélisable
3. **Orchestrateur V2** (Dev Backend) - Dépend des agents
4. **Interface Web** (Dev Frontend) - Parallèle avec orchestrateur
5. **Tests E2E** (QA) - Final