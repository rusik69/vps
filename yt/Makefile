# YouTube Clone Makefile

# Variables
DOCKER_REGISTRY ?= 
VERSION ?= latest
JWT_SECRET ?= your-secret-key-change-in-production
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= postgres
POSTGRES_DB ?= youtube_clone
SERVER_HOST ?= root@your-server.com
PROJECT_DIR ?= /root/yt

.PHONY: help build test run stop clean deploy-local deploy-prod backup-db logs

# Default target
help:
	@echo "YouTube Clone - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  build         - Build Docker images"
	@echo "  test          - Run tests"
	@echo "  run           - Start services locally"
	@echo "  stop          - Stop services"
	@echo "  clean         - Clean up containers and images"
	@echo ""
	@echo "Deployment:"
	@echo "  deploy-local  - Deploy locally with docker-compose"
	@echo "  deploy-prod   - Deploy to production server"
	@echo ""
	@echo "Maintenance:"
	@echo "  backup-db     - Backup PostgreSQL database"
	@echo "  logs          - View service logs"
	@echo ""
	@echo "Variables:"
	@echo "  DOCKER_REGISTRY - Docker registry URL (default: '')"
	@echo "  VERSION         - Image version (default: latest)"
	@echo "  JWT_SECRET      - JWT secret key (default: your-secret-key-change-in-production)"
	@echo "  SERVER_HOST     - Production server host (default: root@your-server.com)"

# Build Docker images
build:
	@echo "Building YouTube Clone images..."
	docker build -t $(DOCKER_REGISTRY)yt-backend:$(VERSION) ./backend
	docker build -t $(DOCKER_REGISTRY)yt-frontend:$(VERSION) ./frontend
	@echo "Build completed successfully!"

# Run tests
test:
	@echo "Running backend tests..."
	cd backend && go test ./...
	@echo "Backend tests completed!"

# Start services locally
run: build
	@echo "Starting YouTube Clone locally..."
	JWT_SECRET=$(JWT_SECRET) \
	POSTGRES_USER=$(POSTGRES_USER) \
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	POSTGRES_DB=$(POSTGRES_DB) \
	docker-compose up -d
	@echo "Services started! Frontend: http://localhost:3000, Backend: http://localhost:8080"

# Stop services
stop:
	@echo "Stopping YouTube Clone services..."
	docker-compose down
	@echo "Services stopped!"

# Clean up containers and images
clean: stop
	@echo "Cleaning up containers and images..."
	docker-compose down -v --rmi local
	docker system prune -f
	@echo "Cleanup completed!"

# Deploy locally (alias for run)
deploy-local: run

# Deploy to production server
deploy-prod: build
	@echo "Deploying YouTube Clone to production..."
	@echo "Pushing images to registry..."
	docker push $(DOCKER_REGISTRY)yt-backend:$(VERSION)
	docker push $(DOCKER_REGISTRY)yt-frontend:$(VERSION)
	
	@echo "Syncing files to server..."
	rsync -avz --delete \
		--exclude 'node_modules' \
		--exclude '.git' \
		--exclude 'backend/main' \
		--exclude 'frontend/dist' \
		. $(SERVER_HOST):$(PROJECT_DIR)/
	
	@echo "Creating external volume on server..."
	ssh $(SERVER_HOST) 'docker volume create yt_postgres_data || true'
	
	@echo "Starting services on production server..."
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && \
		JWT_SECRET=$(JWT_SECRET) \
		POSTGRES_USER=$(POSTGRES_USER) \
		POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		POSTGRES_DB=$(POSTGRES_DB) \
		DOCKER_REGISTRY=$(DOCKER_REGISTRY) \
		VERSION=$(VERSION) \
		docker-compose -f docker-compose.prod.yml up -d'
	
	@echo "Production deployment completed!"
	@echo "Services should be available at http://$(SERVER_HOST)"

# Backup database
backup-db:
	@echo "Creating database backup..."
	mkdir -p backups
	docker exec yt-postgres-prod pg_dump -U $(POSTGRES_USER) $(POSTGRES_DB) | gzip > backups/yt_backup_$$(date +%Y%m%d_%H%M%S).sql.gz
	@echo "Backup created in backups/ directory"

# View logs
logs:
	@echo "Showing service logs..."
	docker-compose logs -f

# Development helpers
dev-backend:
	@echo "Starting backend in development mode..."
	cd backend && \
	DATABASE_URL="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable" \
	JWT_SECRET=$(JWT_SECRET) \
	go run ./cmd

dev-frontend:
	@echo "Starting frontend in development mode..."
	cd frontend && npm run serve

dev-db:
	@echo "Starting PostgreSQL for development..."
	docker run --name yt-dev-postgres -d \
		-p 5432:5432 \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-v $(PWD)/database/schema.sql:/docker-entrypoint-initdb.d/01-schema.sql:ro \
		postgres:15-alpine

# Production maintenance
prod-logs:
	@echo "Showing production logs..."
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml logs -f'

prod-status:
	@echo "Checking production status..."
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml ps'

prod-restart:
	@echo "Restarting production services..."
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && \
		JWT_SECRET=$(JWT_SECRET) \
		POSTGRES_USER=$(POSTGRES_USER) \
		POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		POSTGRES_DB=$(POSTGRES_DB) \
		DOCKER_REGISTRY=$(DOCKER_REGISTRY) \
		VERSION=$(VERSION) \
		docker-compose -f docker-compose.prod.yml restart'

prod-stop:
	@echo "Stopping production services..."
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml down'

# SSL certificate setup (Let's Encrypt)
setup-ssl:
	@echo "Setting up SSL certificates..."
	ssh $(SERVER_HOST) 'mkdir -p $(PROJECT_DIR)/ssl'
	@echo "SSL setup completed! Configure your domain in nginx.conf"