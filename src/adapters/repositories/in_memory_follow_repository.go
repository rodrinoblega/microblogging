package repositories

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
)

type InMemoryFollowRepository struct {
	follows map[uuid.UUID]map[uuid.UUID]bool
}

func NewInMemoryFollowRepository() *InMemoryFollowRepository {
	return &InMemoryFollowRepository{
		follows: make(map[uuid.UUID]map[uuid.UUID]bool),
	}
}

func (r *InMemoryFollowRepository) Save(follow *entities.Follow) error {
	if r.IsFollowing(follow.FollowerID, follow.FollowingID) {
		return errors.New("user already follows this user")
	}

	if _, exists := r.follows[follow.FollowerID]; !exists {
		r.follows[follow.FollowerID] = make(map[uuid.UUID]bool)
	}

	r.follows[follow.FollowerID][follow.FollowingID] = true

	return nil
}

func (r *InMemoryFollowRepository) IsFollowing(followerID, followingID uuid.UUID) bool {
	if followingMap, exists := r.follows[followerID]; exists {
		return followingMap[followingID]
	}

	return false
}

func (r *InMemoryFollowRepository) GetFollowing(userID uuid.UUID) ([]uuid.UUID, error) {
	var followingList []uuid.UUID
	if followingMap, exists := r.follows[userID]; exists {
		for followingID := range followingMap {
			followingList = append(followingList, followingID)
		}
	}

	return followingList, nil
}
