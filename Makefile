# Discord AI Tech News Bot Makefile

.PHONY: help build run test lint lint-fix clean install-tools dev

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the application
build: ## Build the application
	@echo "Building application..."
	go build -o bin/discord-bot ./cmd

# Run the application
run: ## Run the application
	@echo "Running application..."
	go run ./cmd

# Run tests
# test: ## Run tests
# 	@echo "Running tests..."
# 	go test -race -coverprofile=coverage.out ./...

# Run tests with coverage
# test-coverage: test ## Run tests and show coverage
# 	@echo "Coverage report:"
# 	go tool cover -html=coverage.out -o coverage.html
# 	@echo "Coverage report generated: coverage.html"

# Install development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cosmtrek/air@latest

# Run linter
lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run --config .golangci.yml

# Run linter with auto-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running linter with auto-fix..."
	golangci-lint run --config .golangci.yml --fix

# Development with hot reload
dev: ## Start development server with hot reload
	@echo "Starting development server with hot reload..."
	air

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Tidy dependencies
tidy: ## Tidy Go modules
	@echo "Tidying dependencies..."
	go mod tidy
	go mod verify

# Run all checks (lint, test, build)
check: lint test build ## Run all checks

# Prepare for commit (format, tidy, check)
pre-commit: fmt tidy check ## Prepare code for commit

# Display project info
info: ## Show project information
	@echo "Discord AI Tech News Bot"
	@echo "Go version: $(shell go version)"
	@echo "Git branch: $(shell git branch --show-current 2>/dev/null || echo 'unknown')"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
