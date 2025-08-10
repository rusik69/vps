# WireGuard VPN Server Makefile

# Variables
SERVICE_NAME := wireguard
COMPOSE_FILE := docker-compose.prod.yml
VOLUME_NAME := wireguard_config_data

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo "WireGuard VPN Server Management"
	@echo "Available targets:"
	@echo "  deploy-local  - Deploy WireGuard locally for testing"
	@echo "  stop-local    - Stop local WireGuard deployment"
	@echo "  deploy        - Deploy WireGuard to production"
	@echo "  stop          - Stop WireGuard on production"
	@echo "  logs          - Show WireGuard logs on production"
	@echo "  status        - Check WireGuard status on production"
	@echo "  get-config    - Download client configs from production"
	@echo "  restart       - Restart WireGuard on production"

.PHONY: deploy-local
deploy-local: ## Deploy WireGuard locally for testing
	@echo "Creating local WireGuard config volume..."
	docker volume create $(VOLUME_NAME) 2>/dev/null || true
	@echo "Starting WireGuard server locally..."
	docker compose -f $(COMPOSE_FILE) up -d
	@echo "WireGuard deployed locally. Check logs with: docker compose -f $(COMPOSE_FILE) logs -f"

.PHONY: stop-local
stop-local: ## Stop local WireGuard deployment
	@echo "Stopping local WireGuard deployment..."
	docker compose -f $(COMPOSE_FILE) down
	@echo "Local WireGuard stopped"

.PHONY: deploy
deploy: ## Deploy WireGuard to production
	@echo "Deploying WireGuard to production..."
	@if [ -z "$$SSHPASS" ]; then echo "Error: SSHPASS environment variable is not set"; exit 1; fi
	@echo "Testing SSH connectivity..."
	@sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "echo 'SSH connection successful'" || { echo "Error: Cannot connect to msk.govno2.cloud via SSH"; exit 1; }
	@echo "Creating wireguard directory and syncing compose file..."
	sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "mkdir -p ~/wireguard"
	sshpass -e rsync -avz -e "ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no" $(COMPOSE_FILE) root@msk.govno2.cloud:~/wireguard/
	@echo "Creating persistent volume and starting WireGuard..."
	sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "cd ~/wireguard && docker volume create $(VOLUME_NAME) 2>/dev/null || true && docker compose -f $(COMPOSE_FILE) pull && docker compose -f $(COMPOSE_FILE) up -d"
	@echo "WireGuard deployment completed successfully"
	@echo "Client configs will be generated automatically. Use 'make get-config' to download them."

.PHONY: stop
stop: ## Stop WireGuard on production
	@echo "Stopping WireGuard on production..."
	@if [ -z "$$SSHPASS" ]; then echo "Error: SSHPASS environment variable is not set"; exit 1; fi
	@sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "cd ~/wireguard && docker compose -f $(COMPOSE_FILE) down"
	@echo "WireGuard stopped successfully"

.PHONY: logs
logs: ## Show WireGuard logs on production
	@echo "Fetching WireGuard logs from production..."
	@if [ -z "$$SSHPASS" ]; then echo "Error: SSHPASS environment variable is not set"; exit 1; fi
	@sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "cd ~/wireguard && docker compose -f $(COMPOSE_FILE) logs -f --tail=50"

.PHONY: status
status: ## Check WireGuard status on production
	@echo "Checking WireGuard status on production..."
	@if [ -z "$$SSHPASS" ]; then echo "Error: SSHPASS environment variable is not set"; exit 1; fi
	@sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "cd ~/wireguard && docker compose -f $(COMPOSE_FILE) ps && echo '--- Container logs (last 10 lines) ---' && docker compose -f $(COMPOSE_FILE) logs --tail=10"

.PHONY: get-config
get-config: ## Download client configs from production
	@echo "Downloading WireGuard client configs..."
	@if [ -z "$$SSHPASS" ]; then echo "Error: SSHPASS environment variable is not set"; exit 1; fi
	@mkdir -p ./configs
	@echo "Listing available peer configs..."
	@sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "find /var/lib/docker/volumes/$(VOLUME_NAME)/_data/peer* -name '*.conf' 2>/dev/null || echo 'No peer configs found yet. Container may still be initializing.'"
	@echo "Downloading all peer configs to ./configs/"
	@sshpass -e rsync -avz -e "ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no" root@msk.govno2.cloud:/var/lib/docker/volumes/$(VOLUME_NAME)/_data/peer*/. ./configs/ 2>/dev/null || echo "Configs not ready yet or path not accessible"
	@echo "Client configs downloaded to ./configs/ directory"
	@echo "You can also get QR codes with: make qr-codes"

.PHONY: qr-codes
qr-codes: ## Show QR codes for client configs
	@echo "Generating QR codes for client configs..."
	@if [ -z "$$SSHPASS" ]; then echo "Error: SSHPASS environment variable is not set"; exit 1; fi
	@sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "cd ~/wireguard && docker compose -f $(COMPOSE_FILE) exec wireguard /app/show-peer 1" 2>/dev/null || echo "QR code generation failed. Make sure WireGuard is running."

.PHONY: restart
restart: ## Restart WireGuard on production
	@echo "Restarting WireGuard on production..."
	@if [ -z "$$SSHPASS" ]; then echo "Error: SSHPASS environment variable is not set"; exit 1; fi
	@sshpass -e ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no root@msk.govno2.cloud "cd ~/wireguard && docker compose -f $(COMPOSE_FILE) restart"
	@echo "WireGuard restarted successfully"

.PHONY: clean
clean: ## Clean local configs and stop services
	@echo "Cleaning up local WireGuard setup..."
	docker compose -f $(COMPOSE_FILE) down 2>/dev/null || true
	rm -rf ./configs
	@echo "Cleanup completed"