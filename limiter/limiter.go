package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func AuthRateLimiter(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		key := fmt.Sprintf("auth_limit:%s", ip)
		limit := 5
		expiry := 5 * time.Minute // 5 minutes block

		pipe := redisClient.Pipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, expiry)
		_, err := pipe.Exec(ctx)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Rate limiting failed"})
		}

		count := incr.Val()
		if count > int64(limit) {
			ttl, _ := redisClient.TTL(ctx, key).Result()
			return c.Status(429).JSON(fiber.Map{
				"error":       "Too many login attempts",
				"retry_after": fmt.Sprintf("%.0f seconds", ttl.Seconds()),
			})
		}

		return c.Next()
	}
}

func UserRateLimiter(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id") // From JWT middleware
		if userID == nil {
			return c.Next() // Skip if not authenticated
		}

		key := fmt.Sprintf("user_limit:%v", userID)
		limit := 20               // 50 requests
		expiry := 1 * time.Minute // per minute

		pipe := redisClient.Pipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, expiry)
		_, err := pipe.Exec(ctx)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Rate limiting failed"})
		}

		count := incr.Val()
		if count > int64(limit) {
			return c.Status(429).JSON(fiber.Map{
				"error": "User rate limit exceeded. Try again later.",
			})
		}

		return c.Next()
	}
}

func TenantRateLimiter(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Locals("tenant_id") // From middleware
		if tenantID == nil {
			return c.Next()
		}

		key := fmt.Sprintf("tenant_limit:%v", tenantID)
		limit := 50              // 100 requests
		expiry := 1 * time.Minute // per minute per tenant

		pipe := redisClient.Pipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, expiry)
		_, err := pipe.Exec(ctx)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Rate limiting failed"})
		}

		count := incr.Val()
		if count > int64(limit) {
			return c.Status(429).JSON(fiber.Map{
				"error": "Tenant rate limit exceeded.",
			})
		}

		return c.Next()
	}
}
