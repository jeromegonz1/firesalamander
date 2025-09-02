# Contributing to Fire Salamander

Bienvenue dans le projet Fire Salamander ! Ce guide explique comment contribuer efficacement.

## 🚀 Quick Start

1. **Setup session**:
   ```bash
   make session-start
   ```

2. **Voir contexte projet**:
   ```bash
   make context
   ```

3. **Créer branche feature**:
   ```bash
   git checkout -b feature/agent-name-functionality
   ```

## 📋 Standards de développement

### 🎯 Definition of Done
Consultez [.claude/standards/definition_of_done.md](.claude/standards/definition_of_done.md) pour les critères obligatoires.

### 📝 Conventions commits
Suivez [.claude/standards/commit_conventions.md](.claude/standards/commit_conventions.md) pour le format conventional commits.

### 🌿 Branching strategy  
Respectez [.claude/standards/branching_rules.md](.claude/standards/branching_rules.md) pour nommer vos branches.

### 🔍 Pull Request checklist
Utilisez [.claude/standards/checklist_pr.md](.claude/standards/checklist_pr.md) avant de créer votre PR.

## 🏗️ Architecture

### 🤖 Agents Fire Salamander
```
internal/
├── crawler/        # Agent exploration web
├── audit/         # Agent audit technique  
├── semantic/      # Agent analyse sémantique
├── report/        # Agent génération rapports
└── orchestrator/  # Agent coordination
```

### 📋 Contraintes techniques
Voir [.claude/context/constraints.md](.claude/context/constraints.md) pour limites performance et qualité.

## 🧪 Tests & Qualité

### Strategy de test
Consultez [docs/test-strategy.md](docs/test-strategy.md) pour comprendre notre approche TDD.

### Commandes utiles
```bash
make test                    # Tous les tests
make validate-schemas        # Validation JSON Schema
make metrics                 # Métriques progression
```

### Coverage requis
- **Go agents**: ≥ 85%
- **Python semantic**: ≥ 90%
- **Global**: ≥ 85%

## 📊 Sprint Planning

### Epics & User Stories
- Consultez `.claude/epics/` pour epics détaillés
- Sprint actuel dans `.claude/sprints/sprint-X/`
- Estimation en story points (1, 2, 3, 5, 8, 13)

### Daily tracking
Utilisez `.claude/sprints/sprint-X/tasks/day-XX.md` pour suivi quotidien.

## 🔧 Environment Setup

### Prérequis
- Go 1.23+
- Python 3.9+
- Node.js 18+ (pour JSON validation)

### Installation
```bash
git clone https://github.com/jeromegonz1/firesalamander
cd firesalamander
make session-start
```

## 🤝 Code Review Process

1. **Self-review**: Checklist DoD + tests passants
2. **PR creation**: Template description complète
3. **Automated checks**: CI doit être vert
4. **Human review**: Focus architecture et logique métier
5. **Merge**: Squash and merge après approbation

## 📞 Support & Questions

- **Documentation**: `.claude/context/` et `docs/`
- **Context recovery**: `make context` après auto-compact
- **Issues**: GitHub Issues avec templates appropriés
- **Specifications**: `CDC/v4.1-current.md` et `SPECS/`

## 🎨 Code Style

### Go
- Format: `go fmt`
- Linting: `golangci-lint`
- Conventions: Effective Go

### Python  
- Format: `black`
- Linting: `flake8`
- Type hints: `mypy`

### Documentation
- Markdown standard
- Français pour specs métier
- Anglais pour technique OK