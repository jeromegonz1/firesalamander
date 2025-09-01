# 🔥 Fire Salamander - Infrastructure de Test TDD

Documentation technique des agents de test automatisés implémentés pour Fire Salamander v5.

## 📋 Vue d'ensemble

L'infrastructure TDD de Fire Salamander comprend **8 agents de test automatisés** couvrant tous les aspects du cycle de développement :

1. **QA Agent** (Go) - Qualité du code et couverture
2. **API Test Agent** (JavaScript) - Tests d'API (contrat, charge, sécurité)
3. **Frontend Test Agent** (Playwright) - Tests E2E de l'interface
4. **Security Agent** (Python) - Analyse de sécurité OWASP
5. **Performance Agent** (k6) - Tests de charge et performance
6. **SEO Accuracy Agent** (Python) - Validation SEO et conformité
7. **Data Integrity Agent** (Go) - Intégrité des données
8. **Monitoring Agent** (Python) - Surveillance système et santé

## 🏗️ Architecture

```
tests/
├── agents/
│   ├── qa/                    # QA Agent (Go)
│   │   ├── qa_agent.go
│   │   └── qa_agent_test.go
│   ├── api/                   # API Test Agent (JS)
│   │   ├── api_test_agent.js
│   │   ├── test_runner.js
│   │   └── package.json
│   ├── frontend/              # Frontend Agent (Playwright)
│   │   ├── tests/
│   │   ├── playwright.config.js
│   │   └── package.json
│   ├── security/              # Security Agent (Python)
│   │   └── security_agent.py
│   ├── performance/           # Performance Agent (k6)
│   │   └── k6-load-test.js
│   ├── seo/                   # SEO Agent (Python)
│   │   └── seo_agent.py
│   ├── data/                  # Data Integrity Agent (Go)
│   │   └── data_integrity_agent.go
│   └── monitoring/            # Monitoring Agent (Python)
│       └── monitoring_agent.py
├── reports/                   # Rapports générés
└── scripts/                   # Scripts utilitaires
```

## 🔍 Détail des Agents

### 1. QA Agent (Go)
**Responsabilité** : Validation de la qualité du code Go
- **Couverture de code** : Minimum 80% requis
- **Tests unitaires** : Exécution et validation
- **Analyse statique** : go vet, golangci-lint, gosec
- **Complexité** : Analyse cyclomatique avec gocyclo
- **Rapport** : Score global /100 avec détails

**Usage** :
```bash
cd tests/agents/qa
go test -v
```

### 2. API Test Agent (JavaScript)
**Responsabilité** : Tests complets des API REST
- **Tests de contrat** : Validation des endpoints
- **Tests de charge** : Montée en charge progressive
- **Tests de sécurité** : Injection SQL, XSS, authz
- **Validation** : Codes de retour, temps de réponse
- **Rapport** : JSON avec métriques détaillées

**Usage** :
```bash
cd tests/agents/api
node test_runner.js --url=http://localhost:3000
```

### 3. Frontend Test Agent (Playwright)
**Responsabilité** : Tests E2E de l'interface utilisateur
- **Navigation** : Pages principales et workflows
- **Formulaires** : Validation des interactions
- **Responsive** : Tests multi-navigateurs
- **Accessibilité** : Validation ARIA et contraste
- **Captures** : Screenshots et vidéos des échecs

**Usage** :
```bash
cd tests/agents/frontend
npx playwright test
```

### 4. Security Agent (Python)
**Responsabilité** : Analyse de sécurité OWASP
- **Headers de sécurité** : CSP, HSTS, X-Frame-Options
- **Injections** : Tests SQL injection et XSS
- **Authentification** : Validation des contrôles d'accès
- **SSL/TLS** : Configuration et certificats
- **Rapport** : Score sécurité avec recommandations

**Usage** :
```bash
cd tests/agents/security
python3 security_agent.py --url=http://localhost:3000
```

### 5. Performance Agent (k6)
**Responsabilité** : Tests de performance et charge
- **Scénarios** : Ramping, spike, stress tests
- **Métriques** : Temps de réponse, throughput, erreurs
- **Seuils** : P95 < 2s, taux d'erreur < 10%
- **Monitoring** : Ressources système
- **Rapport** : JSON avec graphiques de performance

**Usage** :
```bash
cd tests/agents/performance
BASE_URL=http://localhost:3000 k6 run k6-load-test.js
```

### 6. SEO Accuracy Agent (Python)
**Responsabilité** : Validation SEO et conformité web
- **Méta-données** : Titles, descriptions, keywords
- **Structure** : Headings H1-H6, liens internes
- **Performance** : Temps de chargement, compression
- **Accessibilité** : Alt text, structure sémantique
- **Sitemap/Robots** : Présence et format
- **Rapport** : Score SEO avec recommandations

**Usage** :
```bash
cd tests/agents/seo
python3 seo_agent.py --url=http://localhost:3000
```

### 7. Data Integrity Agent (Go)
**Responsabilité** : Validation de l'intégrité des données
- **Schéma** : Structure tables et contraintes
- **Cohérence** : Validation des relations
- **Qualité** : Détection des anomalies
- **Performance** : Temps de requête
- **Intégrité référentielle** : FK et orphelins
- **Rapport** : Score data avec issues détaillées

**Usage** :
```bash
cd tests/agents/data
go run data_integrity_agent.go --database=../../../fire_salamander_dev.db
```

### 8. Monitoring Agent (Python)
**Responsabilité** : Surveillance système et applicative
- **Disponibilité** : Uptime et monitoring continu
- **Ressources** : CPU, mémoire, disque, réseau
- **Performance** : Temps de réponse endpoints
- **Santé** : Health checks et alertes
- **Charge légère** : Tests de concurrent users
- **Rapport** : Dashboard temps réel

**Usage** :
```bash
cd tests/agents/monitoring
python3 monitoring_agent.py --url=http://localhost:3000 --duration=300
```

## 🤖 CI/CD Integration

### GitHub Actions
Le pipeline CI/CD est configuré dans `.github/workflows/ci.yml` avec 7 jobs :

1. **qa-tests** : QA Agent + couverture
2. **api-tests** : API Test Agent
3. **frontend-tests** : Playwright E2E
4. **security-tests** : Security Agent
5. **performance-tests** : k6 load tests (main branch uniquement)
6. **build-deploy** : Build et déploiement
7. **generate-report** : Rapport consolidé

### Script Global
Le script `test-agents.sh` lance tous les agents séquentiellement :

```bash
./test-agents.sh
```

## 📊 Rapports

Chaque agent génère des rapports dans `tests/reports/` :

- **JSON** : Données brutes pour intégration
- **HTML** : Rapports visuels pour review
- **Logs** : Traces d'exécution pour debug

### Structure des rapports :
```
tests/reports/
├── qa/                    # Couverture + qualité Go
├── api/                   # Tests API + performance
├── frontend/              # Screenshots + vidéos E2E
├── security/              # Analyse sécurité OWASP
├── performance/           # Métriques k6
├── seo/                   # Audit SEO complet
├── data/                  # Intégrité base de données
├── monitoring/            # Surveillance système
└── test_summary.html      # Rapport global consolidé
```

## 🔧 Configuration

### Dépendances
- **Go 1.22+** : QA Agent, Data Integrity Agent
- **Node.js 18+** : API Agent, Frontend Agent
- **Python 3.11+** : Security Agent, SEO Agent, Monitoring Agent
- **k6** : Performance Agent (optionnel)

### Installation
```bash
# Dépendances Python
pip3 install -r requirements.txt

# Dépendances Node.js
cd tests/agents/api && npm install
cd tests/agents/frontend && npm install

# Playwright browsers
npx playwright install --with-deps
```

## 🚀 Utilisation

### Développement local
1. Démarrer Fire Salamander : `go run main.go`
2. Lancer les agents : `./test-agents.sh`
3. Consulter les rapports : `open tests/reports/test_summary.html`

### CI/CD
Les agents se lancent automatiquement sur :
- **Push** vers `main` ou `dev`
- **Pull Request** vers `main`

### Seuils de qualité
- **QA** : Couverture ≥ 80%, score ≥ 80/100
- **Security** : Score ≥ 70/100
- **Performance** : P95 < 2s, erreurs < 10%
- **SEO** : Score ≥ 70/100

## 📈 Métriques et KPI

### Qualité
- Couverture de code
- Complexité cyclomatique
- Violations lint/security

### Performance
- Temps de réponse (P50, P95, P99)
- Throughput (req/s)
- Taux d'erreur

### Sécurité
- Score OWASP
- Vulnérabilités détectées
- Configuration sécurisée

### SEO/UX
- Score SEO
- Performance web
- Accessibilité

## 🔄 Maintenance

### Mise à jour des agents
1. Modifier le code dans `tests/agents/`
2. Tester localement
3. Commit et push (CI/CD validation automatique)

### Ajout d'un nouvel agent
1. Créer le dossier dans `tests/agents/`
2. Implémenter l'agent
3. Ajouter au script `test-agents.sh`
4. Mettre à jour `.github/workflows/ci.yml`

## 📋 Checklist TDD

Avant la Phase 3 (Analyse Sémantique), vérifier :

- [ ] ✅ QA Agent : Tests Go passent avec >80% couverture
- [ ] ✅ API Agent : Endpoints de base fonctionnels
- [ ] ✅ Frontend Agent : Navigation principale OK
- [ ] ✅ Security Agent : Baseline sécurité établie
- [ ] ✅ Performance Agent : Perf baseline documentée
- [ ] ✅ SEO Agent : Structure HTML basique validée
- [ ] ✅ Data Agent : Schéma base cohérent
- [ ] ✅ Monitoring Agent : Surveillance active

**Statut** : ✅ Infrastructure TDD complète et opérationnelle

---

## 🎯 Prochaines étapes

L'infrastructure TDD étant maintenant complète, l'équipe peut procéder à la **Phase 3 : Module Analyse Sémantique Hybride** avec confiance, sachant que chaque modification sera automatiquement validée par les 8 agents de test.

**Ready for Phase 3! 🚀**