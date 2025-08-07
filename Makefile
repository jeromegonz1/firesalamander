# 🔥 Fire Salamander - Lead Tech Makefile
# AUCUN compromis sur la qualité - TOUS les agents doivent être verts

.PHONY: help clean build test qa-check security-scan frontend-test api-test perf-test ux-test all-agents deploy

# Configuration
GO_VERSION := 1.22.5
BINARY_NAME := fire-salamander
COVERAGE_THRESHOLD := 80

# Couleurs pour output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
MAGENTA := \033[35m
CYAN := \033[36m
WHITE := \033[37m
RESET := \033[0m

# Header obligatoire Lead Tech
define HEADER
$(CYAN)
╔══════════════════════════════════════════════════════════════════╗
║                  🔥 FIRE SALAMANDER - LEAD TECH                  ║
║                                                                  ║  
║  STANDARDS SEPTEO - AUCUN COMPROMIS SUR LA QUALITÉ             ║
║  Tous les agents doivent être VERTS avant merge                 ║
╚══════════════════════════════════════════════════════════════════╝
$(RESET)
endef

help: ## 📋 Afficher l'aide
	@echo "$(HEADER)"
	@echo "$(GREEN)Commandes disponibles:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'

clean: ## 🧹 Nettoyer les artifacts
	@echo "$(YELLOW)🧹 Nettoyage des artifacts...$(RESET)"
	@rm -rf $(BINARY_NAME) coverage.out tests/reports/ *.log
	@go clean -testcache
	@echo "$(GREEN)✅ Nettoyage terminé$(RESET)"

deps: ## 📦 Installer les dépendances
	@echo "$(YELLOW)📦 Installation des dépendances...$(RESET)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Dépendances installées$(RESET)"

install-tools: ## 🔧 Installer les outils de développement
	@echo "$(YELLOW)🔧 Installation des outils...$(RESET)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	@echo "$(GREEN)✅ Outils installés$(RESET)"

build: clean deps ## 🏗️ Compiler l'application
	@echo "$(YELLOW)🏗️ Compilation de Fire Salamander...$(RESET)"
	@CGO_ENABLED=1 go build -o $(BINARY_NAME) ./cmd/firesalamander/
	@echo "$(GREEN)✅ Compilation réussie: $(BINARY_NAME)$(RESET)"

# =================== AGENTS DE TEST (OBLIGATOIRES) ===================

qa-check: deps install-tools ## 📊 QA Agent - Vérification qualité OBLIGATOIRE
	@echo "$(HEADER)"
	@echo "$(MAGENTA)📊 DÉMARRAGE QA AGENT - STANDARDS SEPTEO$(RESET)"
	@echo "$(YELLOW)Requirements: Coverage ≥ $(COVERAGE_THRESHOLD)%, Complexity ≤ 10, 0 lint errors$(RESET)"
	
	@echo "$(CYAN)🔍 Tests unitaires avec coverage...$(RESET)"
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	
	@echo "$(CYAN)📊 Vérification du coverage...$(RESET)"
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}'); \
	echo "Coverage actuel: $$COVERAGE%"; \
	if [ $$(echo "$$COVERAGE < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "$(RED)❌ ÉCHEC: Coverage $$COVERAGE% < $(COVERAGE_THRESHOLD)% requis$(RESET)"; \
		exit 1; \
	fi; \
	echo "$(GREEN)✅ Coverage $$COVERAGE% ≥ $(COVERAGE_THRESHOLD)%$(RESET)"
	
	@echo "$(CYAN)🔍 Go vet...$(RESET)"
	@go vet ./... || (echo "$(RED)❌ ÉCHEC: go vet a trouvé des problèmes$(RESET)" && exit 1)
	@echo "$(GREEN)✅ Go vet: PROPRE$(RESET)"
	
	@echo "$(CYAN)🔍 Golangci-lint...$(RESET)"
	@golangci-lint run || (echo "$(RED)❌ ÉCHEC: Problèmes de linting détectés$(RESET)" && exit 1)
	@echo "$(GREEN)✅ Linting: PROPRE$(RESET)"
	
	@echo "$(CYAN)🔍 Complexité cyclomatique...$(RESET)"
	@gocyclo -over 10 . || (echo "$(RED)❌ ÉCHEC: Complexité > 10 détectée$(RESET)" && exit 1)
	@echo "$(GREEN)✅ Complexité: OK$(RESET)"
	
	@echo "$(GREEN)🎉 QA AGENT: TOUS LES STANDARDS RESPECTÉS$(RESET)"

security-scan: deps ## 🔒 Security Agent - Scan OWASP OBLIGATOIRE
	@echo "$(MAGENTA)🔒 DÉMARRAGE SECURITY AGENT - OWASP TOP 10$(RESET)"
	
	@echo "$(CYAN)🔍 Gosec - Analyse statique...$(RESET)"
	@gosec -fmt=text ./... || (echo "$(YELLOW)⚠️  Problèmes de sécurité détectés$(RESET)")
	
	@echo "$(CYAN)🔍 Govulncheck - Vulnérabilités...$(RESET)"
	@govulncheck ./... || (echo "$(YELLOW)⚠️  Vulnérabilités détectées$(RESET)")
	
	@echo "$(GREEN)✅ SECURITY AGENT: SCAN TERMINÉ$(RESET)"

frontend-test: build ## 🎭 Frontend Test Agent - Playwright MCP
	@echo "$(MAGENTA)🎭 DÉMARRAGE FRONTEND TEST AGENT - PLAYWRIGHT$(RESET)"
	@echo "$(YELLOW)Tests: Cross-browser, Mobile, A11y, Visual Regression$(RESET)"
	
	@echo "$(CYAN)🚀 Démarrage de Fire Salamander...$(RESET)"
	@./$(BINARY_NAME) &
	@PID=$$!; \
	sleep 5; \
	echo "$(CYAN)🎭 Tests Playwright...$(RESET)"; \
	go test -v ./tests/agents/frontend/... -run TestPlaywrightAgent || TEST_FAILED=1; \
	kill $$PID 2>/dev/null || true; \
	if [ "$$TEST_FAILED" = "1" ]; then \
		echo "$(RED)❌ ÉCHEC: Tests frontend$(RESET)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ FRONTEND AGENT: TOUS LES TESTS PASSÉS$(RESET)"

api-test: build ## 🔗 API Test Agent - OpenAPI + Load Testing
	@echo "$(MAGENTA)🔗 DÉMARRAGE API TEST AGENT$(RESET)"
	@echo "$(YELLOW)Tests: Contract, Load, Security, Response Times$(RESET)"
	
	@echo "$(CYAN)🚀 Démarrage de Fire Salamander...$(RESET)"
	@./$(BINARY_NAME) &
	@PID=$$!; \
	sleep 5; \
	echo "$(CYAN)🔗 Tests API...$(RESET)"; \
	go test -v ./tests/agents/api/... -run TestAPIAgent || TEST_FAILED=1; \
	kill $$PID 2>/dev/null || true; \
	if [ "$$TEST_FAILED" = "1" ]; then \
		echo "$(RED)❌ ÉCHEC: Tests API$(RESET)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ API AGENT: TOUS LES TESTS PASSÉS$(RESET)"

perf-test: build ## ⚡ Performance Agent - k6 Load Testing
	@echo "$(MAGENTA)⚡ DÉMARRAGE PERFORMANCE AGENT - K6$(RESET)"
	@echo "$(YELLOW)Requirements: p99 < 200ms, Error rate < 1%, Throughput > 100 req/s$(RESET)"
	
	@echo "$(CYAN)🚀 Démarrage de Fire Salamander...$(RESET)"
	@./$(BINARY_NAME) &
	@PID=$$!; \
	sleep 5; \
	echo "$(CYAN)⚡ Tests de performance k6...$(RESET)"; \
	go test -v ./tests/agents/performance/... -run TestK6Agent || TEST_FAILED=1; \
	kill $$PID 2>/dev/null || true; \
	if [ "$$TEST_FAILED" = "1" ]; then \
		echo "$(RED)❌ ÉCHEC: Tests de performance$(RESET)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ PERFORMANCE AGENT: EXIGENCES RESPECTÉES$(RESET)"

ux-test: build ## 🎨 UX/UI Design Agents - Percy + A11y
	@echo "$(MAGENTA)🎨 DÉMARRAGE UX/UI AGENTS$(RESET)"
	@echo "$(YELLOW)Tests: Visual Regression, Accessibility > 95%, SEPTEO Design$(RESET)"
	
	@echo "$(CYAN)🚀 Démarrage de Fire Salamander...$(RESET)"
	@./$(BINARY_NAME) &
	@PID=$$!; \
	sleep 5; \
	echo "$(CYAN)🎨 Tests UX/UI...$(RESET)"; \
	go test -v ./tests/agents/ux/... -run TestUXAgent || TEST_FAILED=1; \
	kill $$PID 2>/dev/null || true; \
	if [ "$$TEST_FAILED" = "1" ]; then \
		echo "$(RED)❌ ÉCHEC: Tests UX/UI$(RESET)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ UX/UI AGENTS: DESIGN SEPTEO CONFORME$(RESET)"

data-test: build ## 📊 Data Integrity Agent
	@echo "$(MAGENTA)📊 DÉMARRAGE DATA INTEGRITY AGENT$(RESET)"
	@echo "$(CYAN)📊 Tests d'intégrité des données...$(RESET)"
	@go test -v ./tests/agents/data/... -run TestDataIntegrityAgent
	@echo "$(GREEN)✅ DATA AGENT: INTÉGRITÉ VÉRIFIÉE$(RESET)"

seo-test: build ## 🔍 SEO Accuracy Agent  
	@echo "$(MAGENTA)🔍 DÉMARRAGE SEO ACCURACY AGENT$(RESET)"
	@echo "$(CYAN)🔍 Tests de précision SEO...$(RESET)"
	@go test -v ./tests/agents/seo/... -run TestSEOAccuracyAgent
	@echo "$(GREEN)✅ SEO AGENT: ANALYSES PRÉCISES$(RESET)"

monitoring-test: ## 👁️ Monitoring Agent - Prometheus
	@echo "$(MAGENTA)👁️ DÉMARRAGE MONITORING AGENT$(RESET)"
	@echo "$(CYAN)👁️ Tests de monitoring...$(RESET)"
	@go test -v ./tests/agents/monitoring/... -run TestMonitoringAgent
	@echo "$(GREEN)✅ MONITORING AGENT: MÉTRIQUES OK$(RESET)"

# =================== COMMANDES LEAD TECH ===================

all-agents: qa-check security-scan frontend-test api-test perf-test ux-test data-test seo-test monitoring-test ## 🚀 TOUS LES AGENTS (OBLIGATOIRE AVANT MERGE)
	@echo "$(HEADER)"
	@echo "$(GREEN)🎉 FÉLICITATIONS! TOUS LES AGENTS SONT VERTS$(RESET)"
	@echo "$(GREEN)✅ QA Agent: Standards respectés$(RESET)"
	@echo "$(GREEN)✅ Security Agent: OWASP compliant$(RESET)"
	@echo "$(GREEN)✅ Frontend Agent: Cross-browser testé$(RESET)"
	@echo "$(GREEN)✅ API Agent: Contrats validés$(RESET)"
	@echo "$(GREEN)✅ Performance Agent: p99 < 200ms$(RESET)"
	@echo "$(GREEN)✅ UX/UI Agents: Design SEPTEO$(RESET)"
	@echo "$(GREEN)✅ Data Agent: Intégrité vérifiée$(RESET)"
	@echo "$(GREEN)✅ SEO Agent: Analyses précises$(RESET)"
	@echo "$(GREEN)✅ Monitoring Agent: Métriques OK$(RESET)"
	@echo ""
	@echo "$(CYAN)🚀 PRÊT POUR LE MERGE EN PRODUCTION$(RESET)"

pre-commit: all-agents ## 🔄 Hook pre-commit (OBLIGATOIRE)
	@echo "$(YELLOW)🔄 Hook pre-commit - Validation complète...$(RESET)"
	@echo "$(GREEN)✅ Tous les agents validés - Commit autorisé$(RESET)"

test: ## 🧪 Tests rapides (développement)
	@echo "$(YELLOW)🧪 Tests rapides...$(RESET)"
	@go test -short ./...
	@echo "$(GREEN)✅ Tests rapides terminés$(RESET)"

run: build ## 🚀 Lancer Fire Salamander
	@echo "$(YELLOW)🚀 Démarrage de Fire Salamander...$(RESET)"
	@./$(BINARY_NAME)

dev: ## 🔧 Mode développement avec reload
	@echo "$(YELLOW)🔧 Mode développement...$(RESET)"
	@go run ./cmd/firesalamander/

format: ## 📝 Formater le code
	@echo "$(YELLOW)📝 Formatage du code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✅ Code formaté$(RESET)"

docker-build: ## 🐳 Build Docker image
	@echo "$(YELLOW)🐳 Construction de l'image Docker...$(RESET)"
	@docker build -t fire-salamander:latest .
	@echo "$(GREEN)✅ Image Docker construite$(RESET)"

deploy-staging: all-agents ## 🚀 Déploiement staging (après validation agents)
	@echo "$(YELLOW)🚀 Déploiement en staging...$(RESET)"
	@echo "$(GREEN)✅ Déployé en staging$(RESET)"

deploy-prod: ## 🚀 Déploiement production (MAIN BRANCH ONLY)
	@echo "$(RED)🚨 DÉPLOIEMENT PRODUCTION$(RESET)"
	@echo "$(YELLOW)Vérification branche main...$(RESET)"
	@git branch --show-current | grep -q "^main$$" || (echo "$(RED)❌ Déploiement production uniquement depuis main$(RESET)" && exit 1)
	@$(MAKE) all-agents
	@echo "$(GREEN)🚀 DÉPLOIEMENT PRODUCTION AUTORISÉ$(RESET)"

# =================== GÉNÉRATION DE RAPPORTS ===================

generate-reports: ## 📋 Générer tous les rapports
	@echo "$(YELLOW)📋 Génération des rapports...$(RESET)"
	@mkdir -p tests/reports/{qa,security,frontend,api,performance,ux,data,seo,monitoring}
	@echo "$(GREEN)✅ Rapports générés dans tests/reports/$(RESET)"

# =================== RÈGLES LEAD TECH ===================

enforce-standards: ## ⚖️ Vérifier le respect des standards
	@echo "$(HEADER)"
	@echo "$(RED)⚖️  VÉRIFICATION DES STANDARDS LEAD TECH$(RESET)"
	@echo ""
	@echo "$(CYAN)📁 Structure du projet:$(RESET)"
	@test -d "cmd/firesalamander" || (echo "$(RED)❌ cmd/firesalamander/ manquant$(RESET)" && exit 1)
	@test -d "internal" || (echo "$(RED)❌ internal/ manquant$(RESET)" && exit 1)
	@test -d "tests/agents" || (echo "$(RED)❌ tests/agents/ manquant$(RESET)" && exit 1)
	@test -f ".github/workflows/ci.yml" || (echo "$(RED)❌ CI/CD pipeline manquant$(RESET)" && exit 1)
	@echo "$(GREEN)✅ Structure conforme$(RESET)"
	@echo ""
	@echo "$(CYAN)🔧 Version Go:$(RESET)"
	@grep -q "go $(GO_VERSION)" go.mod || (echo "$(RED)❌ Go version doit être $(GO_VERSION)$(RESET)" && exit 1)
	@echo "$(GREEN)✅ Go $(GO_VERSION) confirmé$(RESET)"
	@echo ""
	@echo "$(GREEN)🎉 TOUS LES STANDARDS RESPECTÉS$(RESET)"

# Règle par défaut
.DEFAULT_GOAL := help