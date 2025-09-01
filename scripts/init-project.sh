#!/bin/bash
set -e

echo "🚀 Initializing Fire Salamander project structure..."

# Create missing directories
mkdir -p audits
mkdir -p data
mkdir -p logs
mkdir -p config/prompts
mkdir -p tests/{unit,integration}

# Create default config files if missing
if [ ! -f "config/stopwords_fr.txt" ]; then
    echo "Creating French stopwords list..."
    cat > config/stopwords_fr.txt << 'EOF'
le la les un une des du de
et ou à dans sur avec pour par
ce cette ces cet
qui que quoi dont où
il elle ils elles on nous vous
avoir être faire aller venir
très plus moins aussi
ici là maintenant alors
EOF
fi

# Create prompt templates
if [ ! -f "config/prompts/business_keywords.txt" ]; then
    cat > config/prompts/business_keywords.txt << 'EOF'
Tu es un expert SEO français. Analyse ce contenu web et génère 5 expressions de recherche longue traîne (3-5 mots) que les utilisateurs tapent sur Google pour trouver ce type de business.

Contexte métier: {business_context}
Mots-clés existants: {existing_keywords}

Réponds uniquement avec les expressions, une par ligne, sans numérotation.
EOF
fi

# Initialize SQLite database
if [ ! -f "data/firesalamander.db" ]; then
    echo "Creating SQLite database..."
    sqlite3 data/firesalamander.db << 'EOF'
CREATE TABLE audits (
    id TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    config_snapshot TEXT
);

CREATE TABLE feedback (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    audit_id TEXT,
    keyword TEXT,
    rating INTEGER, -- 1 for 👍, -1 for 👎
    reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (audit_id) REFERENCES audits(id)
);
EOF
fi

# Install development dependencies
echo "Installing development tools..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

echo "✅ Project initialization complete!"
echo ""
echo "Next steps:"
echo "  make build    - Build the application"  
echo "  make dev      - Start development server"
echo "  make test     - Run tests"