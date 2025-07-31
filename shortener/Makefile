# URL Shortener Makefile

# Variables
APP_NAME := url-shortener
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_REGISTRY := ghcr.io/rusik69
GO_VERSION := 1.21

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development targets
.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod tidy

.PHONY: build
build: ## Build the application
	CGO_ENABLED=0 GOOS=linux go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

.PHONY: build-local
build-local: ## Build the application for local development
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

.PHONY: run
run: build-local ## Run the application locally
	./bin/$(APP_NAME)

# Testing targets
.PHONY: test
test: ## Run tests
	go test ./... -v -race -coverprofile=coverage.out

.PHONY: test-short
test-short: ## Run short tests
	go test ./... -short -v

.PHONY: coverage
coverage: test ## Generate test coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Linting targets
.PHONY: lint
lint: ## Run linters
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
	goimports -w .

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: check
check: fmt vet lint test ## Run all checks (format, vet, lint, test)

# Docker targets
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run: docker-build ## Run Docker container
	docker run -p 8080:8080 $(DOCKER_IMAGE)

.PHONY: docker-push
docker-push: docker-build ## Push Docker image to registry
	docker tag $(DOCKER_IMAGE) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)

# Docker Compose targets
.PHONY: up
up: ## Start services with docker-compose
	docker-compose up -d

.PHONY: down
down: ## Stop services with docker-compose
	docker-compose down

.PHONY: logs
logs: ## Show docker-compose logs
	docker-compose logs -f

.PHONY: restart
restart: down up ## Restart services

# Database targets
.PHONY: db-migrate
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	# Add migration command here when implemented

.PHONY: db-seed
db-seed: ## Seed database with test data
	@echo "Seeding database..."
	# Add seed command here when implemented

# Cleanup targets
.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	docker system prune -f

.PHONY: clean-all
clean-all: clean ## Clean everything including Docker images
	docker rmi $(DOCKER_IMAGE) 2>/dev/null || true

# CI/CD targets
.PHONY: ci
ci: deps check ## Run CI pipeline
	@echo "CI pipeline completed successfully"

.PHONY: release
release: ci docker-build docker-push ## Build and release
	@echo "Release completed successfully"
