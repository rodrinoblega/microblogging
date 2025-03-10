package adapters

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
	"sync"
)

type InMemoryTweetRepository struct {
	tweets     map[uuid.UUID]*entities.Tweet
	mu         sync.RWMutex
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

	r.mu.Lock()
	defer r.mu.Unlock()
	r.tweets[tweet.ID] = tweet
	return nil
}
