package repositories

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

	tweets := r.collectTweets(userIDs)
	r.sortTweetsByDateDesc(tweets)

	if cursor != nil {
		tweets = r.applyCursor(tweets, *cursor)
	}

	return r.applyLimit(tweets, limit), nil
}

func (r *InMemoryTweetRepository) collectTweets(userIDs []uuid.UUID) []*entities.Tweet {
	var tweets []*entities.Tweet
	for _, userID := range userIDs {
		tweets = append(tweets, r.tweets[userID]...)
	}

	return tweets
}

func (r *InMemoryTweetRepository) sortTweetsByDateDesc(tweets []*entities.Tweet) {
	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.After(tweets[j].CreatedAt)
	})
}

func (r *InMemoryTweetRepository) applyCursor(tweets []*entities.Tweet, cursor uuid.UUID) []*entities.Tweet {
	for i, tweet := range tweets {
		if tweet.ID == cursor {
			if i+1 < len(tweets) {
				return tweets[i+1:]
			}
			return []*entities.Tweet{}
		}
	}

	return []*entities.Tweet{}
}

func (r *InMemoryTweetRepository) applyLimit(tweets []*entities.Tweet, limit int) []*entities.Tweet {
	if len(tweets) > limit {
		return tweets[:limit]
	}

	return tweets
}
