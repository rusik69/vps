# CI/CD Pipeline

This repository uses GitHub Actions for continuous integration and deployment.

## Pipeline Overview

The CI/CD pipeline consists of multiple jobs that run in parallel for efficient testing and building:

### Test Jobs

1. **test-shortener**: Tests and lints the URL shortener service
   - Runs Go tests with PostgreSQL database
   - Performs static code analysis with golangci-lint
   
2. **test-youtube**: Tests and lints the YouTube clone backend
   - Runs comprehensive Go tests including storage, auth, and handler tests
   - Uses GitHub Actions PostgreSQL service for testing
   - Performs static code analysis with golangci-lint
   
3. **test-youtube-frontend**: Tests and builds the YouTube clone frontend
   - Installs Node.js dependencies
   - Runs linting (currently placeholder)
   - Builds the Vue.js frontend

### Build Job

**build**: Builds and pushes Docker images for all services
- Uses matrix strategy to build multiple services in parallel:
  - `url-shortener`: The URL shortener service
  - `yt-backend`: YouTube clone backend API
  - `yt-frontend`: YouTube clone frontend
- Builds multi-architecture images (linux/amd64, linux/arm64)
- Pushes to GitHub Container Registry (ghcr.io)
- Uses Docker layer caching for faster builds

### Deploy Job

**deploy**: Deploys all services to production
- Only runs on manual trigger (`workflow_dispatch`)
- Only runs on the main branch
- Uses unified deployment via the main Makefile
- Deploys all services with shared PostgreSQL database
- Verifies deployment by checking health endpoints

## Triggers

- **Push to main**: Runs tests and builds images
- **Pull Request**: Runs tests only (no build/deploy)
- **Manual Dispatch**: Runs full pipeline including deployment

## Required Secrets

Configure these secrets in your GitHub repository:

- `VPSHOST`: Your production server hostname/IP
- `SSH_PRIVATE_KEY`: SSH private key for server access
- `JWT_SECRET`: Secret key for JWT token signing
- `POSTGRES_PASSWORD`: PostgreSQL database password
- `DOCKER_REGISTRY`: Docker registry URL (optional, defaults to ghcr.io)

## Services Architecture

The pipeline builds and deploys:

1. **URL Shortener** (`url.govno2.cloud`)
   - Go backend with PostgreSQL
   - Nginx routing and static file serving

2. **YouTube Clone** (`yt.govno2.cloud`)
   - Go backend API with file upload support
   - Vue.js frontend with video management
   - Shared PostgreSQL database with separate table prefixes

3. **Shared Infrastructure**
   - Single PostgreSQL database container
   - Nginx reverse proxy with domain-based routing
   - Persistent volumes for data and video files

## Local Development

Run tests locally:

```bash
# Test shortener
cd shortener && make test

# Test YouTube clone
cd yt/backend && go test ./...
cd yt/frontend && npm run build

# Build all services
make build

# Deploy locally
make run-shortener  # or make run-yt
```

## Deployment

The deployment uses the unified docker-compose.prod.yml configuration:

- Single nginx container handling both domains
- Shared PostgreSQL with separate table prefixes
- External volumes for persistence
- Health checks and graceful shutdowns
- Rate limiting and security headers

Manual deployment trigger creates a complete production environment with both services running on their respective domains.