package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
	"net/http"
)

type GetTimeLineUseCaseInterface interface {
	Execute(userID uuid.UUID, cursor *uuid.UUID, limit int) ([]*entities.Tweet, error)
}

type TimelineController struct {
	getTimelineUseCase GetTimeLineUseCaseInterface
}

func NewTimelineController(getTimelineUseCase GetTimeLineUseCaseInterface) *TimelineController {
	return &TimelineController{getTimelineUseCase: getTimelineUseCase}
}

func (tc *TimelineController) GetTimeline(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var cursor *uuid.UUID
	if cursorStr := c.Query("cursor"); cursorStr != "" {
		parsedCursor, err := uuid.Parse(cursorStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor"})
			return
		}
		cursor = &parsedCursor
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		var parsedLimit int
		_, err := fmt.Sscanf(limitStr, "%d", &parsedLimit)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	tweets, err := tc.getTimelineUseCase.Execute(userID, cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting timeline"})
		return
	}

	c.JSON(http.StatusOK, tweets)
}
