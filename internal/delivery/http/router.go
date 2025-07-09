package http

import (
	"net/http"

	"github.com/bgaurav7/gin-microservice-boilerplate/internal/delivery/http/middleware"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/db"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// Router represents the HTTP router
type Router struct {
	engine *gin.Engine
	logger *logger.Logger
	db     *db.Database
}

// NewRouter creates a new HTTP router
func NewRouter(logger *logger.Logger, database *db.Database) *Router {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	engine := gin.New()

	// Add middleware
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logger(logger))

	// Create router
	router := &Router{
		engine: engine,
		logger: logger,
		db:     database,
	}

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
}
