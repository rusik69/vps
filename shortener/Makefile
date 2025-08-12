# URL Shortener Makefile

# Variables
APP_NAME := url-shortener
DOCKER_IMAGE ?= $(APP_NAME)
DOCKER_REGISTRY ?= localhost:5000

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linters"
	@echo "  clean        - Clean build artifacts"
	@echo "  deploy-local - Deploy shortener locally using docker-compose.yml"
	@echo "  stop-local   - Stop local shortener deployment"
	@echo "  backup-db    - Create database backup on production"

.PHONY: build
build: ## Build the application
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

.PHONY: run
run: ## Run the application
	go run ./cmd/$(APP_NAME)

.PHONY: test
test: ## Run tests
	go test ./... -v

.PHONY: lint
lint: ## Run linters
	golangci-lint run ./...

.PHONY: check
check: test lint ## Run all checks

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
.PHONY: deploy-local
deploy-local: ## Deploy shortener locally using docker-compose.yml
	docker compose up -d --build

.PHONY: stop-local
stop-local: ## Stop local shortener deployment
	docker compose down

.PHONY: up
up: ## Start services with docker compose (development)
	docker compose up -d

.PHONY: up-test
up-test: ## Start test services with docker compose
	docker compose -f docker-compose.test.yml up -d


.PHONY: down
down: ## Stop services with docker compose (development)
	docker compose down

.PHONY: down-test
down-test: ## Stop test services with docker compose
	docker compose -f docker-compose.test.yml down


.PHONY: logs
logs: ## Show docker compose logs (development)
	docker compose logs -f

.PHONY: logs-test
logs-test: ## Show test docker compose logs
	docker compose -f docker-compose.test.yml logs -f


.PHONY: restart
restart: down up ## Restart services (development)

.PHONY: restart-test
restart-test: down-test up-test ## Restart test services


# Database targets
.PHONY: db-migrate
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	@go run cmd/migrate/main.go migrate

.PHONY: db-seed
db-seed: ## Seed database with test data
	@echo "Seeding database with test data..."
	@go run cmd/migrate/main.go seed

.PHONY: db-reset
db-reset: ## Reset database (drop and recreate)
	@echo "Resetting database..."
	@go run cmd/migrate/main.go reset

# Cleanup targets
.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	docker system prune -f

.PHONY: test-integration
test-integration: ## Run integration tests in containers
	@echo "Running integration tests in containers..."
	docker compose -f docker-compose.test.yml run --rm test

.PHONY: test-all
test-all: test test-integration ## Run all tests
	@echo "All tests completed"

.PHONY: coverage
coverage: ## Generate test coverage report
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

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
