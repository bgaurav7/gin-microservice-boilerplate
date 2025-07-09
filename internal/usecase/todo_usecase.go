package usecase

import (
	"context"

	"github.com/bgaurav7/gin-microservice-boilerplate/internal/domain/model"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/domain/repository"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
)

// TodoUsecase defines the interface for todo business logic
type TodoUsecase interface {
	// List returns all todos
	List(ctx context.Context) ([]model.Todo, error)

	// Create creates a new todo with the given title
	Create(ctx context.Context, title string) (*model.Todo, error)
}

// todoUsecase implements the TodoUsecase interface
type todoUsecase struct {
	repo   repository.TodoRepository
	logger *logger.Logger
}

// NewTodoUsecase creates a new todo usecase
func NewTodoUsecase(repo repository.TodoRepository, logger *logger.Logger) TodoUsecase {
	return &todoUsecase{
		repo:   repo,
		logger: logger,
	}
}

// List returns all todos
func (u *todoUsecase) List(ctx context.Context) ([]model.Todo, error) {
	u.logger.Info("Listing all todos", map[string]interface{}{})
	return u.repo.GetAll(ctx)
}

// Create creates a new todo with the given title
func (u *todoUsecase) Create(ctx context.Context, title string) (*model.Todo, error) {
	u.logger.Info("Creating new todo", map[string]interface{}{
		"title": title,
	})
	
	todo := &model.Todo{
		Title:     title,
		Completed: false,
	}

	if err := u.repo.Create(ctx, todo); err != nil {
		u.logger.Error("Failed to create todo", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return todo, nil
}
