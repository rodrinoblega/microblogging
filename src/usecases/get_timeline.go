package usecases

import (
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
)

type GetTimelineUseCase struct {
	followRepository FollowRepository
	tweetRepository  TweetRepository
}

func NewGetTimelineUseCase(followRepository FollowRepository, tweetRepository TweetRepository) *GetTimelineUseCase {
	return &GetTimelineUseCase{
		followRepository: followRepository,
		tweetRepository:  tweetRepository,
	}
}

func (gt *GetTimelineUseCase) Execute(userID uuid.UUID) ([]*entities.Tweet, error) {
	followingUsers, err := gt.followRepository.GetFollowing(userID)

	if len(followingUsers) == 0 {
		return []*entities.Tweet{}, nil
	}

	tweets, err := gt.tweetRepository.GetTweetsByUsers(followingUsers)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}
