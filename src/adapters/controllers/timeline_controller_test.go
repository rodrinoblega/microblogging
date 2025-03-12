package controllers

import (
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

type MockGetTimelineUseCase struct {
	mock.Mock
}

func (m *MockGetTimelineUseCase) Execute(userID uuid.UUID, cursor *uuid.UUID, limit int) ([]*entities.Tweet, error) {
	args := m.Called(userID, cursor, limit)
	return args.Get(0).([]*entities.Tweet), args.Error(1)
}

func TestGetTimeline(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockUseCase := new(MockGetTimelineUseCase)
	timelineController := NewTimelineController(mockUseCase)

	router.GET("/timeline/:user_id", timelineController.GetTimeline)

	t.Run("invalid user ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/timeline/invalid-uuid", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid user ID")
	})

	t.Run("invalid cursor", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/timeline/"+uuid.New().String()+"?cursor=invalid-cursor", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid cursor")
	})

	t.Run("use case error", func(t *testing.T) {
		userID := uuid.New()
		mockError := errors.New("error")
		mockUseCase.On("Execute", userID, (*uuid.UUID)(nil), 10).Return([]*entities.Tweet{}, mockError)

		req, err := http.NewRequest(http.MethodGet, "/timeline/"+userID.String(), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Error getting timeline")
	})
}
