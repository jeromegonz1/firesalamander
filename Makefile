# Fire Salamander Makefile

.PHONY: help validate-schemas test init clean build run dev

help:
	@echo "Fire Salamander - Available commands:"
	@echo "  make init              - Initialize project structure"
	@echo "  make validate-schemas  - Validate all JSON schemas"
	@echo "  make test              - Run all tests"
	@echo "  make build             - Build the binary"
	@echo "  make run               - Run the application"
	@echo "  make dev               - Run in development mode"
	@echo "  make clean             - Clean generated files"

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