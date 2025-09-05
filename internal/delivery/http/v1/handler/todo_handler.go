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
// @Summary List todos
// @Description Get list of todos visible to the user
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} internal_domain_model.Todo
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
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
// @Summary Create todo
// @Description Create a new todo entry
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body TodoCreateRequest true "Todo payload"
// @Success 201 {object} internal_domain_model.Todo
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
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
