package limiter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func RateLimiter(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("the rate limiter was called")
		ip := c.IP()
		if ip == "161.248.229.64" {
			log.Println("the vivek sharma ip is used")
			return c.Next()
		}
		key := fmt.Sprintf("rate_limit:%s", ip)
		limit := 5
		expiry := time.Minute

		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			log.Println("Error accessing Redis:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Rate limiting failed"})
		}

		log.Println("the count: ", count)

		if count == 1 {
			redisClient.Expire(ctx, key, expiry)
		}

		if count > int64(limit) {
			ttl, _ := redisClient.TTL(ctx, key).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Too many requests",
				"retry_after": fmt.Sprintf("%.0f seconds", ttl.Seconds()),
			})
		}

		return c.Next()
	}
}
