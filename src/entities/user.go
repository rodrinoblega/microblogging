package entities

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username  string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func NewUser(username string) *User {
	return &User{
		ID:        uuid.New(),
		Username:  username,
		CreatedAt: time.Now(),
	}
}
