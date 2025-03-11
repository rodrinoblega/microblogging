package usecases

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
)

type FollowRepository interface {
	Save(follow *entities.Follow) error
	IsFollowing(followerID, followingID uuid.UUID) bool
	GetFollowing(userID uuid.UUID) ([]uuid.UUID, error)
}

type FollowUserUseCase struct {
	followRepository FollowRepository
}

func NewFollowUserUseCase(followRepository FollowRepository) *FollowUserUseCase {
	return &FollowUserUseCase{followRepository: followRepository}
}

func (fu *FollowUserUseCase) Execute(followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return errors.New("user cannot follow themselves")
	}

	follow := &entities.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return fu.followRepository.Save(follow)
}
