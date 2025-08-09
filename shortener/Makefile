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
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linters"
	@echo "  clean    - Clean build artifacts"

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
.PHONY: up
up: ## Start services with docker compose
	docker compose up -d

.PHONY: down
down: ## Stop services with docker compose
	docker compose down

.PHONY: logs
logs: ## Show docker compose logs
	docker compose logs -f

.PHONY: restart
restart: down up ## Restart services

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
	docker compose run --rm test

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

.PHONY: deploy
deploy: ## Deploy docker-compose.yml to msk.govno2.cloud
	@echo "Deploying docker-compose.yml to msk.govno2.cloud..."
	@if [ -z "$$SSHPASS" ]; then echo "Error: SSHPASS environment variable is not set"; exit 1; fi
	@echo "Creating shortener directory and syncing docker-compose.yml..."
	sshpass -e ssh root@msk.govno2.cloud "mkdir -p ~/shortener"
	sshpass -e rsync -avz docker-compose.yml root@msk.govno2.cloud:~/shortener/
	@echo "Pulling latest images and starting services..."
	sshpass -e ssh root@msk.govno2.cloud "cd ~/shortener && docker compose pull && docker compose up -d"
	@echo "Deployment completed successfully"

.PHONY: stop
stop: ## Stop services on msk.govno2.cloud
	@echo "Stopping services on msk.govno2.cloud..."
	@if [ -z "$$SSHPASS" ]; then echo "Error: SSHPASS environment variable is not set"; exit 1; fi
	@sshpass -e ssh root@msk.govno2.cloud "cd ~/shortener && docker compose down"
	@echo "Services stopped successfully"

.PHONY: release
release: ci docker-build docker-push ## Build and release
	@echo "Release completed successfully"
