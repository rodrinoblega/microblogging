package usecases

import (
	"errors"
	"github.com/rodrinoblega/microblogging/src/entities"
)

type UserRepository interface {
	Save(user *entities.User) error
}

type CreateUserUseCase struct {
	userRepository UserRepository
}

func NewCreateUserUseCase(userRepository UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: userRepository}
}

func (cu *CreateUserUseCase) Execute(username string) (*entities.User, error) {
	if len(username) == 0 {
		return nil, errors.New("username must not be empty")
	}

	user := entities.NewUser(username)
	err := cu.userRepository.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
