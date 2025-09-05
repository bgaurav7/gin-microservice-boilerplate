package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	delivery "github.com/bgaurav7/gin-microservice-boilerplate/internal/delivery/http"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/db"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSwaggerEndpoint(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("ENABLE_SWAGGER", "true")
	os.Setenv("APP_ENVIRONMENT", "dev")
	defer func() {
		os.Unsetenv("ENABLE_SWAGGER")
		os.Unsetenv("APP_ENVIRONMENT")
	}()

	// Create a minimal test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Logger: logger.Config{
			Level: "info",
		},
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret",
			JWTExpiryHours:   1,
			SuperAdminEmail:  "admin@test.com",
		},
	}

	// Initialize logger
	log, err := logger.NewLogger(&cfg.Logger)
	require.NoError(t, err)
	defer log.Sync()

	// Create a mock database (we don't need real DB for this test)
	mockDB := &db.Database{}

	// Create router
	router := delivery.NewRouter(log, mockDB, cfg)

	// Test Swagger JSON endpoint
	t.Run("Swagger JSON endpoint returns 200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/swagger/doc.json", nil)
		w := httptest.NewRecorder()

		router.Handler().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})

	// Test Swagger UI endpoint
	t.Run("Swagger UI endpoint returns 200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/swagger/index.html", nil)
		w := httptest.NewRecorder()

		router.Handler().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
	})
}

func TestSwaggerDisabled(t *testing.T) {
	// Set environment variables to disable Swagger
	os.Setenv("ENABLE_SWAGGER", "false")
	os.Setenv("APP_ENVIRONMENT", "dev")
	defer func() {
		os.Unsetenv("ENABLE_SWAGGER")
		os.Unsetenv("APP_ENVIRONMENT")
	}()

	// Create a minimal test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Logger: logger.Config{
			Level: "info",
		},
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret",
			JWTExpiryHours:   1,
			SuperAdminEmail:  "admin@test.com",
		},
	}

	// Initialize logger
	log, err := logger.NewLogger(&cfg.Logger)
	require.NoError(t, err)
	defer log.Sync()

	// Create a mock database
	mockDB := &db.Database{}

	// Create router
	router := delivery.NewRouter(log, mockDB, cfg)

	// Test that Swagger endpoints return 401 when disabled (auth middleware applied)
	t.Run("Swagger endpoints return 401 when disabled", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/swagger/doc.json", nil)
		w := httptest.NewRecorder()

		router.Handler().ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestSwaggerProductionMode(t *testing.T) {
	// Set environment variables for production mode without auth
	os.Setenv("ENABLE_SWAGGER", "true")
	os.Setenv("APP_ENVIRONMENT", "prod")
	defer func() {
		os.Unsetenv("ENABLE_SWAGGER")
		os.Unsetenv("APP_ENVIRONMENT")
	}()

	// Create a minimal test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Logger: logger.Config{
			Level: "info",
		},
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret",
			JWTExpiryHours:   1,
			SuperAdminEmail:  "admin@test.com",
		},
	}

	// Initialize logger
	log, err := logger.NewLogger(&cfg.Logger)
	require.NoError(t, err)
	defer log.Sync()

	// Create a mock database
	mockDB := &db.Database{}

	// Create router
	router := delivery.NewRouter(log, mockDB, cfg)

	// Test that Swagger endpoints return 401 in production without auth (auth middleware applied)
	t.Run("Swagger endpoints return 401 in production without auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/swagger/doc.json", nil)
		w := httptest.NewRecorder()

		router.Handler().ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestSwaggerProductionModeWithAuth(t *testing.T) {
	// Set environment variables for production mode with auth
	os.Setenv("ENABLE_SWAGGER", "true")
	os.Setenv("APP_ENVIRONMENT", "prod")
	os.Setenv("SWAGGER_BASIC_AUTH_USERNAME", "admin")
	os.Setenv("SWAGGER_BASIC_AUTH_PASSWORD", "password")
	defer func() {
		os.Unsetenv("ENABLE_SWAGGER")
		os.Unsetenv("APP_ENVIRONMENT")
		os.Unsetenv("SWAGGER_BASIC_AUTH_USERNAME")
		os.Unsetenv("SWAGGER_BASIC_AUTH_PASSWORD")
	}()

	// Create a minimal test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Logger: logger.Config{
			Level: "info",
		},
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret",
			JWTExpiryHours:   1,
			SuperAdminEmail:  "admin@test.com",
		},
	}

	// Initialize logger
	log, err := logger.NewLogger(&cfg.Logger)
	require.NoError(t, err)
	defer log.Sync()

	// Create a mock database
	mockDB := &db.Database{}

	// Create router
	router := delivery.NewRouter(log, mockDB, cfg)

	// Test that Swagger endpoints require basic auth in production
	t.Run("Swagger endpoints require basic auth in production", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/swagger/doc.json", nil)
		w := httptest.NewRecorder()

		router.Handler().ServeHTTP(w, req)

		// Should return 401 Unauthorized without basic auth
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test with basic auth
	t.Run("Swagger endpoints work with basic auth in production", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/swagger/doc.json", nil)
		req.SetBasicAuth("admin", "password")
		w := httptest.NewRecorder()

		router.Handler().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})
}