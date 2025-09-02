# Stratégie de Test - Fire Salamander

## 🎯 Vision Test
Tests automatisés robustes garantissant qualité, performance et régression zéro pour les 5 agents Fire Salamander.

## 🏗️ Pyramide de Tests

### 1. Tests Unitaires (Base - 70%)
**Outils**: Go tests, pytest
**Cible**: Couverture ≥ 85% par agent
**Scope**: Fonctions isolées, logique métier

```go
// Exemple Go
func TestNormalizeURL(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"https://example.com/page#section", "https://example.com/page"},
        {"https://example.com/page/?utm_source=test", "https://example.com/page/"},
    }
    // ...
}
```

### 2. Tests Contractuels (Milieu - 20%)
**Outils**: JSON Schema + ajv
**Scope**: Validation API entre agents

```bash
# Validation automatique
make validate-schemas
ajv validate -s contracts/schemas/crawl_result.schema.json -d test-data/crawl_sample.json
```

### 3. Tests d'Intégration (Sommet - 10%)
**Outils**: Go tests + fixtures
**Scope**: Pipeline complet agent→agent

## 🤖 Stratégie par Agent

### Agent Crawler
**Tests unitaires**:
- Normalisation URL
- Extraction contenu HTML
- Respect robots.txt
- Détection langue

**Tests d'intégration**:
- Crawl complet site de test
- Performance avec 100+ pages

**Fixtures**: `test-fixtures/test-site/`

### Agent Technical  
**Tests unitaires**:
- Validation balises SEO
- Calculs scores
- Détection erreurs

**Mocks**: Lighthouse API
**Assertions**: Scores dans ranges attendus

### Agent Semantic
**Tests unitaires** (Python):
- Extraction n-grammes
- Topic modeling
- Ranking keywords

**Dataset**: `data/evaluation/sites-annotations.json`
**Métriques**: Précision ≥ 80% vs golden data

### Agent Report
**Tests unitaires**:
- Génération HTML/JSON/CSV
- Template rendering
- Validation output

**Tests visuels**: Snapshot testing HTML

### Agent Orchestrator
**Tests d'intégration**:
- Pipeline complet
- Gestion erreurs
- Status management

## 📊 Métriques Qualité

### Coverage Cibles
- **Go agents**: ≥ 85%
- **Python semantic**: ≥ 90%
- **Global**: ≥ 85%

### Performance Benchmarks
- **Crawl**: < 2s/page (p95)
- **Technical**: < 500ms/page
- **Semantic**: < 15s/50 pages
- **Report**: < 5s génération

### Reliability
- **Uptime**: 99%+ 
- **Error rate**: < 1%
- **Memory leaks**: 0

## 🔧 Outils & Infrastructure

### Local Development
```bash
make test              # Tous les tests
make test-unit         # Tests unitaires only
make test-integration  # Tests intégration only
make coverage          # Rapport coverage
```

### CI/CD Pipeline
- **GitHub Actions**: Tests automatiques
- **Coverage**: Upload Codecov
- **Quality Gates**: Tests + linting obligatoires

### Test Data Management
- **Fixtures**: Sites HTML statiques
- **Golden data**: Résultats attendus annotés
- **Mocks**: APIs externes (Lighthouse)

## 🚀 Test Automation

### Pre-commit Hooks
- Tests unitaires
- Validation schémas
- Linting

### CI Checks
- Tests Go/Python
- Coverage reports  
- Integration tests
- Performance benchmarks

## 📋 Test Maintenance

### Data Refresh
- Mise à jour fixtures mensuellement
- Refresh golden data si algo change
- Clean test databases

### Test Review
- Review tests avec code
- Refactor tests si besoin
- Monitor test performance