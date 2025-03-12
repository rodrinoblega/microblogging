package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/usecases"
	"net/http"
)

type TweetController struct {
	postTweetUseCase *usecases.PostTweetUseCase
}

func NewTweetController(postTweetUseCase *usecases.PostTweetUseCase) *TweetController {
	return &TweetController{postTweetUseCase: postTweetUseCase}
}

func (tc *TweetController) PostTweet(c *gin.Context) {
	var request struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	tweet, err := tc.postTweetUseCase.Execute(userID, request.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Tweet created successfully. ID %v", tweet.ID)})
}
