# URL Shortener Service

A high-performance URL shortener service built with Go, featuring:
- Persistent PostgreSQL database
- Modern web interface
- Rate limiting
- CAPTCHA protection
- Docker containerization
- GitHub Actions CI/CD

## Features

- Generate short URLs from long URLs
- Redirect to original URLs
- Rate limiting per IP
- CAPTCHA protection for bots
- Analytics tracking
- RESTful API
- Web dashboard

## Prerequisites

- Docker and Docker Compose
- Go 1.21+
- PostgreSQL 15+

## Local Development

1. Clone the repository
2. Run development environment:
```bash
docker-compose up -d
```

3. Access the application:
- Web interface: http://localhost:8080
- API: http://localhost:8080/api

## Building and Testing

Run tests:
```bash
go test ./...
```

Build Docker image:
```bash
docker build -t url-shortener .
```

## Project Structure

```
.
├── cmd/              # Main application entry points
├── internal/         # Internal packages
│   ├── api/         # REST API handlers
│   ├── db/          # Database operations
│   ├── middleware/  # HTTP middleware
│   └── service/     # Business logic
├── pkg/             # Public packages
├── web/             # Frontend code
├── .github/         # GitHub Actions workflows
└── docker/          # Docker configuration
```
