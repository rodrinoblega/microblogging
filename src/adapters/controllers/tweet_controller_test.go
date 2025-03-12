package controllers

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockPostTweetUseCase struct {
	mock.Mock
}

func (m *MockPostTweetUseCase) Execute(userID uuid.UUID, content string) (*entities.Tweet, error) {
	args := m.Called(userID, content)
	return args.Get(0).(*entities.Tweet), args.Error(1)
}

func TestPostTweet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockUseCase := new(MockPostTweetUseCase)
	tweetController := NewTweetController(mockUseCase)

	router.POST("/tweets", tweetController.PostTweet)

	t.Run("invalid request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/tweets", bytes.NewBuffer([]byte("{invalid request")))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request")
	})

	t.Run("user ID invalid", func(t *testing.T) {
		requestBody := []byte(`{"user_id": "invalid-id", "content": "Hello"}`)
		req, err := http.NewRequest(http.MethodPost, "/tweets", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid user ID")
	})

	t.Run("error in execute", func(t *testing.T) {
		userID := uuid.New()
		content := "Hello"
		mockError := errors.New("error")
		mockUseCase.On("Execute", userID, content).Return(&entities.Tweet{}, mockError)

		requestBody := []byte(`{"user_id": "` + userID.String() + `", "content": "` + content + `"}`)
		req, err := http.NewRequest(http.MethodPost, "/tweets", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), mockError.Error())
	})
}
