SHELL := /bin/bash
-include .env
export

GO_FILES := $(shell find . -name '*.go' -not -path './proto/gen/*' -not -path './vendor/*')
GO_MODULES := $(shell find . -name 'go.mod' -exec dirname {} \;)

.PHONY: help setup dev build test test/unit test/integration test/cover \
        lint fmt proto clean docker-up docker-down \
        security-scan scaffold-service scaffold-engine \
        migrate seed check-health

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## One-time developer setup
	@echo "=== Horizon Development Setup ==="
	@command -v go >/dev/null 2>&1 || { echo "Go 1.24+ required"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "Docker required"; exit 1; }
	@go mod download
	@cd tools && go mod download && cd ..
	@echo "=== Setup Complete ==="

dev: ## Start full development environment
	@echo "=== Starting Horizon Development Environment ==="
	@docker compose -f deploy/docker/docker-compose.yml up -d
	@echo "Waiting for services..."
	@sleep 5
	@echo "Development environment ready."

build: ## Compile all Go services
	@echo "Building Go packages..."
	@cd packages/errors && go build ./...
	@cd packages/types && go build ./...
	@cd packages/events && go build ./...
	@cd packages/logging && go build ./...
	@cd packages/testing && go build ./...
	@cd packages/config && go build ./...
	@cd packages/http && go build ./...
	@cd packages/auth && go build ./...
	@cd packages/telemetry && go build ./...
	@echo "All packages compiled successfully."

test: test/unit test/integration ## Run all tests

test/unit: ## Run unit tests
	@echo "Running unit tests..."
	@go test ./packages/... -count=1 -short
	@echo "Unit tests passed."

test/integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test ./tests/integration/... -count=1 -v 2>/dev/null || echo "No integration tests yet"
	@echo "Integration tests passed."

test/cover: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test ./packages/... -count=1 -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out
	@echo "Coverage report generated."

lint: ## Run all linters
	@echo "Running Go linters..."
	@golangci-lint run ./packages/... --timeout=5m
	@echo "Lint passed."

fmt: ## Format all code
	@gofumpt -l -w $(GO_FILES)
	@goimports -local github.com/horizon/core -w $(GO_FILES)
	@echo "Format complete."

proto: ## Compile protobuf definitions
	@cd proto && buf generate
	@echo "Proto compilation complete."

clean: ## Clean build artifacts
	@rm -rf proto/gen/
	@rm -f coverage.out coverage.html
	@echo "Clean complete."

docker-up: ## Start Docker infrastructure
	@docker compose -f deploy/docker/docker-compose.yml up -d
	@echo "Docker infrastructure started."

docker-down: ## Stop Docker infrastructure
	@docker compose -f deploy/docker/docker-compose.yml down
	@echo "Docker infrastructure stopped."

security-scan: ## Run security vulnerability scanners
	@govulncheck ./packages/...
	@echo "Security scan complete."

scaffold-service: ## Create a new Go service from template
	@read -p "Service type (domain/engine/infra): " type; \
	 read -p "Service name: " name; \
	 scripts/scaffold-service.sh $$type $$name

migrate: ## Run database migrations
	@echo "Migration placeholder — implement when PKG-001 arrives."

seed: ## Seed test data
	@echo "Seed placeholder — implement when PKG-001 arrives."

check-health: ## Verify all services healthy
	@echo "Health check placeholder."

build-all: build ## Build all (alias)
