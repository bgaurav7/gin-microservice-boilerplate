package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRBACMiddleware(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	logConfig := &logger.Config{Level: "info"}
	log, _ := logger.NewLogger(logConfig)

	// Create a simple model for testing
	m, _ := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act
`)

	// Create a test enforcer
	e, _ := casbin.NewEnforcer(m)
	
	// Add test policies
	e.AddPolicy("admin", "/api/v1/todos", "GET")
	e.AddPolicy("admin", "/api/v1/todos", "POST")
	e.AddPolicy("user", "/api/v1/todos", "GET")
	e.AddGroupingPolicy("alice@example.com", "admin")
	e.AddGroupingPolicy("bob@example.com", "user")

	// Create config with superadmin
	authConfig := &config.AuthConfig{
		SuperAdminEmail: "admin@example.com",
	}

	// Create middleware with mocked enforcer
	middleware := &RBACMiddleware{
		enforcer: e,
		logger:   log,
		config:   authConfig,
	}

	// Test cases
	tests := []struct {
		name       string
		path       string
		method     string
		userEmail  string
		statusCode int
	}{
		{
			name:       "Public endpoint should be accessible",
			path:       "/healthz",
			method:     "GET",
			userEmail:  "",
			statusCode: http.StatusOK,
		},
		{
			name:       "Admin can access GET endpoint",
			path:       "/api/v1/todos",
			method:     "GET",
			userEmail:  "alice@example.com",
			statusCode: http.StatusOK,
		},
		{
			name:       "Admin can access POST endpoint",
			path:       "/api/v1/todos",
			method:     "POST",
			userEmail:  "alice@example.com",
			statusCode: http.StatusOK,
		},
		{
			name:       "User can access GET endpoint",
			path:       "/api/v1/todos",
			method:     "GET",
			userEmail:  "bob@example.com",
			statusCode: http.StatusOK,
		},
		{
			name:       "User cannot access POST endpoint",
			path:       "/api/v1/todos",
			method:     "POST",
			userEmail:  "bob@example.com",
			statusCode: http.StatusForbidden,
		},
		{
			name:       "Unauthorized user cannot access protected endpoint",
			path:       "/api/v1/todos",
			method:     "GET",
			userEmail:  "unknown@example.com",
			statusCode: http.StatusForbidden,
		},
		{
			name:       "Superadmin can access any endpoint",
			path:       "/api/v1/todos",
			method:     "POST",
			userEmail:  "admin@example.com",
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gin context
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			
			// Add test route with middleware
			r.Use(func(c *gin.Context) {
				// Set user email in context if provided
				if tt.userEmail != "" {
					c.Set("userEmail", tt.userEmail)
				}
				c.Next()
			})
			r.Use(middleware.Authorize())
			
			// Add test route
			r.Any(tt.path, func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// Create request
			req, _ := http.NewRequest(tt.method, tt.path, nil)

			// Process request
			r.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
