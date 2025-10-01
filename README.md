# 🚀 Pulse - HTTP Job Scheduler

[![Release](https://img.shields.io/github/v/release/lucasbonna/pulse)](https://github.com/lucasbonna/pulse/releases)
[![Docker](https://img.shields.io/badge/docker-ghcr.io-blue)](https://ghcr.io/lucasbonna/pulse)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Pulse is a lightweight, self-contained HTTP job scheduler written in Go. It allows you to schedule and monitor HTTP requests with configurable intervals, making it perfect for health checks, webhooks, API monitoring, and periodic tasks.

## ✨ Features

- 🔄 **HTTP Job Scheduling**: Schedule GET, POST, PUT, PATCH, DELETE requests
- ⏰ **Flexible Intervals**: Configure execution intervals from 1 second to 24 hours
- 🗄️ **SQLite Database**: Lightweight, embedded database with WAL mode
- 🔒 **Bearer Token Authentication**: Secure API access
- 🐳 **Docker Ready**: Single container deployment
- 📊 **Job Management**: Create, update, delete, and monitor jobs via REST API
- 🚀 **Concurrent Execution**: Parallel job execution with overlap protection
- 📝 **Request Logging**: Track job execution history and status

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐
│   REST API      │    │   Job Scheduler  │    │   SQLite DB
│   (Chi Router)  │◄──►│   (1s interval)  │◄──►│   (WAL mode)
└─────────────────┘    └──────────────────┘
         │                        │
         ▼                        ▼
┌─────────────────┐    ┌──────────────────┐
│  Authentication │    │  HTTP Requests   │
│  (Bearer Token) │    │  (Configurable)  │
└─────────────────┘    └──────────────────┘
```

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
# Pull and run
docker run -d \
  --name pulse \
  -p 8080:8080 \
  -e PORT=8080 \
  -e TOKEN=your_secret_token \
  -v pulse_data:/app/data \
  ghcr.io/lucasbonna/pulse:latest
```

### Using Docker Compose

```yaml
version: '3.8'
services:
  pulse:
    image: ghcr.io/lucasbonna/pulse:latest
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - TOKEN=your_secret_token
    volumes:
      - pulse_data:/app/data
    restart: unless-stopped

volumes:
  pulse_data:
```

### From Source

```bash
# Clone repository
git clone https://github.com/lucasbonna/pulse.git
cd pulse

# Create .env file
echo "PORT=:8080" > .env
echo "TOKEN=your_secret_token" >> .env

# Install dependencies
go mod tidy

# Run the application
go run cmd/server/main.go
```

## 📖 API Documentation

### Authentication

All API requests require a Bearer token in the Authorization header:

```bash
Authorization: Bearer your_secret_token
```

### Base URL

```
http://localhost:8080/api
```

### Endpoints

#### Create Job
```http
POST /api/jobs
Content-Type: application/json
Authorization: Bearer your_secret_token

{
  "name": "Health Check API",
  "url": "https://api.example.com/health",
  "method": "GET",
  "headers": "User-Agent: Pulse/1.0",
  "interval_seconds": 300,
  "active": true
}
```

#### List All Jobs
```http
GET /api/jobs
Authorization: Bearer your_secret_token
```

#### Update Job
```http
PATCH /api/jobs/{id}
Content-Type: application/json
Authorization: Bearer your_secret_token

{
  "name": "Updated Health Check",
  "interval_seconds": 600,
  "active": false
}
```

#### Delete Job
```http
DELETE /api/jobs/{id}
Authorization: Bearer your_secret_token
```

### Request/Response Examples

**Create Job Response:**
```json
{
  "id": 1,
  "name": "Health Check API",
  "url": "https://api.example.com/health",
  "method": "GET",
  "headers": "User-Agent: Pulse/1.0",
  "interval_seconds": 300,
  "next_run_at": "2024-01-15T10:30:00Z",
  "active": true
}
```

**Job List Response:**
```json
[
  {
    "id": 1,
    "name": "Health Check API",
    "url": "https://api.example.com/health",
    "method": "GET",
    "interval_seconds": 300,
    "next_run_at": "2024-01-15T10:30:00Z",
    "active": true
  }
]
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | Yes |
| `TOKEN` | Bearer token for API auth | - | Yes |

### Job Configuration

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `name` | string | Job identifier | 1-100 chars |
| `url` | string | Target URL | Valid URL |
| `method` | string | HTTP method | GET, POST, PUT, PATCH, DELETE |
| `headers` | string | HTTP headers (optional) | Max 1000 chars |
| `interval_seconds` | int | Execution interval | 1-86400 (1s to 24h) |
| `active` | bool | Job status | true/false |

## 🐳 Docker Deployment

### Single Container

```bash
docker run -d \
  --name pulse \
  --restart unless-stopped \
  -p 8080:8080 \
  -e PORT=8080 \
  -e TOKEN=super_secret_token \
  -v /opt/pulse/data:/app/data \
  ghcr.io/lucasbonna/pulse:latest
```

### Docker Swarm Stack

```yaml
version: '3.8'

services:
  pulse:
    image: ghcr.io/lucasbonna/pulse:latest
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - TOKEN=super_secret_token
    volumes:
      - pulse_data:/app/data
    deploy:
      replicas: 1
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
      resources:
        limits:
          memory: 128M
        reservations:
          memory: 64M
    healthcheck:
      test: ["CMD", "timeout", "3", "sh", "-c", "nc -z localhost 8080"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  pulse_data:
    driver: local
```

Deploy with:
```bash
docker stack deploy -c docker-stack.yml pulse
```

## 📊 Monitoring & Logs

### View Logs
```bash
# Docker container
docker logs -f pulse

# Docker Swarm
docker service logs -f pulse_pulse
```

### Database Location
- **Container**: `/app/data/db.sqlite`
- **Host volume**: Mounted directory + `/db.sqlite`

### Log Levels
- Job execution status
- HTTP request results
- Scheduler activities
- API requests (with Chi middleware)

## 🔄 How It Works

1. **Scheduler**: Runs every 1 second, checks for due jobs
2. **Job Execution**: Makes HTTP requests based on job configuration
3. **Concurrency**: Prevents overlapping executions of the same job
4. **Next Run**: Calculates next execution time after completion
5. **Persistence**: Stores jobs and history in SQLite database

## 🛠️ Development

### Prerequisites
- Go 1.21+
- SQLite (handled by modernc.org/sqlite)

### Local Development
```bash
# Clone repository
git clone https://github.com/lucasbonna/pulse.git
cd pulse

# Install dependencies
go mod tidy

# Create .env file
cp .env.example .env

# Install Air for hot reload (optional)
go install github.com/cosmtrek/air@latest

# Run with hot reload
air

# Or run directly
go run cmd/server/main.go
```

### Database Migrations
SQLite schema is embedded in the binary and auto-applied on startup.

### Generate DB Code
```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate after modifying queries
sqlc generate
```

## 📝 Use Cases

- **API Health Monitoring**: Regular health checks for your services
- **Webhook Scheduling**: Periodic webhook deliveries
- **Data Synchronization**: Trigger sync processes at intervals
- **Cache Warming**: Keep caches warm with scheduled requests
- **Status Page Updates**: Update external status pages
- **Backup Triggers**: Initiate backup processes via API calls

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Convention
This project uses [Conventional Commits](https://conventionalcommits.org/):
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks

## 🗺️ Roadmap

We have exciting features planned for future releases! Here's what's coming:

### 🔥 **Phase 1: Core Improvements** (v1.1 - v1.2)

#### **Execution Tracking & Reliability**
- [ ] **Job execution history** - Track last 10-20 runs with status and timing
- [ ] **Configurable timeouts** - Per-job timeout settings (default: 30s)
- [ ] **Headers implementation** - Proper parsing and usage of custom headers
- [ ] **Health check endpoints** - `/health` and `/ready` for better Docker integration

#### **Monitoring & Observability**
- [ ] **Job statistics API** - Success rates, average response times, execution trends
- [ ] **Prometheus metrics** - `/metrics` endpoint for monitoring systems

### 📊 **Phase 2: Advanced Features** (v1.3 - v2.0)

#### **Enhanced Scheduling**
- [ ] **Cron expressions** - Support for complex scheduling beyond simple intervals
- [ ] **One-time jobs** - Execute jobs just once at a specific time

#### **Reliability & Alerts**
- [ ] **Circuit breaker** - Automatically disable failing jobs temporarily
- [ ] **Webhook notifications** - Get alerts when jobs fail consistently
- [ ] **Dead letter queue** - Handle jobs that fail repeatedly

#### **Real-time Updates**
- [ ] **Server-Sent Events** - Live execution notifications

### 🎨 **Phase 3: User Experience** (v2.1+)

#### **Visual Dashboard** (Separate Project)
- [ ] **Web UI** - React-based dashboard consuming the Pulse API
- [ ] **Real-time monitoring** - Live job status and execution graphs
- [ ] **Job management** - Create, edit, and manage jobs through UI
- [ ] **Charts & analytics** - Visual representation of job performance

#### **Developer Tools**
- [ ] **OpenAPI/Swagger** - Auto-generated API documentation
- [ ] **Job templates** - Reusable job configurations

### 🚀 **Phase 4: Enterprise Features** (v3.0+)

#### **Security & Multi-tenancy**
- [ ] **Multiple API keys** - Different permission levels
- [ ] **User authentication** - Job ownership and access control

#### **Advanced Integrations**
- [ ] **Message queue support** - Redis, RabbitMQ integration
- [ ] **External databases** - PostgreSQL, MySQL support for larger deployments

---

**Want to contribute?** Check out our [Contributing Guidelines](#-contributing) and help us build these features!

**Have a feature request?** [Open an issue](https://github.com/lucasbonna/pulse/issues) and let's discuss it!

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Chi Router](https://github.com/go-chi/chi) - Lightweight HTTP router
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) - Pure Go SQLite driver
- [sqlc](https://github.com/sqlc-dev/sqlc) - Type-safe SQL code generation
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading

---

**Made by [Lucas Bonna](https://github.com/lucasbonna)**

*⭐ Star this repository if you find it useful!*
