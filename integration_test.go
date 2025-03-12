package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/setup"
	"github.com/rodrinoblega/microblogging/src/entities"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegration(t *testing.T) {
	router := initialize()

	userID := uuid.New()
	var tweet *entities.Tweet

	t.Run("post tweet", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"user_id": userID.String(),
			"content": "First tweet!",
		})

		req, _ := http.NewRequest("POST", "/api/tweet", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		json.Unmarshal(resp.Body.Bytes(), &tweet)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("follow user", func(t *testing.T) {
		followerID := uuid.New()
		followingID := uuid.New()

		requestBody, _ := json.Marshal(map[string]string{
			"follower_id":  followerID.String(),
			"following_id": followingID.String(),
		})

		req, _ := http.NewRequest("POST", "/api/follow", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("get timeline", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/timeline/"+userID.String()+"?limit=10", nil)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var tweets []map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &tweets)

		assert.NotNil(t, tweets)
	})

	t.Run("get timeline with cursor", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/timeline/"+userID.String()+"?limit=10&cursor="+tweet.ID.String(), nil)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var tweets []map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &tweets)

		assert.NotNil(t, tweets)
	})

	t.Run("create user", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"username": "rnoblega",
		})

		req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})

}

func initialize() *gin.Engine {
	gin.SetMode(gin.TestMode)

	appDependencies := setup.InitializeTestDependencies()

	router := SetupRouter(appDependencies)

	return router
}
