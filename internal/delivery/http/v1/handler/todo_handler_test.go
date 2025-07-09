package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgaurav7/gin-microservice-boilerplate/internal/delivery/http/v1/handler"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/domain/model"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTodoUsecase is a mock implementation of the TodoUsecase interface
type MockTodoUsecase struct {
	mock.Mock
}

func (m *MockTodoUsecase) List(ctx context.Context) ([]model.Todo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Todo), args.Error(1)
}

func (m *MockTodoUsecase) Create(ctx context.Context, title string) (*model.Todo, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Todo), args.Error(1)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestTodoHandler_GetAll(t *testing.T) {
	mockUsecase := new(MockTodoUsecase)
	log, _ := logger.NewLogger(&logger.Config{Level: "info"})
	todoHandler := handler.NewTodoHandler(mockUsecase, log)
	router := setupRouter()
	router.GET("/api/v1/todos", todoHandler.GetAll)

	t.Run("Success", func(t *testing.T) {
		expectedTodos := []model.Todo{
			{ID: 1, Title: "Test Todo 1", Completed: false},
			{ID: 2, Title: "Test Todo 2", Completed: true},
		}

		mockUsecase.On("List", mock.Anything).Return(expectedTodos, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/todos", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []model.Todo
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedTodos, response)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Empty List", func(t *testing.T) {
		mockUsecase.On("List", mock.Anything).Return([]model.Todo{}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/todos", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []model.Todo
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Empty(t, response)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUsecase.On("List", mock.Anything).Return(nil, errors.New("database error")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/todos", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestTodoHandler_Create(t *testing.T) {
	mockUsecase := new(MockTodoUsecase)
	log, _ := logger.NewLogger(&logger.Config{Level: "info"})
	todoHandler := handler.NewTodoHandler(mockUsecase, log)
	router := setupRouter()
	router.POST("/api/v1/todos", todoHandler.Create)

	t.Run("Success", func(t *testing.T) {
		title := "Test Todo"
		expectedTodo := &model.Todo{
			ID:        1,
			Title:     title,
			Completed: false,
		}

		mockUsecase.On("Create", mock.Anything, title).Return(expectedTodo, nil).Once()

		reqBody, _ := json.Marshal(handler.TodoCreateRequest{Title: title})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response model.Todo
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedTodo.ID, response.ID)
		assert.Equal(t, expectedTodo.Title, response.Title)
		assert.Equal(t, expectedTodo.Completed, response.Completed)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Missing Title", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Empty Title", func(t *testing.T) {
		reqBody, _ := json.Marshal(handler.TodoCreateRequest{Title: ""})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error", func(t *testing.T) {
		title := "Test Todo"
		mockUsecase.On("Create", mock.Anything, title).Return(nil, errors.New("database error")).Once()

		reqBody, _ := json.Marshal(handler.TodoCreateRequest{Title: title})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}
