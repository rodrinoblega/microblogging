package controllers

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockFollowUserUseCase struct {
	mock.Mock
}

func (m *MockFollowUserUseCase) Execute(followerID, followingID uuid.UUID) error {
	args := m.Called(followerID, followingID)
	return args.Error(0)
}

func TestFollowUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockUseCase := new(MockFollowUserUseCase)
	followController := NewFollowController(mockUseCase)

	router.POST("/follow", followController.FollowUser)

	t.Run("invalid request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/follow", bytes.NewBuffer([]byte("{invalid json")))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request")
	})

	t.Run("invalid follower ID", func(t *testing.T) {
		requestBody := []byte(`{"follower_id": "id-inválido", "following_id": "` + uuid.New().String() + `"}`)
		req, err := http.NewRequest(http.MethodPost, "/follow", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid follower ID")
	})

	t.Run("invalid following ID", func(t *testing.T) {
		requestBody := []byte(`{"follower_id": "` + uuid.New().String() + `", "following_id": "id-inválido"}`)
		req, err := http.NewRequest(http.MethodPost, "/follow", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid following ID")
	})

	t.Run("error in execute", func(t *testing.T) {
		followerID := uuid.New()
		followingID := uuid.New()
		mockError := errors.New("error")
		mockUseCase.On("Execute", followerID, followingID).Return(mockError)

		requestBody := []byte(`{"follower_id": "` + followerID.String() + `", "following_id": "` + followingID.String() + `"}`)
		req, err := http.NewRequest(http.MethodPost, "/follow", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), mockError.Error())
	})
}
