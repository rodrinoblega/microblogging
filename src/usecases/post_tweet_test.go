package usecases

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/adapters/repositories"
	"github.com/rodrinoblega/microblogging/src/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPostTweetUseCase(t *testing.T) {
	tweetRepository := repositories.NewInMemoryTweetRepository()
	postTweetUseCase := NewPostTweetUseCase(tweetRepository)

	userID := uuid.New()

	t.Run("should post a valid tweet", func(t *testing.T) {
		tweet, err := postTweetUseCase.Execute(userID, "This is a valid tweet!")

		require.NoError(t, err)
		require.NotNil(t, tweet)

		assert.Equal(t, userID, tweet.UserID)
		assert.Equal(t, "This is a valid tweet!", tweet.Content)
	})

	t.Run("should return error for empty tweet", func(t *testing.T) {
		tweet, err := postTweetUseCase.Execute(userID, "")

		assert.Error(t, err)
		assert.Nil(t, tweet)
		assert.EqualError(t, err, "tweet content must be between 1 and 280 characters")

	})

	t.Run("should return error for tweet exceeding character limit", func(t *testing.T) {
		longTweet := string(make([]byte, 281))
		tweet, err := postTweetUseCase.Execute(userID, longTweet)

		assert.Error(t, err)
		assert.Nil(t, tweet)
		assert.EqualError(t, err, "tweet content must be between 1 and 280 characters")
	})

	t.Run("Test get tweets by users with pagination ", func(t *testing.T) {
		user1 := uuid.New()
		user2 := uuid.New()

		tweets := []*entities.Tweet{
			{ID: uuid.New(), UserID: user1, Content: "Tweet 1", CreatedAt: time.Now().Add(-5 * time.Minute)},
			{ID: uuid.New(), UserID: user1, Content: "Tweet 2", CreatedAt: time.Now().Add(-3 * time.Minute)},
			{ID: uuid.New(), UserID: user2, Content: "Tweet 3", CreatedAt: time.Now().Add(-2 * time.Minute)},
			{ID: uuid.New(), UserID: user2, Content: "Tweet 4", CreatedAt: time.Now().Add(-1 * time.Minute)},
		}

		for _, tweet := range tweets {
			err := tweetRepository.Save(tweet)
			assert.NoError(t, err)
		}

		// Tweets of user 1 and limit 1
		result, err := tweetRepository.GetTweetsByUsers([]uuid.UUID{user1}, nil, 1)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Tweet 2", result[0].Content)

		// Next page using last ID as cursor
		cursor := result[0].ID
		result, err = tweetRepository.GetTweetsByUsers([]uuid.UUID{user1}, &cursor, 1)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Tweet 1", result[0].Content)

		// Try next page, no more tweets
		cursor = result[0].ID
		result, err = tweetRepository.GetTweetsByUsers([]uuid.UUID{user1}, &cursor, 1)
		assert.NoError(t, err)
		assert.Len(t, result, 0)

		result, err = tweetRepository.GetTweetsByUsers([]uuid.UUID{user1, user2}, nil, 3)
		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "Tweet 4", result[0].Content)
		assert.Equal(t, "Tweet 3", result[1].Content)
		assert.Equal(t, "Tweet 2", result[2].Content)
	})

	t.Run("should return error if saving tweet fails", func(t *testing.T) {
		tweetRepository.ShouldFail = true

		tweet, err := postTweetUseCase.Execute(userID, "This tweet will fail to save")

		assert.Error(t, err)
		assert.Nil(t, tweet)
		assert.EqualError(t, err, "simulated error")
	})

	t.Run("", func(t *testing.T) {
		var maptest map[string]string

		if _, exists := maptest["asd"]; exists {
			fmt.Println(maptest["asd"])
		}
	})
}
