# VPS Projects

A collection of projects designed to run on VPS servers. This repository includes multiple microservices, starting with a URL shortener service.

## URL Shortener Service

A modern URL shortener service with the following features:
- Persistent database storage
- Simple yet attractive web interface
- Rate limiting
- CAPTCHA protection
- Docker packaging
- Automated testing and CI/CD

## Project Structure
```
vps/
├── url-shortener/           # Main URL shortener service
│   ├── cmd/                 # Application entry points
│   ├── internal/            # Internal packages
│   │   ├── api/            # API handlers
│   │   ├── db/             # Database operations
│   │   ├── middleware/     # HTTP middleware
│   │   └── web/            # Web interface
│   ├── pkg/                # Reusable packages
│   ├── tests/              # Integration tests
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── go.mod
└── .github/                # GitHub Actions workflows
```

## Prerequisites
- Go 1.21+
- Docker and Docker Compose
- PostgreSQL (for development)

## Getting Started

1. Clone the repository
2. Build and run using Docker:
```bash
docker-compose up --build
```

3. Access the web interface at http://localhost:8080

## Testing

Run tests locally:
```bash
go test ./...
```

## License

MIT License
