package http

import (
	"net/http"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/delivery/http/handler"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/delivery/http/middleware"
	v1 "github.com/bgaurav7/gin-microservice-boilerplate/internal/delivery/http/v1"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/db"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/jwt"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// Router represents the HTTP router
type Router struct {
	engine         *gin.Engine
	logger         *logger.Logger
	db             *db.Database
	tokenService   *jwt.TokenService
	config         *config.Config
	authMiddleware *middleware.AuthMiddleware
	rbacMiddleware *middleware.RBACMiddleware
}

// NewRouter creates a new HTTP router
func NewRouter(logger *logger.Logger, database *db.Database, cfg *config.Config) *Router {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	engine := gin.New()

	// Add middleware
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logger(logger))

	// Create JWT token service
	tokenService := jwt.NewTokenService(&cfg.Auth)

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService, logger, &cfg.Auth)

	// Create RBAC middleware
	rbacMiddleware, err := middleware.NewRBACMiddleware(logger, &cfg.Auth)
	if err != nil {
		logger.Error("Failed to create RBAC middleware", map[string]interface{}{"error": err.Error()})
		// Continue without RBAC if it fails to initialize
	}

	// Create router
	router := &Router{
		engine:         engine,
		logger:         logger,
		db:             database,
		tokenService:   tokenService,
		config:         cfg,
		authMiddleware: authMiddleware,
		rbacMiddleware: rbacMiddleware,
	}

	// Apply auth middleware globally for JWT parsing
	engine.Use(authMiddleware.Authenticate())

	// Register routes
	router.registerRoutes()

	return router
}

// Handler returns the HTTP handler
func (r *Router) Handler() http.Handler {
	return r.engine
}

// registerRoutes registers all routes
func (r *Router) registerRoutes() {
	// Root path
	r.engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Gin Microservice Boilerplate")
	})

	// Health check
	r.engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Readiness check
	r.engine.GET("/readyz", func(c *gin.Context) {
		// Check database connection
		if err := r.db.Ping(); err != nil {
			r.logger.Error("Database connection failed", map[string]interface{}{"error": err.Error()})
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "reason": "database connection failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// Auth routes
	authHandler := handler.NewAuthHandler(r.tokenService, r.logger, &r.config.Auth)
	r.engine.POST("/auth", authHandler.Authenticate)

	// API v1 routes - protected by auth middleware and RBAC
	apiV1 := r.engine.Group("/api/v1")
	apiV1.Use(r.authMiddleware.RequireAuthentication())

	// Apply RBAC middleware if available
	if r.rbacMiddleware != nil {
		apiV1.Use(r.rbacMiddleware.Authorize())
		r.logger.Info("RBAC middleware applied to /api/v1 routes", nil)
	} else {
		r.logger.Warn("RBAC middleware not available, skipping RBAC enforcement", nil)
	}

	v1.RegisterRoutes(apiV1, r.db, r.logger)
}
