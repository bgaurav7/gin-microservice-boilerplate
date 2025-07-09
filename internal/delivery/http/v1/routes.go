package v1

import (
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/delivery/http/v1/handler"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/db"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/repository"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/usecase"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all API v1 routes
func RegisterRoutes(router *gin.RouterGroup, database *db.Database, logger *logger.Logger) {
	// Initialize repositories
	todoRepo := repository.NewTodoRepository(database, logger)

	// Initialize usecases
	todoUsecase := usecase.NewTodoUsecase(todoRepo, logger)

	// Initialize handlers
	todoHandler := handler.NewTodoHandler(todoUsecase, logger)

	// Register todo routes
	todoRoutes := router.Group("/todos")
	{
		todoRoutes.GET("", todoHandler.GetAll)
		todoRoutes.POST("", todoHandler.Create)
	}
}
