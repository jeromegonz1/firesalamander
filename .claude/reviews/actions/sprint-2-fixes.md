# Actions Sprint 2 - Corrections Post-Revue

## 🔴 P1 - Corrections Critiques (Bloquant)

### 1. Corriger les erreurs de compilation
**Fichier**: `internal/seo/analyzer.go:174`
**Action**: Résoudre les références aux constantes non définies
```go
// Lignes 174-179 - Définir ou importer les constantes manquantes
constants.TitleMinLength
constants.TitleMaxLength  
constants.MetaDescMinLength
constants.MetaDescMaxLength
constants.MinContentWords
constants.OptimalContentWords
```

### 2. Résoudre le panic runtime Orchestrator
**Fichier**: `internal/orchestrator/orchestrator.go:130`
**Action**: Ajouter vérifications nil avant déréférencement
```go
// Ligne 130 - Ajouter nil checks
if o.CrawlerConfig == nil {
    return fmt.Errorf("crawler config is nil")
}
```

### 3. Réparer les tests cassés
**Fichier**: `cmd/server/main_test.go:97`
**Action**: Implémenter ou supprimer les fonctions non définies
- `resultsHandler`
- `setupServer` 
- `constants.TestQueryURLParam`

### 4. Définir les interfaces core
**Nouveau fichier**: `internal/interfaces/interfaces.go`
```go
type PageCrawler interface {
    Crawl(ctx context.Context, seedURL string, outputDir string) (*CrawlResult, error)
}

type TechnicalAnalyzer interface {
    Analyze(crawlResult *CrawlResult, auditID string) (*TechResult, error)
}

type ReportGenerator interface {
    Generate(results AuditResults, format string) (string, error)
}
```

## 🟠 P2 - Améliorations Importantes

### 5. Refactoriser Orchestrator.runAudit()
**Action**: Diviser en étapes pipeline
- `prepareCrawl()`
- `executeCrawl()`
- `runTechnicalAnalysis()`
- `runSemanticAnalysis()`
- `generateReport()`

### 6. Externaliser le hardcoding
**Configuration**: Créer `config/limits.yaml`
```yaml
crawler:
  default_port: 8080
  max_urls: 300
  semantic_service_url: "http://localhost:5000"
```

### 7. Augmenter couverture tests
**Objectif**: Atteindre 85% pour tous les modules
- Crawler: 46.1% → 85% (+38.9%)
- Audit: 69.0% → 85% (+16%)
- Semantic: 71.7% → 85% (+13.3%)

### 8. Intégrer gofmt pre-commit
**Fichier**: `.pre-commit-config.yaml`
```yaml
repos:
  - repo: local
    hooks:
      - id: gofmt
        name: gofmt
        entry: gofmt -w
        language: system
        files: \.go$
```

## 🟡 P3 - Optimisations

### 9. Réduire complexité Python
**Fichier**: `semantic_analyzer.py`
**Action**: Décomposer les fonctions complexes (39 → <10)

### 10. Diviser les gros fichiers
**Cibles**:
- `internal/seo/analyzer.go` (794 lignes → <300)
- `internal/seo/technical_auditor.go` (748 lignes → <300)

## Timeline et Ressources

| Priorité | Effort Estimé | Assigné | Deadline |
|----------|---------------|---------|----------|
| P1 | 3-5 jours | Dev Lead | Sprint 2.1 |
| P2 | 2-3 jours | Team | Sprint 2.2 |
| P3 | 1-2 jours | Junior | Sprint 2.3 |

## Critères de Validation

### Tests d'Acceptation P1
- [ ] `go build` réussit sans erreurs
- [ ] `go test ./...` passe à 100%
- [ ] Aucun panic en exécution normale
- [ ] Interfaces définies et utilisées

### Métriques P2  
- [ ] Couverture >85% tous modules
- [ ] Zéro hardcoding détecté
- [ ] `gofmt -l` retourne vide
- [ ] Complexité cyclomatique <15

### Validation P3
- [ ] Aucun fichier >500 lignes
- [ ] Complexité Python <10
- [ ] Documentation techniques à jour

---
**Next Review**: Post Sprint 2.1 - Validation corrections P1