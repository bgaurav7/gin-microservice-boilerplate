package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bgaurav7/gin-microservice-boilerplate/internal/domain/model"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTodoRepository is a mock implementation of the TodoRepository interface
type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) GetAll(ctx context.Context) ([]model.Todo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Todo), args.Error(1)
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *model.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func TestTodoUsecase_List(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	log, _ := logger.NewLogger(&logger.Config{Level: "info"})
	todoUsecase := usecase.NewTodoUsecase(mockRepo, log)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		expectedTodos := []model.Todo{
			{ID: 1, Title: "Test Todo 1", Completed: false},
			{ID: 2, Title: "Test Todo 2", Completed: true},
		}

		mockRepo.On("GetAll", ctx).Return(expectedTodos, nil).Once()

		todos, err := todoUsecase.List(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedTodos, todos)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("database error")
		mockRepo.On("GetAll", ctx).Return(nil, expectedError).Once()

		todos, err := todoUsecase.List(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, todos)
		mockRepo.AssertExpectations(t)
	})
}

func TestTodoUsecase_Create(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	log, _ := logger.NewLogger(&logger.Config{Level: "info"})
	todoUsecase := usecase.NewTodoUsecase(mockRepo, log)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		title := "Test Todo"
		mockRepo.On("Create", ctx, mock.MatchedBy(func(todo *model.Todo) bool {
			return todo.Title == title && !todo.Completed
		})).Return(nil).Once()

		todo, err := todoUsecase.Create(ctx, title)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, title, todo.Title)
		assert.False(t, todo.Completed)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		title := "Test Todo"
		expectedError := errors.New("database error")
		mockRepo.On("Create", ctx, mock.MatchedBy(func(todo *model.Todo) bool {
			return todo.Title == title && !todo.Completed
		})).Return(expectedError).Once()

		todo, err := todoUsecase.Create(ctx, title)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, todo)
		mockRepo.AssertExpectations(t)
	})
}
