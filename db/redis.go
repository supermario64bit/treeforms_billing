package db

import (
	"context"
	"os"
	"strconv"
	"time"
	"treeforms_billing/logger"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func GetRedis() *redis.Client {

	if redisClient != nil {
		return redisClient
	}

	// Connect to Redis
	redisDbStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDbStr)
	if err != nil {
		logger.HighlightedDanger("Unable to get redis db. Message: " + err.Error())
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // Redis address
		Password: os.Getenv("REDIS_PASS"), // No password by default
		DB:       redisDB,                 // Default DB
	})

	return redisClient
}

func SetRedisCache(key string, value interface{}, expTime time.Duration) error {
	if redisClient != nil {
		redisClient = GetRedis()
	}

	return redisClient.Set(context.Background(), key, value, expTime).Err()
}

func GetFromRedisCache(key string) (string, error) {
	if redisClient != nil {
		redisClient = GetRedis()
	}

	return redisClient.Get(context.Background(), key).Result()
}
