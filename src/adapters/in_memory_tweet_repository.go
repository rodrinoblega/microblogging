package adapters

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
)

type InMemoryTweetRepository struct {
	tweets     map[uuid.UUID]*entities.Tweet
	ShouldFail bool
}

func NewInMemoryTweetRepository() *InMemoryTweetRepository {
	return &InMemoryTweetRepository{
		tweets: make(map[uuid.UUID]*entities.Tweet),
	}
}

func (r *InMemoryTweetRepository) Save(tweet *entities.Tweet) error {
	if r.ShouldFail {
		return errors.New("simulated error")
	}

	r.tweets[tweet.ID] = tweet
	return nil
}
