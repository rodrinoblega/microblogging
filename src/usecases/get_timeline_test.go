package usecases

import (
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/adapters"
	"github.com/rodrinoblega/microblogging/src/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTimelineUseCase(t *testing.T) {
	followRepository := adapters.NewInMemoryFollowRepository()
	tweetRepository := adapters.NewInMemoryTweetRepository()
	timelineUseCase := NewGetTimelineUseCase(followRepository, tweetRepository)

	user1 := uuid.New()
	user2 := uuid.New()

	t.Run("should get valid timelines", func(t *testing.T) {

		err := followRepository.Save(&entities.Follow{FollowerID: user1, FollowingID: user2})
		assert.NoError(t, err)

		tweet1 := &entities.Tweet{ID: uuid.New(), UserID: user2, Content: "This is a test tweet 1"}
		err = tweetRepository.Save(tweet1)
		assert.NoError(t, err)

		tweet2 := &entities.Tweet{ID: uuid.New(), UserID: user2, Content: "This is a test tweet 2"}
		err = tweetRepository.Save(tweet2)
		assert.NoError(t, err)

		timelineUser1, err := timelineUseCase.Execute(user1)
		assert.NoError(t, err)
		assert.Len(t, timelineUser1, 2)
		assert.Equal(t, "This is a test tweet 1", timelineUser1[0].Content)
		assert.Equal(t, "This is a test tweet 2", timelineUser1[1].Content)

		timelineUser2, err := timelineUseCase.Execute(user2)
		assert.NoError(t, err)
		assert.Len(t, timelineUser2, 0)
	})

	t.Run("should return error on gettimeline", func(t *testing.T) {
		tweetRepository.ShouldFail = true

		err := followRepository.Save(&entities.Follow{FollowerID: user1, FollowingID: user2})
		assert.NoError(t, err)

		timelineUser1, err := timelineUseCase.Execute(user1)

		assert.Error(t, err)
		assert.Nil(t, timelineUser1)
		assert.EqualError(t, err, "simulated error")
	})

}
