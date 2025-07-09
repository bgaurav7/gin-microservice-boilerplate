package http

import (
	"net/http"

	"github.com/bgaurav7/gin-microservice-boilerplate/internal/delivery/http/middleware"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// Router represents the HTTP router
type Router struct {
	engine *gin.Engine
	logger *logger.Logger
}

// NewRouter creates a new HTTP router
func NewRouter(logger *logger.Logger) *Router {
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
}
