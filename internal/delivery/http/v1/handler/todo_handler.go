package handler

import (
	"net/http"

	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/usecase"
	"github.com/gin-gonic/gin"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	todoUsecase usecase.TodoUsecase
	logger      *logger.Logger
}

// TodoCreateRequest represents the request body for creating a todo
type TodoCreateRequest struct {
	Title string `json:"title" binding:"required"`
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(todoUsecase usecase.TodoUsecase, logger *logger.Logger) *TodoHandler {
	return &TodoHandler{
		todoUsecase: todoUsecase,
		logger:      logger,
	}
}

// GetAll godoc
// @Summary Get all todos
// @Description Get all todos
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} model.Todo
// @Router /api/v1/todos [get]
func (h *TodoHandler) GetAll(c *gin.Context) {
	todos, err := h.todoUsecase.List(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get todos", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Create godoc
// @Summary Create a new todo
// @Description Create a new todo
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body TodoCreateRequest true "Todo object"
// @Success 201 {object} model.Todo
// @Failure 400 {object} map[string]string "Invalid request"
// @Router /api/v1/todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	var req TodoCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	todo, err := h.todoUsecase.Create(c.Request.Context(), req.Title)
	if err != nil {
		h.logger.Error("Failed to create todo", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}
