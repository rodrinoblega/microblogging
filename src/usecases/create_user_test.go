package usecases

import (
	"github.com/rodrinoblega/microblogging/src/adapters/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateUserUseCase(t *testing.T) {
	userRepository := repositories.NewInMemoryUserRepository()
	createUserUseCase := NewCreateUserUseCase(userRepository)

	t.Run("should create a valid user", func(t *testing.T) {
		user, err := createUserUseCase.Execute("rnoblega")

		require.NoError(t, err)
		require.NotNil(t, user)

		assert.Equal(t, "rnoblega", user.Username)
	})

	t.Run("should return error for empty username", func(t *testing.T) {
		user, err := createUserUseCase.Execute("")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, "username must not be empty")
	})
}
