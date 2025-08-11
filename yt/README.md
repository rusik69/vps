# YouTube Clone

A lightweight YouTube clone built with Go backend and Vue.js frontend, containerized with Docker.

## Features

- **User Authentication**: JWT-based registration and login
- **Video Management**: Upload, view, edit, and delete videos
- **Video Streaming**: Direct video URL support with thumbnail preview
- **Responsive UI**: Dark theme with Tailwind CSS
- **Docker Support**: Complete containerization for development and production
- **Production Ready**: Nginx reverse proxy with SSL support

## Architecture

- **Backend**: Go with Gorilla Mux router, JWT authentication, PostgreSQL
- **Frontend**: Vue.js 3 with Vue Router, Axios for API calls
- **Database**: PostgreSQL with proper indexing and triggers
- **Reverse Proxy**: Nginx with rate limiting and security headers
- **Containerization**: Docker Compose orchestration

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Make (optional, but recommended)

### Development Setup

1. **Clone and setup**:
   ```bash
   cd /path/to/vps/yt
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Start services**:
   ```bash
   make run
   ```

3. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432

### Production Deployment

1. **Configure environment**:
   ```bash
   cp .env.example .env
   # Set production values, especially JWT_SECRET and database credentials
   ```

2. **Deploy to production**:
   ```bash
   make deploy-prod SERVER_HOST=root@your-server.com
   ```

## Available Commands

```bash
# Development
make build         # Build Docker images
make test          # Run tests
make run           # Start services locally
make stop          # Stop services
make clean         # Clean up containers and images

# Deployment
make deploy-local  # Deploy locally
make deploy-prod   # Deploy to production

# Maintenance
make backup-db     # Backup database
make logs          # View logs
make prod-status   # Check production status
make prod-restart  # Restart production services
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### Videos
- `GET /api/videos` - Get all videos
- `GET /api/videos/{id}` - Get specific video
- `POST /api/videos` - Create new video (auth required)
- `PUT /api/videos/{id}` - Update video (auth required)
- `DELETE /api/videos/{id}` - Delete video (auth required)
- `GET /api/my-videos` - Get user's videos (auth required)

### Health
- `GET /health` - Health check

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Videos Table
```sql
CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    url TEXT NOT NULL,
    thumbnail_url TEXT,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    views INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |
| `POSTGRES_USER` | Database username | `postgres` |
| `POSTGRES_PASSWORD` | Database password | `postgres` |
| `POSTGRES_DB` | Database name | `youtube_clone` |
| `PORT` | Backend port | `8080` |
| `SERVER_HOST` | Production server | `root@your-server.com` |

### Production Considerations

1. **Security**:
   - Change default JWT secret
   - Use strong database passwords
   - Configure SSL certificates
   - Set up proper firewall rules

2. **Performance**:
   - Configure nginx worker processes
   - Adjust rate limiting based on needs
   - Set up database connection pooling
   - Enable proper caching headers

3. **Monitoring**:
   - Check application logs regularly
   - Monitor database performance
   - Set up health check endpoints

## Development

### Backend Development
```bash
# Start PostgreSQL
make dev-db

# Run backend locally
make dev-backend
```

### Frontend Development
```bash
# Install dependencies
cd frontend && npm install

# Start development server
make dev-frontend
```

### Running Tests
```bash
# Backend tests
cd backend && go test ./...

# Or using make
make test
```

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 3000, 5432, 8080 are available
2. **Database connection**: Check PostgreSQL is running and credentials are correct
3. **CORS issues**: Verify API endpoints are properly configured
4. **Docker issues**: Try `make clean` to reset containers

### Logs

```bash
# View all logs
make logs

# Production logs
make prod-logs

# Specific service logs
docker-compose logs frontend
docker-compose logs backend
docker-compose logs postgres
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## License

This project is licensed under the MIT License.