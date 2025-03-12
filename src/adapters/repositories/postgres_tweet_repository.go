package repositories

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rodrinoblega/microblogging/src/entities"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type PostgresTweetRepository struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewPostgresTweetRepository(db *gorm.DB, cache *redis.Client) *PostgresTweetRepository {
	return &PostgresTweetRepository{db: db, cache: cache}
}

func (r *PostgresTweetRepository) Save(tweet *entities.Tweet) error {
	return r.db.Create(tweet).Error
}

func (r *PostgresTweetRepository) GetTweetsByUsers(userIDs []uuid.UUID, cursor *uuid.UUID, limit int) ([]*entities.Tweet, error) {
	var tweets []*entities.Tweet
	ctx := context.Background()

	cacheKey := generateCacheKey(userIDs, cursor, limit)

	cachedTweets, err := r.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		if err = json.Unmarshal([]byte(cachedTweets), &tweets); err == nil {
			log.Printf("Retrieving tweets from cache")
			return tweets, nil
		}
	}

	query := r.db.Where("user_id IN ?", userIDs)

	if cursor != nil {
		var cursorTweet entities.Tweet
		if err := r.db.First(&cursorTweet, *cursor).Error; err != nil {
			return nil, err
		}
		query = query.Where("created_at < ?", cursorTweet.CreatedAt)
	}

	log.Printf("Retrieving tweets from database")
	err = query.Order("created_at desc").Limit(limit).Find(&tweets).Error
	if err != nil {
		return nil, err
	}

	serializedTweets, err := json.Marshal(tweets)
	if err == nil {
		log.Printf("Catching database result")
		r.cache.Set(ctx, cacheKey, serializedTweets, time.Hour).Err()
	}

	return tweets, nil
}

func generateCacheKey(userIDs []uuid.UUID, cursor *uuid.UUID, limit int) string {
	key := "tweets:"
	for _, id := range userIDs {
		key += id.String() + ":"
	}
	if cursor != nil {
		key += cursor.String() + ":"
	}
	key += strconv.Itoa(limit)
	return key
}
