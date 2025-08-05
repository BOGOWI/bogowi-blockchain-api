# BOGOWI Blockchain API - Go Development Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint
BINARY_NAME=bogowi-api
BINARY_LINUX=$(BINARY_NAME)-linux

# Build info
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT_SHA ?= $(shell git rev-parse HEAD)
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT_SHA) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all build clean test coverage deps format lint security docker help

all: clean deps format lint test build

## Build commands
build: ## Build the binary
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v ./...

build-linux: ## Build for Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BINARY_LINUX) ./...

build-all: ## Build for all platforms
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BINARY_NAME)-linux-amd64 ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BINARY_NAME)-linux-arm64 ./...
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BINARY_NAME)-darwin-amd64 ./...
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BINARY_NAME)-darwin-arm64 ./...
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BINARY_NAME)-windows-amd64.exe ./...

## Test commands
test: ## Run tests
	$(GOTEST) -v -race ./...

test-coverage: ## Run tests with coverage (standard)
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

coverage-enhanced: ## Generate enhanced HTML coverage report with gocov
	@echo "🚀 Generating enhanced coverage report with gocov..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && \
	gocov test ./... | gocov-html > coverage-enhanced.html
	@echo "✅ Enhanced coverage report: coverage-enhanced.html"
	@command -v open >/dev/null 2>&1 && open coverage-enhanced.html || echo "📄 Report ready to view"

coverage-all: ## Generate all coverage reports
	@echo "📊 Generating all coverage reports..."
	$(MAKE) test-coverage
	$(MAKE) coverage-enhanced
	@echo "🎉 All coverage reports generated!"
	@echo "  📄 coverage.html (standard)"
	@echo "  📄 coverage-enhanced.html (enhanced)"

coverage-api: ## Run API tests with coverage
	@echo "🔍 Running API tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage-api.out ./internal/api/...
	$(GOCMD) tool cover -html=coverage-api.out -o coverage-api.html
	@echo "📊 API coverage report generated: coverage-api.html"
	@echo "📈 API coverage summary:"
	@$(GOCMD) tool cover -func=coverage-api.out | grep total

coverage-install: ## Install coverage tools
	@echo "📦 Installing coverage tools..."
	go install github.com/axw/gocov/gocov@latest
	go install github.com/matm/gocov-html/cmd/gocov-html@latest
	@echo "✅ Coverage tools installed!"

## Development commands
clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_LINUX)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) verify

deps-update: ## Update dependencies
	$(GOMOD) tidy
	$(GOGET) -u all

format: ## Format Go code
	$(GOFMT) -s -w .

format-check: ## Check if code is formatted
	@unformatted=$$($(GOFMT) -s -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not formatted:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

lint: ## Run linter
	$(GOLINT) run --timeout=5m

lint-fix: ## Run linter with auto-fix
	$(GOLINT) run --fix --timeout=5m

security: ## Run security scan
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Installing..."; \
		curl -sfL https://raw.githubusercontent.com/securecodewarrior/gosec/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest; \
		gosec ./...; \
	fi

vet: ## Run go vet
	$(GOCMD) vet ./...

## Development workflow
dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	$(GOMOD) download
	@if ! command -v $(GOLINT) >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest; \
	fi
	@echo "Development setup complete!"

pre-commit: format lint vet test ## Run pre-commit checks

ci: deps format-check lint vet test ## Run CI checks

## Docker commands
docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME):latest .

docker-run: ## Run Docker container
	docker run -p 3001:3001 --env-file .env $(BINARY_NAME):latest

## Run commands
run: ## Run the application
	$(GOCMD) run main.go

run-dev: ## Run in development mode with hot reload (requires air)
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Air not installed. Run: go install github.com/cosmtrek/air@latest"; \
		$(GOCMD) run main.go; \
	fi

## Help
help: ## Show this help message
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
