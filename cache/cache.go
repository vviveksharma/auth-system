package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var defaultClient ICacheClient

func NewCacheClient(client *redis.Client) ICacheClient {
	return &CacheClient{
		Client: client,
	}
}

// Init - Initialize the default cache client (call once at startup)
func Init(redisClient *redis.Client) {
	defaultClient = NewCacheClient(redisClient)
	log.Println("✅ Cache client initialized")
}

type CacheClient struct {
	Client *redis.Client
}

type ICacheClient interface {
	SetValue(key string, value string) error
	SetAPIValue(key string, value interface{}) error
	// Generic cache operations
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	Exists(key string) (bool, error)
}

func ConnectCache() *redis.Client {
	url := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	addr := url + ":" + port
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		DB:           0,
		PoolSize:     50,
		MinIdleConns: 10,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	ctx := context.Background()
	maxRetries := 5
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		if err := rdb.Ping(ctx).Err(); err == nil {
			log.Println("✅ Connected to Redis successfully")
			return rdb
		}

		log.Printf("⚠️ Failed to connect to Redis (attempt %d/%d)", i+1, maxRetries)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	log.Println("❌ Failed to connect to Redis after all retries")
	return nil
}

func (r *CacheClient) SetValue(key string, value string) error {
	err := r.Client.Set(ctx, key, value, time.Hour*24*365).Err()
	if err != nil {
		log.Println("failed to set value in Cache:", err)
	}
	return err
}

func (r *CacheClient) SetAPIValue(key string, value interface{}) error {
	err := r.Client.Set(ctx, key, value, 10*time.Minute).Err()
	if err != nil {
		log.Println("failed to set value in Cache:", err)
	}
	return err
}

// ==================== GENERIC CACHE OPERATIONS ====================

// Set - Generic function to cache any data with custom TTL
func (r *CacheClient) Set(key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := r.Client.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		log.Printf("❌ Failed to set cache key '%s': %v", key, err)
		return fmt.Errorf("redis set error: %w", err)
	}

	log.Printf("✅ Cached: %s (TTL: %v)", key, ttl)
	return nil
}

// Get - Generic function to retrieve cached data
func (r *CacheClient) Get(key string, dest interface{}) error {
	data, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss: key '%s' not found", key)
	}
	if err != nil {
		return fmt.Errorf("redis get error: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	log.Printf("✅ Cache HIT: %s", key)
	return nil
}

// Delete - Remove key from cache
func (r *CacheClient) Delete(key string) error {
	if err := r.Client.Del(ctx, key).Err(); err != nil {
		log.Printf("⚠️ Failed to delete cache key '%s': %v", key, err)
		return err
	}

	log.Printf("🗑️ Cache invalidated: %s", key)
	return nil
}

// Exists - Check if key exists in cache
func (r *CacheClient) Exists(key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// ==================== HELPER FUNCTIONS (Use these directly) ====================

// Set - Cache any data with custom TTL (no initialization needed)
func Set(key string, value interface{}, ttl time.Duration) error {
	if defaultClient == nil {
		return fmt.Errorf("cache not initialized, call cache.Init() first")
	}
	return defaultClient.Set(key, value, ttl)
}

// Get - Retrieve cached data (no initialization needed)
func Get(key string, dest interface{}) error {
	if defaultClient == nil {
		return fmt.Errorf("cache not initialized, call cache.Init() first")
	}
	return defaultClient.Get(key, dest)
}

// Delete - Remove key from cache (no initialization needed)
func Delete(key string) error {
	if defaultClient == nil {
		return fmt.Errorf("cache not initialized, call cache.Init() first")
	}
	return defaultClient.Delete(key)
}

// Exists - Check if key exists (no initialization needed)
func Exists(key string) (bool, error) {
	if defaultClient == nil {
		return false, fmt.Errorf("cache not initialized, call cache.Init() first")
	}
	return defaultClient.Exists(key)
}
