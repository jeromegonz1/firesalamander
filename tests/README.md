# 🧪 Fire Salamander - Test Infrastructure

Infrastructure de tests automatisés pour Fire Salamander avec approche TDD.

## 📁 Structure

```
tests/
├── unit/           # Tests unitaires Go (>80% coverage)
├── integration/    # Tests d'intégration API/DB
├── e2e/           # Tests E2E Playwright
├── performance/   # Tests de performance k6
└── agents/        # Agents de test automatisés
    ├── qa/        # QA Agent (Go + coverage)
    ├── api/       # API Test Agent
    ├── frontend/  # Frontend Test Agent
    ├── seo/       # SEO Accuracy Agent
    ├── data/      # Data Integrity Agent
    ├── security/  # Security Agent (OWASP ZAP)
    ├── perf/      # Performance Agent
    └── monitor/   # Monitoring Agent
```

## 🤖 Agents de Test

### QA Agent
- **Objectif** : Qualité du code Go avec >80% coverage
- **Outils** : `go test`, `go vet`, `golangci-lint`, `gocov`
- **Automatisation** : CI/CD sur chaque PR

### API Test Agent
- **Objectif** : Tests de contrat, charge et sécurité API
- **Outils** : Postman/Newman, Artillery, OWASP ZAP
- **Coverage** : Tous les endpoints REST

### Frontend Test Agent  
- **Objectif** : Tests E2E de l'interface utilisateur
- **Outils** : Playwright, Visual regression
- **Scénarios** : User journeys complets

### SEO Accuracy Agent
- **Objectif** : Vérifier la précision des analyses SEO
- **Méthodes** : Comparaison avec outils de référence
- **Métriques** : Précision, recall des détections

### Data Integrity Agent
- **Objectif** : Cohérence des données crawlées
- **Vérifications** : Intégrité référentielle, validations
- **Monitoring** : Alertes sur anomalies

### Security Agent
- **Objectif** : Tests de sécurité automatisés
- **Outils** : OWASP ZAP, Bandit, Safety
- **Coverage** : OWASP Top 10

### Performance Agent
- **Objectif** : Tests de performance et scalabilité  
- **Outils** : k6, Apache Bench
- **Métriques** : Latence, throughput, ressources

### Monitoring Agent
- **Objectif** : Health checks et monitoring
- **Checks** : Endpoints, DB, cache, ressources
- **Alerting** : Slack/Email sur problèmes

## 🚀 Commandes

```bash
# Tests complets
make test-all

# Par catégorie
make test-unit
make test-integration
make test-e2e
make test-performance

# Agents spécifiques
make qa-check
make security-scan
make performance-test

# Coverage
make coverage
```

## 📊 Métriques de Qualité

- **Code Coverage** : >80% pour Go
- **API Tests** : 100% endpoints couverts
- **E2E Tests** : User journeys critiques
- **Security** : 0 vulnérabilités High/Critical
- **Performance** : <200ms P95 latency

## 🔄 Intégration CI/CD

Tests automatiques sur :
- ✅ Chaque commit (unit tests)
- ✅ Chaque PR (full test suite)
- ✅ Nightly builds (performance + security)
- ✅ Pre-production (E2E complets)