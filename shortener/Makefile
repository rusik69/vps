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
up: ## Start services with docker compose (development)
	docker compose up -d

.PHONY: up-test
up-test: ## Start test services with docker compose
	docker compose -f docker-compose.test.yml up -d

.PHONY: up-prod
up-prod: ## Start production services with docker compose
	docker compose -f docker-compose.prod.yml up -d

.PHONY: down
down: ## Stop services with docker compose (development)
	docker compose down

.PHONY: down-test
down-test: ## Stop test services with docker compose
	docker compose -f docker-compose.test.yml down

.PHONY: down-prod
down-prod: ## Stop production services with docker compose
	docker compose -f docker-compose.prod.yml down

.PHONY: logs
logs: ## Show docker compose logs (development)
	docker compose logs -f

.PHONY: logs-test
logs-test: ## Show test docker compose logs
	docker compose -f docker-compose.test.yml logs -f

.PHONY: logs-prod
logs-prod: ## Show production docker compose logs
	docker compose -f docker-compose.prod.yml logs -f

.PHONY: restart
restart: down up ## Restart services (development)

.PHONY: restart-test
restart-test: down-test up-test ## Restart test services

.PHONY: restart-prod
restart-prod: down-prod up-prod ## Restart production services

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

.PHONY: deploy
deploy: ## Deploy production compose file to msk.govno2.cloud
	@echo "Deploying production compose file to msk.govno2.cloud..."
	@if [ -z "$$SSHPASS" ]; then echo "Error: SSHPASS environment variable is not set"; exit 1; fi
	@echo "Testing SSH connectivity..."
	@sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "echo 'SSH connection successful'" || { echo "Error: Cannot connect to msk.govno2.cloud via SSH"; exit 1; }
	@echo "Creating shortener directory and syncing production compose file..."
	sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "mkdir -p ~/shortener"
	sshpass -e rsync -avz -e "ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no" docker-compose.prod.yml root@msk.govno2.cloud:~/shortener/
	@echo "Pulling latest images and starting services..."
	sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "cd ~/shortener && docker compose -f docker-compose.prod.yml pull && docker compose -f docker-compose.prod.yml up -d"
	@echo "Deployment completed successfully"

.PHONY: stop
stop: ## Stop services on msk.govno2.cloud
	@echo "Stopping services on msk.govno2.cloud..."
	@if [ -z "$$SSHPASS" ]; then echo "Error: SSHPASS environment variable is not set"; exit 1; fi
	@sshpass -e ssh root@msk.govno2.cloud "cd ~/shortener && docker compose -f docker-compose.prod.yml down"
	@echo "Services stopped successfully"

.PHONY: release
release: ci docker-build docker-push ## Build and release
	@echo "Release completed successfully"
