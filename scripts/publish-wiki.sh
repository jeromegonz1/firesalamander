#!/bin/bash
# Publier la documentation sur GitHub Wiki

set -e

WIKI_DIR="../firesalamander.wiki"
REPO_URL="https://github.com/jeromegonz1/firesalamander.wiki.git"

echo "🦎 Fire Salamander - Publication Wiki"
echo "======================================"

# Clone wiki if not exists
if [ ! -d "$WIKI_DIR" ]; then
  echo "📥 Clonage du wiki GitHub..."
  git clone $REPO_URL $WIKI_DIR
else
  echo "📂 Wiki déjà cloné, mise à jour..."
  cd $WIKI_DIR
  git pull origin master
  cd -
fi

echo "📋 Copie des fichiers documentation..."

# Copy main files
cp docs/wiki/Home.md $WIKI_DIR/
cp docs/architecture.md $WIKI_DIR/Architecture.md
cp docs/user-guide/getting-started.md $WIKI_DIR/Installation.md
cp docs/user-guide/first-audit.md $WIKI_DIR/User-Guide.md
cp docs/troubleshooting/faq.md $WIKI_DIR/Troubleshooting.md

# Create API reference page
cat > $WIKI_DIR/API-Reference.md << 'EOF'
# API Reference

## JSON-RPC 2.0 Endpoints

### start_crawl
Démarre l'exploration d'un site.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "start_crawl",
  "params": {
    "audit_id": "FS-001",
    "seed_url": "https://example.fr",
    "max_urls": 300
  },
  "id": "request-1"
}
```

**Response:**
```json
{
  "jsonrpc": "2.0", 
  "result": {
    "audit_id": "FS-001",
    "status": "started",
    "estimated_duration": "5m"
  },
  "id": "request-1"
}
```

## REST API

### POST /api/audit
Démarre un audit complet.

### GET /api/audit/:id/status
Récupère le statut d'un audit.

### GET /api/audit/:id/report/:format
Télécharge le rapport (html/json/csv).

Voir [SPECS/technical/](https://github.com/jeromegonz1/firesalamander/tree/main/SPECS/technical) pour les détails.
EOF

echo "✅ Fichiers copiés avec succès"

# Commit and push to wiki
cd $WIKI_DIR

echo "📤 Commit et push vers GitHub Wiki..."

git add .
git commit -m "📚 Update documentation - $(date '+%Y-%m-%d %H:%M')" || echo "Pas de changements à committer"
git push origin master

echo "🎉 Wiki publié avec succès!"
echo "🌐 Voir: https://github.com/jeromegonz1/firesalamander/wiki"

cd -