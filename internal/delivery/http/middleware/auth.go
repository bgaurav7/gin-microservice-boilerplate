package middleware

import (
	"net/http"
	"strings"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/jwt"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware represents the authentication middleware
type AuthMiddleware struct {
	tokenService *jwt.TokenService
	logger       *logger.Logger
	config       *config.AuthConfig
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(tokenService *jwt.TokenService, logger *logger.Logger, config *config.AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
		logger:       logger,
		config:       config,
	}
}

// Authenticate is a middleware that authenticates requests using JWT tokens
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for certain paths
		if c.Request.URL.Path == "/healthz" || c.Request.URL.Path == "/readyz" ||
			c.Request.URL.Path == "/auth" || c.Request.URL.Path == "/public" {
			c.Next()
			return
		}

		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Error("Authorization header is missing", map[string]interface{}{
				"path": c.Request.URL.Path,
			})
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is missing",
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			m.logger.Error("Invalid authorization format", map[string]interface{}{
				"path": c.Request.URL.Path,
			})
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header format must be 'Bearer {token}'",
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := authHeader[len(prefix):]

		// Validate the token
		claims, err := m.tokenService.ValidateToken(tokenString)
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

		// Set user information in context
		c.Set("userEmail", claims.Email)
		c.Set("userID", claims.Subject)

		// Check if user is a super admin
		isSuperAdmin := m.tokenService.IsSuperAdmin(claims.Email)
		c.Set("isSuperAdmin", isSuperAdmin)

		m.logger.Info("User authenticated", map[string]interface{}{
			"email":        claims.Email,
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
