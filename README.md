# Gin Microservice Boilerplate

A production-ready Go microservice boilerplate with clean architecture, versioned APIs, RBAC, OIDC, and more.

## Features

### Architecture
- **Clean Architecture** with layers: domain, usecase, delivery (HTTP), infrastructure
- **Dependency Injection** via constructor pattern
- **Interface-first design** for testability

### Database
- **PostgreSQL** integration with Neon (or any PostgreSQL provider)
- **GORM ORM** for database operations and model management
- **Golang-migrate** for SQL schema migrations
- Database health check via `/readyz` endpoint

### Core Libraries
- `gin-gonic/gin` for HTTP routing
- `gorm.io/gorm` for ORM with PostgreSQL
- `golang-migrate/migrate` for database migrations
- `spf13/viper` for configuration (YAML + ENV)
- `uber-go/zap` for JSON-only logging
- `casbin/casbin/v2` for RBAC
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
- Structured JSON logging with trace IDs
- OpenTelemetry tracing
- Docker and Kubernetes deployment
- GitHub Actions CI

## Project Structure

```
gin-microservice-boilerplate/
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
   git clone https://github.com/yourusername/gin-microservice-boilerplate.git
   cd gin-microservice-boilerplate
   ```

2. Choose your environment configuration:
   ```bash
   # For development environment (default)
   export APP_ENVIRONMENT=dev
   
   # For production environment
   export APP_ENVIRONMENT=prod
   ```

3. Set up your PostgreSQL database:
   - You can use a local PostgreSQL instance for development
   - Or use a cloud provider like Neon (https://neon.tech) for production
   - Configuration is automatically loaded from the appropriate config file (dev.yaml or prod.yaml)

4. Run database migrations:
   ```bash
   make migrate
   ```
   This will apply any pending migrations to your database.

5. Run the application with live reload:
   ```bash
   make run
   ```
   
   The application will:
   - Connect to the PostgreSQL database
   - Start the HTTP server

5. Access the API at http://localhost:8080

   Available endpoints:
   - `/` - Welcome message
   - `/healthz` - Health check endpoint (returns 200 OK if the service is running)
   - `/readyz` - Readiness check endpoint (returns 200 OK if the database connection is healthy, 503 Service Unavailable otherwise)

### Configuration

This project uses a layered configuration system with environment-specific YAML files. The configuration files are located in the `config` directory:

- `config.yaml` - Common configuration shared across all environments
- `dev.yaml` - Development environment-specific configuration
- `prod.yaml` - Production environment-specific configuration

The application first loads the common configuration from `config.yaml`, then merges the environment-specific configuration on top of it based on the `APP_ENVIRONMENT` environment variable. If not set, it defaults to `dev`.

The configuration loading process follows this order of precedence (highest to lowest):

1. Environment variables
2. Environment-specific YAML file (`dev.yaml` or `prod.yaml`)
3. Common configuration file (`config.yaml`)

You can override any configuration value using environment variables. For example:

```bash
export APP_ENVIRONMENT=prod
export DB_HOST=my-custom-db-host.example.com
```

This will load the common configuration, merge the production configuration on top, and then override the `database.host` value with `my-custom-db-host.example.com`.

### Using Make Commands

- `make run` - Run the application with live reload
- `make build` - Build the application
- `make test` - Run tests
- `make lint` - Run linters
- `make migrate` - Run database migrations
- `make swagger` - Generate Swagger documentation

### Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations. 

#### Generating Migration Files

To generate a new migration file:

```bash
# Install golang-migrate CLI if not already installed
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Generate a new migration file
~/go/bin/migrate create -ext sql -dir migrations -seq your_migration_name
```

This will create two files:
- `migrations/NNNNNN_your_migration_name.up.sql` - For applying the migration
- `migrations/NNNNNN_your_migration_name.down.sql` - For reverting the migration

Edit these files to add your SQL statements.

#### Migration Execution

Migrations are handled separately from application startup using the `make migrate` command. This follows best practices for production environments by:

1. Separating database schema changes from application deployment
2. Avoiding race conditions in multi-instance deployments
3. Allowing for controlled migration execution and rollback

To run migrations:

```bash
make migrate
```

This command will:
1. Connect to the database using the configuration from your environment-specific YAML file
2. Apply any pending migrations from the `migrations` directory
3. Log the migration status

### Docker

Build and run the application using Docker:

```bash
docker build -t gin-microservice-boilerplate .
docker run -p 8080:8080 --env-file .env gin-microservice-boilerplate
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

The application uses a simple JWT-based authentication system. This is a stateless authentication mechanism that issues signed JWT tokens containing user identity information.

### Setup

1. **Configure JWT Secret and Expiry**:
   - Edit `config/config.yaml` or set environment variables:
   ```yaml
   auth:
     jwt_secret: "supersecretkey"
     jwt_expiry_hours: 1
     superadmin_email: "admin@example.com"
   ```

2. **Set Superadmin Email** (optional):
   - Edit `config/config.yaml` or set the `SUPERADMIN_EMAIL` environment variable to grant a specific email superadmin privileges

### Authentication Flow

> **Note:** This is a simplified authentication system for development and testing purposes only. In a production environment, this should be replaced with a proper authentication solution that includes secure password handling, user registration, and additional security measures as per your specific requirements.

1. Client sends a POST request to `/auth` with email:
   ```bash
   curl -X POST http://localhost:8080/auth \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com"}'
   ```

2. Server validates the request and returns a JWT token:
   ```json
   {
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
   }
   ```

3. Client includes this token in subsequent API requests

### Protected Endpoints

All API endpoints under `/api/v1/*` require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer your_jwt_token
```

### User Identity

The JWT token contains the following claims:
- `sub`: User's email address (used as the subject identifier)
- `exp`: Token expiration time (default: 1 hour)
- `iat`: Token issue time

The auth middleware injects these values into the Gin context, making them available to handlers via:
- `c.Get("userEmail")` - User's email address
- `c.Get("isSuperAdmin")` - Boolean indicating if user is a superadmin

### Superadmin Access

Users with email matching the `auth.superadmin_email` config value are automatically granted superadmin privileges. This is checked by the auth middleware during token validation.

### Role-Based Access Control (RBAC)

The application uses Casbin for Role-Based Access Control (RBAC) to restrict access to resources based on user roles.

#### RBAC Configuration

1. **Model Definition**:
   - Located at `internal/infrastructure/rbac/model.conf`
   - Defines the RBAC model with subjects (users), objects (resources), and actions (HTTP methods)

2. **Policy Rules**:
   - Located at `internal/infrastructure/rbac/policy.csv`
   - Contains role definitions and permissions in the format:
     - `p, role, resource, action` (permission rule)
     - `g, user_email, role` (role assignment)

3. **Example Policy**:
   ```csv
   p, admin, /api/v1/todos, GET
   p, admin, /api/v1/todos, POST
   p, user, /api/v1/todos, GET
   g, alice@example.com, admin
   g, bob@example.com, user
   ```

#### Access Control Flow

1. **Authentication**: JWT middleware authenticates the user and sets `userEmail` in the context
2. **Authorization**: RBAC middleware checks if the user has permission to access the requested resource
3. **Superadmin Override**: Users with the configured superadmin email bypass RBAC checks
4. **Policy Enforcement**: For regular users, access is granted only if a matching policy rule exists

#### Adding New Roles and Permissions

To add new roles or permissions, edit the `policy.csv` file:

```csv
# Add a new permission rule
p, manager, /api/v1/users, GET

# Assign a user to a role
g, carol@example.com, manager
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
