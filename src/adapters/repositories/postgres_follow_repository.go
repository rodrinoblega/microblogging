package repositories

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
	"gorm.io/gorm"
)

type PostgresFollowRepository struct {
	db *gorm.DB
}

func NewPostgresFollowRepository(db *gorm.DB) *PostgresFollowRepository {
	return &PostgresFollowRepository{db: db}
}

func (r *PostgresFollowRepository) Save(follow *entities.Follow) error {
	if r.IsFollowing(follow.FollowerID, follow.FollowingID) {
		return errors.New("the user is already being followed")
	}
	return r.db.Create(follow).Error
}

func (r *PostgresFollowRepository) GetFollowing(userID uuid.UUID) ([]uuid.UUID, error) {
	var follows []entities.Follow

	err := r.db.Where("follower_id = ?", userID).Find(&follows).Error
	if err != nil {
		return nil, err
	}

	followingIDs := make([]uuid.UUID, len(follows))
	for i, follow := range follows {
		followingIDs[i] = follow.FollowingID
	}

	return followingIDs, nil
}

func (r *PostgresFollowRepository) IsFollowing(followerID, followingID uuid.UUID) bool {
	var count int64
	r.db.Model(&entities.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count)
	return count > 0
}
