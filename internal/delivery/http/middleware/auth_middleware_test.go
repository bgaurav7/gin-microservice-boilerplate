package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/jwt"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Setup test config
	authConfig := &config.AuthConfig{
		JWTSecret:      "test-secret",
		JWTExpiryHours: 1,
		SuperAdminEmail: "admin@example.com",
	}

	// Setup logger and token service
	log, _ := logger.NewLogger(&logger.Config{Level: "info"})
	tokenService := jwt.NewTokenService(authConfig)

	// Setup auth middleware
	authMiddleware := NewAuthMiddleware(tokenService, log, authConfig)

	// Setup gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add middleware and test route
	router.Use(authMiddleware.Authenticate())
	
	// Add protected route that requires authentication
	router.GET("/protected", authMiddleware.RequireAuthentication(), func(c *gin.Context) {
		userEmail, _ := c.Get("userEmail")
		isSuperAdmin, _ := c.Get("isSuperAdmin")
		c.JSON(http.StatusOK, gin.H{
			"userEmail": userEmail,
			"isSuperAdmin": isSuperAdmin,
		})
	})
	
	// Add public route
	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "public"})
	})

	// Add health check route (should be skipped by auth middleware)
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Test cases
	t.Run("Public endpoints should be accessible without token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/public", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Health check endpoint should be accessible without token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/healthz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Protected endpoint should require token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Protected endpoint should be accessible with valid token", func(t *testing.T) {
		// Generate token for test user
		token, _ := tokenService.GenerateToken("user@example.com")
		
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Superadmin should have isSuperAdmin flag set to true", func(t *testing.T) {
		// Generate token for superadmin user
		token, _ := tokenService.GenerateToken("admin@example.com")
		
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "admin@example.com", response["userEmail"])
		assert.Equal(t, true, response["isSuperAdmin"])
	})

	t.Run("Regular user should have isSuperAdmin flag set to false", func(t *testing.T) {
		// Generate token for regular user
		token, _ := tokenService.GenerateToken("user@example.com")
		
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user@example.com", response["userEmail"])
		assert.Equal(t, false, response["isSuperAdmin"])
	})

	t.Run("Invalid token should be rejected", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
