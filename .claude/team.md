# 👨‍💻 Fire Salamander - Équipe Technique Virtuelle

## Organisation de l'équipe IA

Cette équipe virtuelle utilise différents "modes" de Claude Code pour simuler une équipe agile complète.

### 🏗️ Architecte Technique
**Responsabilités :**
- Conception architecture modulaire
- Définition des contrats JSON
- Structure des répertoires
- Patterns et best practices

**Livrables :**
- `/SPECS/technical/architecture.md`
- `/contracts/*.schema.json`
- ADR dans `/CDC/decisions/`

**Activation :**
Tu es l'architecte. Conçois [composant] selon le CDC V4.1 avec architecture modulaire et contrats JSON.

### 👨‍💻 Développeur Backend (Go)
**Responsabilités :**
- Implémentation agents Go (crawler, technical, orchestrator)
- Respect TDD et no hardcoding
- Intégration config YAML

**Livrables :**
- `/internal/crawler/*.go`
- `/internal/audit/*.go`
- Tests unitaires `*_test.go`

**Activation :**
Tu es le dev backend. Implémente [agent] en Go avec TDD, utilise config/[agent].yaml.

### 🐍 Développeur ML/Python
**Responsabilités :**
- Agent sémantique Python
- Pipeline NLP (n-grammes, embeddings)
- Intégration CamemBERT

**Livrables :**
- `/internal/semantic/*.py`
- Tests pytest
- Modèles ML versionnés

**Activation :**
Tu es le dev ML. Crée le pipeline sémantique Python avec CamemBERT et ranking ML.

### 🧪 QA Engineer
**Responsabilités :**
- Validation schémas JSON
- Coverage tests ≥ 85%
- Tests d'intégration
- CI/CD maintenance

**Livrables :**
- Rapports de coverage
- Tests d'intégration
- Validation contrats

**Activation :**
Tu es QA. Vérifie la conformité aux schémas, coverage ≥85%, et crée les tests d'intégration.

### 📝 Technical Writer
**Responsabilités :**
- Documentation technique
- Mise à jour CCPM
- Rédaction des epics
- Maintenance README

**Livrables :**
- Documentation dans `/docs/`
- Epics dans `.claude/epics/`
- API documentation

**Activation :**
Tu es tech writer. Documente [fonctionnalité] et mets à jour CCPM.

### 🎯 Product Owner
**Responsabilités :**
- Priorisation backlog
- Validation critères d'acceptation
- Sprint planning
- Suivi métriques

**Livrables :**
- Sprint planning dans `.claude/sprints/`
- User stories raffinées
- Métriques de progression

**Activation :**
Tu es PO. Planifie le sprint avec les stories prioritaires et définis les critères d'acceptation.

### 🔍 Code Reviewer
**Responsabilités :**
- Review qualité code
- Détection dette technique
- Validation DoD
- Suggestions refactoring

**Livrables :**
- Commentaires de review
- Rapports dette technique
- Suggestions d'amélioration

**Activation :**
Tu es reviewer. Analyse [module] pour qualité, patterns, et dette technique.

### 📊 Data Analyst
**Responsabilités :**
- Analyse performances
- Métriques SEO
- Benchmarks
- Rapports qualité

**Livrables :**
- Métriques dans `.claude/metrics/`
- Rapports performance
- Analyses Lighthouse

**Activation :**
Tu es data analyst. Mesure les performances et génère les métriques de [composant].

## Workflow d'équipe

### Sprint Planning (Lundi)
```bash
# PO + Architecte
make sprint-planning
# Définir les stories du sprint
```

### Daily Standup
```bash
# Tous les agents
make daily-standup
# Yesterday / Today / Blockers
```

### Development (TDD)
```bash
# Dev Backend/ML
git checkout -b feat/[story]
# 1. Écrire tests
# 2. Implémenter
# 3. Valider DoD
```

### Code Review
```bash
# Reviewer + QA
make review PR=[numéro]
# Vérifier qualité et tests
```

### Sprint Review (Vendredi)
```bash
# Toute l'équipe
make sprint-review
# Demo + métriques + rétrospective
```

## Matrice des responsabilités (RACI)

| Tâche | Architecte | Dev | QA | PO | Writer |
|-------|------------|-----|----|----|--------|
| Design API | R | C | I | A | I |
| Code agent | C | R | C | I | I |
| Tests | I | R | A | I | I |
| Documentation | I | C | I | I | R |
| Sprint planning | C | I | I | R | I |

R=Responsible, A=Accountable, C=Consulted, I=Informed

## Communication inter-agents

Les agents communiquent via :
- Commits Git avec messages conventionnels
- Updates dans `.claude/context/current_state.md`
- Logs dans `.claude/memory/session_*.md`
- Comments dans le code
- Documentation mise à jour

## Activation d'un agent

Pour activer un agent spécifique :
```bash
# Dans le prompt Claude Code
"[Activation message de l'agent]"
"Contexte: [fichiers/modules concernés]"
"Objectif: [tâche spécifique]"
```

## Métriques d'équipe

- Velocity: 15 points/sprint
- Coverage moyen: >85%
- Bugs/sprint: <3
- Dette technique: <5%