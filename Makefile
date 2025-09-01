# Fire Salamander Makefile

.PHONY: help validate-schemas test init clean build run dev context

help:
	@echo "Fire Salamander - Available commands:"
	@echo "  make context           - Show project context after auto-compact"
	@echo "  make init              - Initialize project structure"
	@echo "  make validate-schemas  - Validate all JSON schemas"
	@echo "  make test              - Run all tests"
	@echo "  make build             - Build the binary"
	@echo "  make run               - Run the application"
	@echo "  make dev               - Run in development mode"
	@echo "  make clean             - Clean generated files"

context:
	@echo "🦎 Fire Salamander - Contexte Projet"
	@echo "═══════════════════════════════════════"
	@echo ""
	@echo "📋 CONTEXTE CCMP ACTUEL:"
	@echo "────────────────────────"
	@cat .claude/context/current_state.md | head -40
	@echo ""
	@echo "📋 DÉCISIONS TECHNIQUES:"
	@echo "─────────────────────────"
	@cat .claude/context/decisions.md | head -15
	@echo ""
	@echo "📈 DERNIERS COMMITS:"
	@echo "──────────────────────"
	@git log --oneline -5 --color=always
	@echo ""
	@echo "📂 FICHIERS MODIFIÉS:"
	@echo "───────────────────────"
	@git status --porcelain || echo "Aucun fichier modifié"
	@echo ""
	@echo "🧪 ÉTAT DES TESTS:"
	@echo "─────────────────"
	@echo "Go tests:"
	@go test ./internal/crawler ./internal/audit ./internal/orchestrator ./internal/semantic ./internal/report 2>/dev/null | grep -E "(PASS|FAIL|ok)" | tail -10 || echo "❌ Erreur tests Go"
	@echo ""
	@echo "Python tests (Agent Sémantique):"
	@cd internal/semantic/python && source venv/bin/activate && python -m pytest --tb=no -q 2>/dev/null | tail -5 || echo "❌ Agent sémantique non testé"
	@echo ""
	@echo "📊 RÉSUMÉ:"
	@echo "─────────"
	@echo "Architecture: 5 agents implémentés avec TDD"
	@echo "Repository: https://github.com/jeromegonz1/firesalamander"
	@echo "CDC: CDC/v4.1-current.md"
	@echo "Specs: SPECS/functional/full-specifications.md"
	@echo ""

validate-schemas:
	@echo "🔍 Validating JSON schemas..."
	@command -v ajv >/dev/null 2>&1 || { echo "Installing ajv-cli..."; npm install -g ajv-cli; }
	@for schema in SPECS/technical/api-contracts/*.schema.json; do \
		echo "✓ Validating $$schema"; \
		ajv validate -s "$$schema" --strict=false || exit 1; \
	done
	@echo "✅ All schemas valid"

test:
	@echo "🧪 Running tests..."
	@go test -v ./...

build:
	@echo "🔨 Building Fire Salamander..."
	@go build -o fire-salamander ./cmd/server

run: build
	@echo "🔥 Starting Fire Salamander..."
	@./fire-salamander

dev:
	@echo "🔥 Starting Fire Salamander in dev mode..."
	@go run ./cmd/server/main.go

init:
	@echo "🚀 Initializing project structure..."
	@bash scripts/init-project.sh
	@echo "✅ Project structure initialized"

clean:
	@echo "🧹 Cleaning generated files..."
	@rm -f fire-salamander
	@rm -rf audits/*/
	@rm -f *.log
	@echo "✅ Clean complete"

# Development helpers
fmt:
	@go fmt ./...

lint:
	@golangci-lint run

deps:
	@go mod tidy
	@go mod download

# CCPM Commands
.PHONY: ccpm-init ccpm-update ccpm-context ccpm-session

ccpm-init:
	@bash scripts/ccpm-update.sh
	@echo "CCPM initialized"

ccpm-context:
	@bash scripts/ccpm-context.sh

ccpm-session:
	@echo "Creating new session log..."
	@bash scripts/ccpm-update.sh

ccpm-update:
	@echo "Updating current state..."
	@${EDITOR:-nano} .claude/context/current_state.md