package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewCacheClient(client *redis.Client) ICacheClient {
	return &CacheClient{
		Client: client,
	}
}

type CacheClient struct {
	Client *redis.Client
}

type ICacheClient interface {
	SetValue(key string, value string) error
	SetAPIValue(key string, value interface{}) error
}

func ConnectCache() *redis.Client {
	var ctx = context.Background()
	CacheClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	_, err := CacheClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Cache: %v", err)
	}

	log.Println("Cache connected successfully")
	return CacheClient
}

func (r *CacheClient) SetValue(key string, value string) error {
	var ctx = context.Background()
	err := r.Client.Set(ctx, key, value, time.Hour*24*365).Err()
	if err != nil {
		log.Println("failed to set value in Cache:", err)
	}
	return err
}
func (r *CacheClient) SetAPIValue(key string, value interface{}) error {
	var ctx = context.Background()
	err := r.Client.Set(ctx, key, value, 10*time.Minute).Err()
	if err != nil {
		log.Println("failed to set value in Cache:", err)
	}
	return err
}
