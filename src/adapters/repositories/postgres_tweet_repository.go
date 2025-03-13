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
	"strings"
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
	/*ctx := context.Background()
	cacheKey := generateCacheKey(userIDs, cursor, limit)

	if tweets, err := r.getTweetsFromCache(ctx, cacheKey); err == nil {
		log.Println("Tweets retrieved from cache")
		return tweets, nil
	} else {
		log.Printf("Error retrieving data from cache: %v", err)
	}*/

	tweets, err := r.getTweetsFromDB(userIDs, cursor, limit)
	if err != nil {
		return nil, err
	}

	/*if err = r.storeTweetsInCache(ctx, cacheKey, tweets); err != nil {
		log.Printf("Error storing in cache: %v", err)
	}*/

	return tweets, nil
}

func (r *PostgresTweetRepository) getTweetsFromCache(ctx context.Context, cacheKey string) ([]*entities.Tweet, error) {
	cachedData, err := r.cache.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var tweets []*entities.Tweet
	if err := json.Unmarshal([]byte(cachedData), &tweets); err != nil {
		return nil, err
	}

	return tweets, nil
}

func (r *PostgresTweetRepository) getTweetsFromDB(userIDs []uuid.UUID, cursor *uuid.UUID, limit int) ([]*entities.Tweet, error) {
	var tweets []*entities.Tweet
	query := r.db.Where("user_id IN ?", userIDs)

	if cursor != nil {
		var cursorTweet entities.Tweet
		if err := r.db.First(&cursorTweet, *cursor).Error; err != nil {
			return nil, err
		}
		query = query.Where("created_at < ?", cursorTweet.CreatedAt)
	}

	log.Println("Retrieving data from database")
	if err := query.Order("created_at desc").Limit(limit).Find(&tweets).Error; err != nil {
		return nil, err
	}

	return tweets, nil
}

func (r *PostgresTweetRepository) storeTweetsInCache(ctx context.Context, cacheKey string, tweets []*entities.Tweet) error {
	serializedData, err := json.Marshal(tweets)
	if err != nil {
		return err
	}

	if err := r.cache.Set(ctx, cacheKey, serializedData, time.Hour).Err(); err != nil {
		return err
	}

	log.Println("Results stored in cache")

	return nil
}

func generateCacheKey(userIDs []uuid.UUID, cursor *uuid.UUID, limit int) string {
	var keyBuilder strings.Builder
	keyBuilder.WriteString("tweets:")

	for _, id := range userIDs {
		keyBuilder.WriteString(id.String())
		keyBuilder.WriteString(":")
	}

	if cursor != nil {
		keyBuilder.WriteString(cursor.String())
		keyBuilder.WriteString(":")
	}

	keyBuilder.WriteString(strconv.Itoa(limit))

	return keyBuilder.String()
}
