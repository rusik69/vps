# VPS Services Unified Makefile

# Variables
DOCKER_REGISTRY ?= 
VERSION ?= latest
JWT_SECRET ?= your-secret-key-change-in-production
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= postgres
POSTGRES_DB ?= shared_db
SERVER_HOST ?= root@your-server.com
PROJECT_DIR ?= /root/vps

.PHONY: help build test run stop clean deploy deploy-prod backup-db logs

# Default target
help:
	@echo "VPS Services - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  build            - Build all Docker images"
	@echo "  test             - Run tests for all services"
	@echo "  run-shortener    - Start shortener service locally"
	@echo "  run-yt           - Start YouTube clone locally"
	@echo "  stop             - Stop all services"
	@echo "  clean            - Clean up containers and images"
	@echo ""
	@echo "Deployment:"
	@echo "  deploy           - Deploy all services to production (alias for deploy-prod)"
	@echo "  deploy-prod      - Deploy all services to production"
	@echo ""
	@echo "Maintenance:"
	@echo "  backup-db        - Backup PostgreSQL database"
	@echo "  logs             - View production service logs"
	@echo "  status           - Check production service status"
	@echo "  restart          - Restart production services"
	@echo ""
	@echo "Variables:"
	@echo "  DOCKER_REGISTRY  - Docker registry URL (default: '')"
	@echo "  VERSION          - Image version (default: latest)"
	@echo "  JWT_SECRET       - JWT secret key (default: your-secret-key-change-in-production)"
	@echo "  SERVER_HOST      - Production server host (default: root@your-server.com)"

# Build all Docker images
build:
	@echo "Building all service images..."
	docker build -t $(DOCKER_REGISTRY)shortener:$(VERSION) ./shortener
	docker build -t $(DOCKER_REGISTRY)yt-backend:$(VERSION) ./yt/backend
	docker build -t $(DOCKER_REGISTRY)yt-frontend:$(VERSION) ./yt/frontend
	@echo "All builds completed successfully!"

# Run tests for all services
test:
	@echo "Running shortener tests..."
	cd shortener && go test ./...
	@echo "Running YouTube clone backend tests..."
	cd yt/backend && go test ./...
	@echo "All tests completed!"

# Start shortener service locally
run-shortener:
	@echo "Starting shortener service locally..."
	cd shortener && make run

# Start YouTube clone locally
run-yt:
	@echo "Starting YouTube clone locally..."
	cd yt && make run

# Stop all services
stop:
	@echo "Stopping all services..."
	cd shortener && make stop || true
	cd yt && make stop || true
	docker-compose -f docker-compose.prod.yml down || true
	@echo "All services stopped!"

# Clean up containers and images
clean: stop
	@echo "Cleaning up containers and images..."
	cd shortener && make clean || true
	cd yt && make clean || true
	docker-compose -f docker-compose.prod.yml down -v --rmi local || true
	docker system prune -f
	@echo "Cleanup completed!"

# Deploy all services to production (alias for deploy-prod)
deploy: deploy-prod

# Deploy all services to production
deploy-prod: build
	@echo "Deploying all services to production..."
	@echo "Pushing images to registry..."
	docker push $(DOCKER_REGISTRY)shortener:$(VERSION)
	docker push $(DOCKER_REGISTRY)yt-backend:$(VERSION)
	docker push $(DOCKER_REGISTRY)yt-frontend:$(VERSION)
	
	@echo "Syncing files to server..."
	rsync -avz --delete \
		--exclude 'node_modules' \
		--exclude '.git' \
		--exclude 'shortener/main' \
		--exclude 'yt/backend/main' \
		--exclude 'yt/frontend/dist' \
		--exclude 'wireguard/config' \
		. $(SERVER_HOST):$(PROJECT_DIR)/
	
	@echo "Creating external volumes on server..."
	ssh $(SERVER_HOST) 'docker volume create shared_postgres_data || true'
	ssh $(SERVER_HOST) 'docker volume create yt_videos_data || true'
	
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
	@echo "Services should be available at:"
	@echo "  - Shortener: http://url.govno2.cloud"
	@echo "  - YouTube Clone: http://yt.govno2.cloud"


# Backup database
backup-db:
	@echo "Creating database backup..."
	mkdir -p backups
	ssh $(SERVER_HOST) 'docker exec shared-postgres-prod pg_dump -U $(POSTGRES_USER) $(POSTGRES_DB)' | gzip > backups/vps_backup_$$(date +%Y%m%d_%H%M%S).sql.gz
	@echo "Backup created in backups/ directory"

# View production logs
logs:
	@echo "Showing production service logs..."
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml logs -f'

# Check production status
status:
	@echo "Checking production service status..."
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml ps'

# Restart production services
restart:
	@echo "Restarting production services..."
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && \
		JWT_SECRET=$(JWT_SECRET) \
		POSTGRES_USER=$(POSTGRES_USER) \
		POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		POSTGRES_DB=$(POSTGRES_DB) \
		DOCKER_REGISTRY=$(DOCKER_REGISTRY) \
		VERSION=$(VERSION) \
		docker-compose -f docker-compose.prod.yml restart'

# Stop production services
stop-prod:
	@echo "Stopping production services..."
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml down'

# View specific service logs
logs-shortener:
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml logs -f shortener'

logs-yt-backend:
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml logs -f yt-backend'

logs-yt-frontend:
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml logs -f yt-frontend'

logs-nginx:
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml logs -f nginx'

logs-postgres:
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml logs -f postgres'

# Development helpers
dev-db:
	@echo "Starting PostgreSQL for development..."
	docker run --name vps-dev-postgres -d \
		-p 5432:5432 \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-v $(PWD)/shortener/database/schema.sql:/docker-entrypoint-initdb.d/01-shortener-schema.sql:ro \
		-v $(PWD)/yt/database/schema.sql:/docker-entrypoint-initdb.d/02-yt-schema.sql:ro \
		postgres:15-alpine

dev-shortener:
	@echo "Starting shortener in development mode..."
	cd shortener && make dev-backend

dev-yt-backend:
	@echo "Starting YouTube clone backend in development mode..."
	cd yt && make dev-backend

dev-yt-frontend:
	@echo "Starting YouTube clone frontend in development mode..."
	cd yt && make dev-frontend

# SSL certificate setup
setup-ssl:
	@echo "Setting up SSL certificates..."
	ssh $(SERVER_HOST) 'mkdir -p $(PROJECT_DIR)/ssl'
	@echo "SSL directory created. Configure your SSL certificates and update nginx.conf"

# Health checks
health-check:
	@echo "Checking service health..."
	@echo "Testing shortener service..."
	curl -f http://url.govno2.cloud/health || echo "Shortener health check failed"
	@echo "Testing YouTube clone service..."
	curl -f http://yt.govno2.cloud/api/videos || echo "YouTube clone health check failed"


# Monitoring
monitor:
	@echo "Starting monitoring dashboard..."
	@echo "Service status:"
	$(MAKE) status
	@echo ""
	@echo "Recent logs (last 50 lines):"
	ssh $(SERVER_HOST) 'cd $(PROJECT_DIR) && docker-compose -f docker-compose.prod.yml logs --tail=50'