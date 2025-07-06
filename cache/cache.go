package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(client *redis.Client) IRedisClient {
	return &RedisClient{
		Client: client,
	}
}

type RedisClient struct {
	Client *redis.Client
}

type IRedisClient interface {
	SetValue(key string, value string) error
	SetAPIValue(key string, value interface{}) error
}

func ConnectCache() *redis.Client {
	var ctx = context.Background()
	RedisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return RedisClient
}

func (r *RedisClient) SetValue(key string, value string) error {
	var ctx = context.Background()
	err := r.Client.Set(ctx, key, value, time.Hour*24*365).Err()
	if err != nil {
		log.Println("failed to set value in Redis:", err)
	}
	return err
}
func (r *RedisClient) SetAPIValue(key string, value interface{}) error {
	var ctx = context.Background()
	err := r.Client.Set(ctx, key, value, 10*time.Minute).Err()
	if err != nil {
		log.Println("failed to set value in Redis:", err)
	}
	return err
}
