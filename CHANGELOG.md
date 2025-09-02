# Changelog

Toutes les modifications notables de Fire Salamander sont documentées ici.

Le format s'inspire de [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added
- Documentation complète avec architecture Mermaid
- Tests d'intégration E2E complets
- Gestion erreurs avec fallbacks gracieux

## [0.1.0] - 2025-09-02

### Added
- **Pipeline 5 agents** avec communication JSON-RPC 2.0
- **Agent Crawler** (Go) - Exploration BFS avec respect robots.txt/sitemap
- **Agent Technical** (Go) - Analyse SEO + Lighthouse + Core Web Vitals
- **Agent Orchestrator** (Go) - Coordination pipeline avec audit_ids
- **Agent Semantic** (Python) - NLP français avec CamemBERT/DistilCamemBERT
- **Agent Report Engine** (Go) - Génération HTML/JSON/CSV avec branding SEPTEO
- **Configuration YAML** centralisée sans hardcoding
- **Tests TDD complets** - 43 Go + 15 Python + 17 intégration = 75 tests
- **CI/CD GitHub Actions** avec validation pre-commit
- **CCPM intégration** pour mémoire contextuelle entre sessions
- **Architecture modulaire** avec JSON Schema contracts
- **Dataset d'évaluation** 6 sites français annotés
- **Fixtures de test** HTML + robots.txt + sitemap
- **Environnement agile** avec Definition of Done et standards
- **Sprint planning** avec vélocité établie (35-40 pts/sprint)

### Changed
- Architecture refactorisée de monolithe vers multi-agents
- Communication inter-agents via JSON-RPC au lieu de REST
- Configuration externalisée en YAML
- Tests passés de manuels à TDD automatisés

### Fixed
- Cycles d'import entre packages config et crawler
- Variables non utilisées dans les boucles
- Précision flottante dans les tests (assert.InDelta)
- Dépendances manquantes dans legacy code
- Échecs GitHub Actions CI (paths, timeouts, templates)
- Pre-commit hooks trop stricts
- Tokenizer NLTK manquant (fallback regex)
- Fonction template 'mul' non définie

### Security
- Validation entrée utilisateur
- Rate limiting configurable
- Timeouts stricts anti-DoS
- Sandboxing des analyses
- Pas de stockage données personnelles

### Performance
- Crawling concurrent (5 workers par défaut)
- Cache intelligent (TTL configurable)
- Timeout adaptatifs par type de site
- Mode léger DistilCamemBERT pour performance

### Technical Debt
- Suppression 26k lignes de sur-ingénierie v2
- Simplification de main.go complexe
- Archivage code legacy dans archive/
- Nettoyage logs et fichiers temporaires

## [0.0.1] - 2025-08-01

### Added
- Projet initial avec structure basique
- Templates HTML/CSS existants
- Configuration serveur Go basique

## Types de changements
- `Added` pour les nouvelles fonctionnalités
- `Changed` pour les changements aux fonctionnalités existantes
- `Deprecated` pour les fonctionnalités bientôt supprimées
- `Removed` pour les fonctionnalités supprimées
- `Fixed` pour les corrections de bugs
- `Security` pour les corrections de vulnérabilités
- `Performance` pour les améliorations de performance
- `Technical Debt` pour le nettoyage du code