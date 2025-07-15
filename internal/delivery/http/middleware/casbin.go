package middleware

import (
	"net/http"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/rbac"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// RBACMiddleware represents the RBAC middleware
type RBACMiddleware struct {
	enforcer *casbin.Enforcer
	logger   *logger.Logger
	config   *config.AuthConfig
}

// NewRBACMiddleware creates a new RBAC middleware
func NewRBACMiddleware(logger *logger.Logger, config *config.AuthConfig) (*RBACMiddleware, error) {
	enforcer, err := rbac.NewEnforcer()
	if err != nil {
		return nil, err
	}

	return &RBACMiddleware{
		enforcer: enforcer,
		logger:   logger,
		config:   config,
	}, nil
}

// Authorize is a middleware that authorizes requests using Casbin
func (m *RBACMiddleware) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authorization for certain paths
		if c.Request.URL.Path == "/healthz" || c.Request.URL.Path == "/readyz" ||
			c.Request.URL.Path == "/auth" || c.Request.URL.Path == "/public" {
			c.Next()
			return
		}

		// Get user email from context (set by auth middleware)
		userEmail, exists := c.Get("userEmail")
		if !exists {
			m.logger.Error("User email not found in context", nil)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		// Check if user is superadmin
		email, ok := userEmail.(string)
		if !ok {
			m.logger.Error("User email is not a string", nil)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			c.Abort()
			return
		}

		// Superadmin override - always allow access
		if email == m.config.SuperAdminEmail {
			m.logger.Info("Superadmin access granted", map[string]interface{}{
				"email": email,
				"path":  c.Request.URL.Path,
				"method": c.Request.Method,
			})
			c.Next()
			return
		}

		// Check if user has permission
		obj := c.Request.URL.Path
		act := c.Request.Method
		allowed, err := m.enforcer.Enforce(email, obj, act)
		if err != nil {
			m.logger.Error("Casbin enforcement error", map[string]interface{}{
				"error": err.Error(),
				"email": email,
				"path":  obj,
				"method": act,
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			c.Abort()
			return
		}

		if !allowed {
			m.logger.Warn("Access denied", map[string]interface{}{
				"email": email,
				"path":  obj,
				"method": act,
			})
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			c.Abort()
			return
		}

		m.logger.Info("Access granted", map[string]interface{}{
			"email": email,
			"path":  obj,
			"method": act,
		})
		c.Next()
	}
}
