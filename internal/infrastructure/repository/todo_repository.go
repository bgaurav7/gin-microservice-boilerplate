package repository

import (
	"context"

	"github.com/bgaurav7/gin-microservice-boilerplate/internal/domain/model"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/domain/repository"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/db"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
)

// todoRepository implements the TodoRepository interface
type todoRepository struct {
	db     *db.Database
	logger *logger.Logger
}

// NewTodoRepository creates a new todo repository
func NewTodoRepository(db *db.Database, logger *logger.Logger) repository.TodoRepository {
	return &todoRepository{
		db:     db,
		logger: logger,
	}
}

// GetAll retrieves all todos from the repository
func (r *todoRepository) GetAll(ctx context.Context) ([]model.Todo, error) {
	var todos []model.Todo
	result := r.db.DB.Find(&todos)
	if result.Error != nil {
		r.logger.Error("Failed to get todos", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return todos, nil
}

// Create adds a new todo to the repository
func (r *todoRepository) Create(ctx context.Context, todo *model.Todo) error {
	result := r.db.DB.Create(todo)
	if result.Error != nil {
		r.logger.Error("Failed to create todo", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}
