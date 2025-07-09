# Go Microservice Boilerplate

A production-ready Go microservice boilerplate with clean architecture, versioned APIs, RBAC, OIDC, and more.

## Features

### Architecture
- **Clean Architecture** with layers: domain, usecase, delivery (HTTP), infrastructure
- **Dependency Injection** via constructor pattern
- **Interface-first design** for testability

### Core Libraries
- `gin-gonic/gin` for HTTP routing
- `gorm.io/gorm` for ORM with PostgreSQL
- `golang-migrate/migrate` for database migrations
- `spf13/viper` for configuration (YAML + ENV)
- `uber-go/zap` for JSON-only logging
- `casbin/casbin/v2` for RBAC
- `coreos/dex` for OIDC GitHub SSO
- `stretchr/testify` for testing
- `air-verse/air` for live reload
- `swaggo/swag` + `gin-swagger` for OpenAPI docs
- `opentelemetry.io/otel` with `otelgin` and `otelzap` for tracing/logging

### API Features
- Versioned API structure (`/api/v1/`, `/api/v2/`)
- Todo CRUD API (`GET /todos`, `POST /todos`) behind Casbin RBAC
- Role-based access control (user, admin, superadmin)
- Health (`/healthz`) and readiness (`/readyz`) endpoints
- Swagger docs at `/swagger/index.html`

### Infrastructure
- PostgreSQL database with GORM
- Dex OIDC GitHub login with static superadmin email
- Structured JSON logging with trace IDs
- OpenTelemetry tracing
- Docker and Kubernetes deployment
- GitHub Actions CI

## Project Structure

```
go-microservice-boilerplate/
├── cmd/
│   └── server/
│       └── main.go                      # App bootstrap (logger, router, DI)
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       ├── router.go                # Root router + version groups
│   │       └── v1/
│   │       │   ├── routes.go            # Register v1 routes
│   │       │   └── handler/
│   │       │       └── todo_handler.go
│   │       └── v2/
│   │       │   ├── routes.go            # Placeholder for future v2
│   │       │   └── handler/
│   │       │       └── todo_handler.go
│   │       └── middleware/
│   │           ├── auth.go              # Dex OIDC auth middleware
│   │           ├── casbin.go            # RBAC enforcement
│   │           ├── otel.go              # OpenTelemetry tracing
│   │           └── logger.go            # Request logging
│   ├── domain/
│   │   └── todo.go                      # Core entity
│   ├── usecase/
│   │   └── todo_usecase.go              # Business logic
│   └── infrastructure/
│       ├── db/
│       │   ├── postgres.go              # PostgreSQL + GORM
│       │   └── migration.go             # Database migrations
│       ├── dex/
│       │   └── client.go                # Dex OIDC client
│       ├── rbac/
│       │   ├── casbin.go                # Casbin RBAC enforcer
│       │   ├── model.conf               # RBAC model
│       │   └── policy.csv               # RBAC policy
│       ├── logger/
│       │   └── zap.go                   # JSON-only logger via otelzap
│       └── telemetry/
│           └── otel.go                  # Tracer + exporter setup
├── config/
│   ├── config.yaml                      # Default configuration
│   └── config.go                        # Configuration loader
├── api/
│   └── docs/                            # Swagger JSON/YAML
├── migrations/
│   └── 001_init.up.sql                  # Initial schema
├── test/
│   └── todo_test.go                     # Integration tests
├── deploy/
│   └── k8s/
│       ├── deployment.yaml              # Kubernetes deployment
│       ├── service.yaml                 # Kubernetes service
│       ├── dex.yaml                     # Dex deployment
│       ├── configmap.yaml               # ConfigMap
│       └── secret.yaml                  # Secrets
├── .github/
│   └── workflows/
│       └── ci.yaml                      # GitHub Actions CI
├── .air.toml                            # Live reload config
├── .env.example                         # Environment variables
├── .golangci.yml                        # Linter configuration
├── Dockerfile                           # Multi-stage build
├── docker-compose.yml                   # Local development
├── Makefile                             # Build commands
├── go.mod                               # Go modules
└── go.sum                               # Go dependencies
```

## Getting Started

### Prerequisites
- Go 1.21+
- Docker or Podman
- PostgreSQL (or use the provided Docker Compose setup)

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-microservice-boilerplate.git
   cd go-microservice-boilerplate
   ```

2. Create a `.env` file from the example:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. Start the dependencies using Docker Compose:
   ```bash
   docker-compose up -d postgres dex
   ```

4. Run the application with live reload:
   ```bash
   make run
   ```

5. Access the API at http://localhost:8080

### Using Make Commands

- `make run` - Run the application with live reload
- `make build` - Build the application
- `make test` - Run tests
- `make lint` - Run linters
- `make migrate` - Run database migrations
- `make swagger` - Generate Swagger documentation

### Docker

Build and run the application using Docker:

```bash
docker build -t go-microservice-boilerplate .
docker run -p 8080:8080 --env-file .env go-microservice-boilerplate
```

Or use Docker Compose to run the entire stack:

```bash
docker-compose up
```

### Kubernetes

Deploy to Kubernetes:

```bash
kubectl apply -f deploy/k8s/
```

Or use Podman:

```bash
podman kube play deploy/k8s/deployment.yaml
```

## API Documentation

Swagger documentation is available at `/swagger/index.html` when the application is running.

## Authentication

The application uses Dex for OIDC authentication with GitHub. Configure your GitHub OAuth application in the Dex configuration.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
