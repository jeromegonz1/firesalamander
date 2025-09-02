# Guide de démarrage rapide

## Prérequis
- Go 1.21+
- Python 3.11+
- Node.js 18+ (pour Lighthouse)
- 4GB RAM minimum

## Installation

### 1. Cloner le repository
```bash
git clone https://github.com/jeromegonz1/firesalamander.git
cd firesalamander
```

### 2. Installer les dépendances
```bash
make install
```

### 3. Configurer
```bash
cp .env.example .env
# Éditer .env avec vos paramètres
```

## Premier audit

### Via l'interface web

1. Lancer le serveur : `make server`
2. Ouvrir http://localhost:8080
3. Entrer l'URL à auditer
4. Cliquer "Analyser"

### Via l'API
```bash
curl -X POST http://localhost:8080/api/audit \
  -H "Content-Type: application/json" \
  -d '{"seed_url": "https://example.com"}'
```

## Structure d'un audit

1. **Crawl** : Exploration du site (300 pages max)
2. **Analyse technique** : Lighthouse + règles SEO
3. **Analyse sémantique** : Extraction mots-clés FR
4. **Rapport** : PDF avec scores et recommandations

## Prochaines étapes

- [Comprendre les rapports](interpreting-reports.md)
- [Configuration avancée](../developer-guide/setup.md)