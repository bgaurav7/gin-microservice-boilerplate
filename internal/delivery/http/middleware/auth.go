package middleware

import (
	"net/http"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/dex"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware represents the authentication middleware
type AuthMiddleware struct {
	dexClient *dex.Client
	logger    *logger.Logger
	config    *config.AuthConfig
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(dexClient *dex.Client, logger *logger.Logger, config *config.AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{
		dexClient: dexClient,
		logger:    logger,
		config:    config,
	}
}

// Authenticate is a middleware that authenticates requests using JWT tokens
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for certain paths
		if c.Request.URL.Path == "/healthz" || c.Request.URL.Path == "/readyz" ||
			c.Request.URL.Path == "/auth/login" || c.Request.URL.Path == "/auth/callback" {
			c.Next()
			return
		}

		// Verify the token
		userInfo, err := m.dexClient.VerifyToken(c.Request.Context(), c.Request)
		if err != nil {
			m.logger.Error("Authentication failed", map[string]interface{}{
				"error": err.Error(),
				"path":  c.Request.URL.Path,
			})
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication failed",
			})
			c.Abort()
			return
		}

		// Check if token is still valid
		if !m.dexClient.TokenValid(userInfo) {
			m.logger.Error("Token expired", map[string]interface{}{
				"path": c.Request.URL.Path,
			})
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token expired",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("userEmail", userInfo.Email)
		c.Set("userID", userInfo.Subject)
		c.Set("userName", userInfo.Name)
		
		// Check if user is a super admin
		isSuperAdmin := m.dexClient.IsSuperAdmin(userInfo.Email, m.config.SuperAdminEmail)
		c.Set("isSuperAdmin", isSuperAdmin)

		m.logger.Info("User authenticated", map[string]interface{}{
			"email":        userInfo.Email,
			"path":         c.Request.URL.Path,
			"isSuperAdmin": isSuperAdmin,
		})

		c.Next()
	}
}

// RequireAuthentication is a middleware that requires authentication
func (m *AuthMiddleware) RequireAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		if _, exists := c.Get("userEmail"); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireSuperAdmin is a middleware that requires super admin privileges
func (m *AuthMiddleware) RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is a super admin
		isSuperAdmin, exists := c.Get("isSuperAdmin")
		if !exists || !isSuperAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Super admin privileges required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
