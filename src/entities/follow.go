package entities

import (
	"github.com/google/uuid"
	"time"
)

type Follow struct {
	FollowerID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	FollowingID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
