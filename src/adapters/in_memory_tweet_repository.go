package adapters

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
	"sort"
	"time"
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

	if tweet.CreatedAt.IsZero() {
		tweet.CreatedAt = time.Now()
	}

	if _, exists := r.tweets[tweet.UserID]; !exists {
		r.tweets[tweet.UserID] = []*entities.Tweet{}
	}

	r.tweets[tweet.UserID] = append(r.tweets[tweet.UserID], tweet)

	return nil
}

func (r *InMemoryTweetRepository) GetTweetsByUsers(userIDs []uuid.UUID, cursor *uuid.UUID, limit int) ([]*entities.Tweet, error) {
	if r.ShouldFail {
		return nil, errors.New("simulated error")
	}

	var tweets []*entities.Tweet
	for _, userID := range userIDs {
		if userTweets, exists := r.tweets[userID]; exists {
			tweets = append(tweets, userTweets...)
		}
	}

	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.After(tweets[j].CreatedAt)
	})

	if cursor != nil {
		index := -1
		for i, tweet := range tweets {
			if tweet.ID == *cursor {
				index = i + 1
				break
			}
		}
		if index == -1 || index >= len(tweets) {
			return []*entities.Tweet{}, nil
		}
		tweets = tweets[index:]
	}

	// Aplicar limit
	if len(tweets) > limit {
		tweets = tweets[:limit]
	}

	return tweets, nil
}
