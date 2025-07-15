package handler

import (
	"bytes"
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

func TestAuthHandler_Authenticate(t *testing.T) {
	// Setup test config
	authConfig := &config.AuthConfig{
		JWTSecret:       "test-secret",
		JWTExpiryHours:  1,
		SuperAdminEmail: "admin@example.com",
	}

	// Setup logger and token service
	log, _ := logger.NewLogger(&logger.Config{Level: "info"})
	tokenService := jwt.NewTokenService(authConfig)

	// Setup auth handler
	authHandler := NewAuthHandler(tokenService, log, authConfig)

	// Setup gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth", authHandler.Authenticate)

	// Test cases
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		checkToken     bool
	}{
		{
			name:           "Valid email",
			requestBody:    map[string]interface{}{"email": "user@example.com"},
			expectedStatus: http.StatusOK,
			checkToken:     true,
		},
		{
			name:           "Empty email",
			requestBody:    map[string]interface{}{"email": ""},
			expectedStatus: http.StatusBadRequest,
			checkToken:     false,
		},
		{
			name:           "Missing email field",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			checkToken:     false,
		},
		{
			name:           "Invalid email format",
			requestBody:    map[string]interface{}{"email": "not-an-email"},
			expectedStatus: http.StatusBadRequest,
			checkToken:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check token if expected
			if tt.checkToken {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check if token exists and is not empty
				token, exists := response["token"]
				assert.True(t, exists)
				assert.NotEmpty(t, token)

				// Verify token is valid
				claims, err := tokenService.ValidateToken(token)
				assert.NoError(t, err)
				assert.Equal(t, tt.requestBody["email"], claims.Email)
			}
		})
	}
}
