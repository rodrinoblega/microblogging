package repositories

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
)

type InMemoryUserRepository struct {
	users      map[uuid.UUID][]*entities.User
	ShouldFail bool
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[uuid.UUID][]*entities.User),
	}
}

func (r *InMemoryUserRepository) Save(user *entities.User) error {
	if r.ShouldFail {
		return errors.New("simulated error")
	}

	if _, exists := r.users[user.ID]; !exists {
		r.users[user.ID] = []*entities.User{}
	}

	r.users[user.ID] = append(r.users[user.ID], user)

	return nil
}
