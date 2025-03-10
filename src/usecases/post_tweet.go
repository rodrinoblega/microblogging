package usecases

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
)

type TweetRepository interface {
	Save(tweet *entities.Tweet) error
}

type PostTweetUseCase struct {
	tweetRepository TweetRepository
}

func NewPostTweetUseCase(tweetRepository TweetRepository) *PostTweetUseCase {
	return &PostTweetUseCase{tweetRepository: tweetRepository}
}

func (pt *PostTweetUseCase) Execute(userID uuid.UUID, content string) (*entities.Tweet, error) {
	if len(content) == 0 || len(content) > 280 {
		return nil, errors.New("tweet content must be between 1 and 280 characters")
	}

	tweet := entities.NewTweet(userID, content)
	err := pt.tweetRepository.Save(tweet)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}
