package database

import (
	"fmt"
	"github.com/rodrinoblega/microblogging/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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

	log.Printf("Successfully connected to PostgreSQL database")

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
	if err != nil {
		log.Printf("Failed to connect to PostgreSQL database: %v", err)
		panic(err)
	}

	return db
}
