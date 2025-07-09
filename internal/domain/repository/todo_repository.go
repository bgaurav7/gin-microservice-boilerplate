package repository

import (
	"context"

	"github.com/bgaurav7/gin-microservice-boilerplate/internal/domain/model"
)

// TodoRepository defines the interface for todo repository operations
type TodoRepository interface {
	// GetAll retrieves all todos from the repository
	GetAll(ctx context.Context) ([]model.Todo, error)

	// Create adds a new todo to the repository
	Create(ctx context.Context, todo *model.Todo) error
}
