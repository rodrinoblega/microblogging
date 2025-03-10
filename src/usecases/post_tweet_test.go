package usecases

import (
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/adapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPostTweetUseCase(t *testing.T) {
	tweetRepository := adapters.NewInMemoryTweetRepository()
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

	t.Run("should return error if saving tweet fails", func(t *testing.T) {
		tweetRepository.ShouldFail = true

		tweet, err := postTweetUseCase.Execute(userID, "This tweet will fail to save")

		assert.Error(t, err)
		assert.Nil(t, tweet)
		assert.EqualError(t, err, "simulated error")
	})
}
