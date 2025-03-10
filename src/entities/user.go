package entities

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID
	Username  string
	Followers map[uuid.UUID]bool
}

func NewUser(username string) *User {
	return &User{
		ID:        uuid.New(),
		Username:  username,
		Followers: make(map[uuid.UUID]bool),
	}
}
