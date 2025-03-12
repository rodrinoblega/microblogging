package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type FollowUserUseCaseInterface interface {
	Execute(followerID, followingID uuid.UUID) error
}

type FollowController struct {
	followUserUseCase FollowUserUseCaseInterface
}

func NewFollowController(followUserUseCase FollowUserUseCaseInterface) *FollowController {
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
