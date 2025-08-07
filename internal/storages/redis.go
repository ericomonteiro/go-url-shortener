package storages

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// RedisClient represents a Redis database client
var RedisClient *redis.Client
var ctx = context.Background()

// NewRedis initializes a new Redis client
func NewRedis() (*redis.Client, error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // default Redis address
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	// Test the connection
	err := RedisClient.Ping(ctx).Err()
	if err != nil {
		log.Printf("Error connecting to Redis: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to Redis")
	return RedisClient, nil
}

// CloseRedis closes the Redis connection
func CloseRedis() {
	if RedisClient != nil {
		_ = RedisClient.Close()
	}
}
