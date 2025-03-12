package database

import (
	"fmt"
	"github.com/rodrinoblega/microblogging/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var (
	once             sync.Once
	instancePostgres *gorm.DB
)

func NewPostgres(env *config.Config) *gorm.DB {
	once.Do(func() {
		instancePostgres = postgresDB(env)
	})

	return instancePostgres
}

func postgresDB(env *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		env.PostgresHost,
		env.PgUser,
		env.PgPassword,
		env.PgDatabase,
		env.PostgresPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	fmt.Println("Antes")
	if err != nil {
		panic(err)
	}

	return db
}
