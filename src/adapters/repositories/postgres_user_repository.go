package repositories

import (
	"github.com/rodrinoblega/microblogging/src/entities"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(user *entities.User) error {
	return r.db.Create(user).Error
}
