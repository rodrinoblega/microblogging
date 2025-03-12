package usecases

import (
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/adapters/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFollowUserUseCase(t *testing.T) {
	followRepository := repositories.NewInMemoryFollowRepository()
	followUserUseCase := NewFollowUserUseCase(followRepository)

	user1 := uuid.New()
	user2 := uuid.New()
	user3 := uuid.New()

	t.Run("should follow a user successfully", func(t *testing.T) {
		err := followUserUseCase.Execute(user1, user2)

		require.NoError(t, err)

		assert.True(t, followRepository.IsFollowing(user1, user2))
	})

	t.Run("should return error if already following", func(t *testing.T) {
		err := followUserUseCase.Execute(user2, user3)
		require.NoError(t, err)

		err = followUserUseCase.Execute(user2, user3)

		assert.Error(t, err)
		assert.EqualError(t, err, "user already follows this user")
	})

	t.Run("should not allow a user to follow themselves", func(t *testing.T) {
		err := followUserUseCase.Execute(user1, user1)

		assert.Error(t, err)
		assert.EqualError(t, err, "user cannot follow themselves")
	})
}
