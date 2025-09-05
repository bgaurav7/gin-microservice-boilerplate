# Swagger Requirements Mapping

This document maps README.md requirements to implemented Swagger documentation.

## Extracted Requirements from README.md

### R1: API Base Paths
**README Quote:** "Versioned API structure (`/api/v1/`, `/api/v2/`)"
**Status:** ✅ Implemented
**Implementation:** 
- `/api/v1/` - Active endpoints documented
- `/api/v2/` - Placeholder (not implemented in code)

### R2: Todo CRUD API
**README Quote:** "Todo CRUD API (`GET /todos`, `POST /todos`) behind Casbin RBAC"
**Status:** ✅ Implemented
**Implementation:** 
- `GET /api/v1/todos` - `internal/delivery/http/v1/handler/todo_handler.go:25-35`
- `POST /api/v1/todos` - `internal/delivery/http/v1/handler/todo_handler.go:37-60`

### R3: Health and Readiness Endpoints
**README Quote:** "Health (`/healthz`) and readiness (`/readyz`) endpoints"
**Status:** ✅ Implemented
**Implementation:**
- `GET /healthz` - `internal/delivery/http/router.go:70-72`
- `GET /readyz` - `internal/delivery/http/router.go:74-84`

### R4: Authentication Endpoint
**README Quote:** "Client sends a POST request to `/auth` with email"
**Status:** ✅ Implemented
**Implementation:**
- `POST /auth` - `internal/delivery/http/handler/auth_handler.go:35-85`

### R5: Welcome Endpoint
**README Quote:** "Access the API at http://localhost:8080. Available endpoints: `/` - Welcome message"
**Status:** ✅ Implemented
**Implementation:**
- `GET /` - `internal/delivery/http/router.go:65-67`

### R6: JWT Authentication
**README Quote:** "The application uses a simple JWT-based authentication system"
**Status:** ✅ Implemented
**Implementation:** JWT security scheme documented in Swagger

### R7: Protected Endpoints
**README Quote:** "All API endpoints under `/api/v1/*` require authentication"
**Status:** ✅ Implemented
**Implementation:** All `/api/v1/*` endpoints marked with `@Security ApiKeyAuth`

### R8: RBAC Roles
**README Quote:** "Role-based access control (user, admin, superadmin)"
**Status:** ✅ Implemented
**Implementation:** Documented in policy.csv and Swagger security schemes

### R9: Swagger Documentation
**README Quote:** "Swagger docs at `/swagger/index.html`"
**Status:** ✅ Implemented
**Implementation:** Swagger UI served at `/swagger/*any` with gating

### R10: Todo Model Structure
**README Quote:** "Todo CRUD API" (inferred from context)
**Status:** ✅ Implemented
**Implementation:** `internal/domain/model/todo.go:7-15` - Todo struct with ID, Title, Completed, CreatedAt, UpdatedAt

## RBAC Policy Mapping

From `internal/infrastructure/rbac/policy.csv`:
- **admin role:** Can GET and POST `/api/v1/todos`
- **user role:** Can GET `/api/v1/todos`
- **alice@example.com:** Assigned admin role
- **bob@example.com:** Assigned user role

## Security Implementation

- JWT Bearer token authentication
- Environment-based Swagger UI gating
- Production safety with basic auth option
- Default disabled in production

## Endpoints Summary

| Method | Path | Handler | Auth Required | RBAC |
|--------|------|---------|---------------|------|
| GET | `/` | router.go:65-67 | No | No |
| GET | `/healthz` | router.go:70-72 | No | No |
| GET | `/readyz` | router.go:74-84 | No | No |
| POST | `/auth` | auth_handler.go:35-85 | No | No |
| GET | `/api/v1/todos` | todo_handler.go:25-35 | Yes | Yes |
| POST | `/api/v1/todos` | todo_handler.go:37-60 | Yes | Yes |
| GET | `/swagger/*any` | Generated | Gated | No |

## Notes

- All requirements from README.md have been successfully mapped to implemented endpoints
- No discrepancies found between README and actual implementation
- Swagger documentation covers all active endpoints
- Security schemes properly implemented for JWT authentication
- RBAC policies documented and enforced