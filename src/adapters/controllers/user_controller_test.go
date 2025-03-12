package controllers

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rodrinoblega/microblogging/src/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockCreateUserUseCase struct {
	mock.Mock
}

func (m *MockCreateUserUseCase) Execute(username string) (*entities.User, error) {
	args := m.Called(username)
	return args.Get(0).(*entities.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockUseCase := new(MockCreateUserUseCase)
	userController := NewUserController(mockUseCase)

	router.POST("/users", userController.CreateUser)

	t.Run("invalid request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("{invalid json")))
		assert.NoError(t, err)

		// Record the response
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request")
	})

	t.Run("error in execute", func(t *testing.T) {
		username := "testuser"
		mockError := errors.New("use case error")
		mockUseCase.On("Execute", username).Return(&entities.User{}, mockError)

		requestBody := []byte(`{"username": "` + username + `"}`)
		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), mockError.Error())
	})
}
