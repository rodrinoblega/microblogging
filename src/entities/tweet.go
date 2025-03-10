package entities

import (
	"github.com/google/uuid"
	"time"
)

type Tweet struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
}

func NewTweet(userID uuid.UUID, content string) *Tweet {
	return &Tweet{
		ID:        uuid.New(),
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}
}
