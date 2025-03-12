package database

import (
	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	onceRedis     sync.Once
	instanceRedis *redis.Client
)

func NewRedis() *redis.Client {
	onceRedis.Do(func() {
		instanceRedis = redisDB()
	})

	return instanceRedis
}

func redisDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}
