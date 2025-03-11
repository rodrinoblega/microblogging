package adapters

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
)

type InMemoryTweetRepository struct {
	tweets     map[uuid.UUID][]*entities.Tweet
	ShouldFail bool
}

func NewInMemoryTweetRepository() *InMemoryTweetRepository {
	return &InMemoryTweetRepository{
		tweets: make(map[uuid.UUID][]*entities.Tweet),
	}
}

func (r *InMemoryTweetRepository) Save(tweet *entities.Tweet) error {
	if r.ShouldFail {
		return errors.New("simulated error")
	}

	if _, exists := r.tweets[tweet.UserID]; !exists {
		r.tweets[tweet.UserID] = []*entities.Tweet{}
	}

	r.tweets[tweet.UserID] = append(r.tweets[tweet.UserID], tweet)

	return nil
}

func (r *InMemoryTweetRepository) GetTweetsByUsers(userIDs []uuid.UUID) ([]*entities.Tweet, error) {
	if r.ShouldFail {
		return nil, errors.New("simulated error")
	}

	var tweets []*entities.Tweet
	for _, userID := range userIDs {
		if userTweets, exists := r.tweets[userID]; exists {
			tweets = append(tweets, userTweets...)
		}
	}

	return tweets, nil
}
