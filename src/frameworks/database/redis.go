package database

import (
	"github.com/go-redis/redis/v8"
	"github.com/rodrinoblega/microblogging/config"
	"log"
	"sync"
)

var (
	onceRedis     sync.Once
	instanceRedis *redis.Client
)

func NewRedis(env *config.Config) *redis.Client {
	onceRedis.Do(func() {
		instanceRedis = redisDB(env)
	})

	log.Printf("Successfully connected to Redis database")

	return instanceRedis
}

func redisDB(env *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.SentryEndpoint,
		Password: "",
		DB:       0,
	})

	return rdb
}
