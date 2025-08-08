# Makefile for go-vcard

.PHONY: help test test-coverage test-verbose clean build lint format examples

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Testing
test: ## Run tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-coverage: ## Run tests with coverage report
	go test -cover ./...

test-coverage-html: ## Generate HTML coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code quality
lint: ## Run linter
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	golangci-lint run

format: ## Format code
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	go vet ./...

# Build and examples
build: ## Build the module
	go build ./...

examples: ## Set up examples directory
	cd examples && go mod tidy

run-basic-example: examples ## Run basic example
	cd examples && go run basic/main.go

run-http-example: examples ## Run HTTP example
	cd examples && go run http/main.go

# Cleanup
clean: ## Clean generated files
	rm -f coverage.out coverage.html
	rm -f examples/*.vcf
	find . -name "*.vcf" -delete

# Development
dev-setup: ## Set up development environment
	go mod tidy
	go mod download

# Release tasks
tag: ## Create a new git tag (use: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make tag VERSION=v1.0.0"; exit 1; fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

# Documentation
docs: ## Generate documentation
	@echo "Documentation available at:"
	@echo "  - README.md"
	@echo "  - wiki/ directory"
	@echo "  - https://pkg.go.dev/go.rumenx.com/vcard"

# CI/CD helpers
ci-test: ## Run tests in CI environment
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

ci-lint: ## Run linting in CI environment
	golangci-lint run --timeout=5m

# Dependencies
deps: ## Download dependencies
	go mod download

deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

# Check everything
check: format vet lint test ## Run all checks (format, vet, lint, test)

# Quick development cycle
dev: format vet test ## Quick development cycle (format, vet, test)
