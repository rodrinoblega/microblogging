package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/setup"
	"github.com/rodrinoblega/microblogging/src/adapters/controllers"
	"github.com/rodrinoblega/microblogging/src/adapters/repositories"
	"github.com/rodrinoblega/microblogging/src/usecases"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegration(t *testing.T) {
	router := inititialize()

	userID := uuid.New()

	t.Run("post tweet", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"user_id": userID.String(),
			"content": "First tweet!",
		})

		req, _ := http.NewRequest("POST", "/api/tweet", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

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

}

func inititialize() *gin.Engine {
	gin.SetMode(gin.TestMode)

	tweetRepo := repositories.NewInMemoryTweetRepository()
	followRepo := repositories.NewInMemoryFollowRepository()

	postTweetUseCase := usecases.NewPostTweetUseCase(tweetRepo)
	followUserUseCase := usecases.NewFollowUserUseCase(followRepo)
	getTimelineUseCase := usecases.NewGetTimelineUseCase(followRepo, tweetRepo)

	tweetController := controllers.NewTweetController(postTweetUseCase)
	followController := controllers.NewFollowController(followUserUseCase)
	getTimelineController := controllers.NewTimelineController(getTimelineUseCase)

	appDependencies := &setup.AppDependencies{
		FollowController:   followController,
		TweetController:    tweetController,
		TimelineController: getTimelineController,
	}

	router := SetupRouter(appDependencies)

	return router
}
