# 🔥 Fire Salamander

Outil d'audit SEO automatisé développé par SEPTEO Digital Services.

## 🎯 Objectifs
- Crawler intelligent avec respect du robots.txt
- Audit technique complet (Lighthouse + règles SEO)
- Analyse sémantique française avec suggestions IA
- Rapports PDF et web actionables

## 🏗️ Architecture
Fire Salamander utilise une architecture multi-agents :

```
Orchestrator → Crawler → Technical Analyzer → Semantic Analyzer → Report Engine
```

Consultez [docs/architecture.md](docs/architecture.md) pour le schéma complet avec diagramme Mermaid.

## 🚀 Démarrage rapide

### Prérequis
- Go 1.22.5+
- Node.js (pour validation des schémas)
- SQLite

### Installation
```bash
git clone https://github.com/jeromegonz1/firesalamander.git
cd firesalamander
make init
make build
```

### Lancement
```bash
make run
# Ou en mode développement :
make dev
```

L'interface sera accessible sur http://localhost:8080

## 📋 Commandes disponibles

```bash
make help              # Afficher l'aide
make validate-schemas  # Valider les contrats JSON
make test              # Lancer les tests
make build             # Compiler l'application
make run               # Démarrer l'application
make dev               # Mode développement
make clean             # Nettoyer les fichiers générés
```

## 📁 Structure du projet

```
fire-salamander/
├── CDC/                     # Cahier des charges
│   ├── v4.1-current.md     # Version actuelle
│   └── decisions/           # ADR (Architecture Decision Records)
├── SPECS/                   # Spécifications
│   ├── functional/          # Specs agents
│   ├── technical/           # Architecture et contrats API
│   └── test-scenarios/      # Scénarios de test
├── config/                  # Configuration centralisée
├── templates/               # Templates HTML/CSS
├── internal/                # Code Go (privé)
├── cmd/server/             # Point d'entrée
└── scripts/                # Scripts utilitaires
```

## 🔧 Configuration

La configuration est centralisée dans `config/` :
- `config.yaml` - Configuration principale
- `crawler.yaml` - Paramètres de crawling
- `semantic.yaml` - Moteur sémantique
- `tech_rules.yaml` - Règles d'audit technique

## 📊 Documentation

- **[CDC V4.1](CDC/v4.1-current.md)** - Cahier des charges complet
- **[Spécifications fonctionnelles](SPECS/functional/)** - Détail des agents
- **[Architecture technique](SPECS/technical/architecture.md)** - Vue d'ensemble
- **[Flux de données](SPECS/technical/data-flow.md)** - Pipeline de traitement

## 🧪 Tests et validation

```bash
make test                    # Tests unitaires
make validate-schemas        # Validation des contrats API
```

## 🎨 Interface utilisateur

L'interface utilise :
- **Templates Go** pour le rendu serveur
- **Alpine.js** pour l'interactivité
- **CSS SEPTEO** pour le branding
- **Font Awesome** pour les icônes

## 🤖 Agents

1. **Orchestrateur** - Coordination du pipeline
2. **Crawler** - Exploration intelligente des sites
3. **Audit Technique** - Analyse SEO et performance (Lighthouse)
4. **Analyse Sémantique** - Compréhension métier et suggestions
5. **Reporting** - Génération de rapports PDF/web

## 📚 Documentation

- [Architecture du système](docs/architecture.md) - Vue d'ensemble avec schéma Mermaid
- [Spécifications fonctionnelles](SPECS/functional/) - User stories et scénarios BDD
- [Spécifications techniques](SPECS/technical/) - API contracts et intégration
- [Guide de développement](CONTRIBUTING.md) - Standards et processus agile

## 📈 Roadmap

- **Phase 0** ✅ - Structure et fondations
- **Phase 1** 🔄 - Crawler et audit technique
- **Phase 2** 📅 - Analyse sémantique baseline  
- **Phase 3** 📅 - ML et IA enrichie

## 🏢 SEPTEO Digital Services

Développé par l'équipe SEPTEO pour automatiser les audits SEO internes.

---
*Powered by SEPTEO* 🔥