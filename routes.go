package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rodrinoblega/microblogging/setup"
)

func SetupRouter(dependencies *setup.AppDependencies) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")

	api.POST("/follow", dependencies.FollowController.FollowUser)

	api.POST("/tweet", dependencies.TweetController.PostTweet)

	api.GET("/timeline/:user_id", dependencies.TimelineController.GetTimeline)

	api.POST("/users", dependencies.UserController.CreateUser)

	return router
}
