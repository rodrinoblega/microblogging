package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/usecases"
	"net/http"
)

type FollowController struct {
	followUserUseCase *usecases.FollowUserUseCase
}

func NewFollowController(followUserUseCase *usecases.FollowUserUseCase) *FollowController {
	return &FollowController{followUserUseCase: followUserUseCase}
}

func (fc *FollowController) FollowUser(c *gin.Context) {
	var request struct {
		FollowerID  string `json:"follower_id"`
		FollowingID string `json:"following_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	followerID, err := uuid.Parse(request.FollowerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid follower ID"})
		return
	}

	followingID, err := uuid.Parse(request.FollowingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid following ID"})
		return
	}

	err = fc.followUserUseCase.Execute(followerID, followingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User followed successfully"})
}
