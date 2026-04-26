package guardrail

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// GuardRail is the main struct that holds the configuration and dependencies
type GuardRail struct {
	config    Config
	db        *gorm.DB
	redis     *redis.Client
	jwtSecret []byte
}

// New creates a new GuardRail middleware instance
// Usage:
//
//	app.Use(guardrail.New(guardrail.Config{
//	    DB: myDB,
//	    JWTSecret: "your-secret-key",
//	}))
func New(config Config) (*GuardRail, error) {
	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Set defaults
	config.setDefaults()

	gr := &GuardRail{
		config:    config,
		db:        config.DB,
		redis:     config.RedisClient,
		jwtSecret: []byte(config.JWTSecret),
	}

	return gr, nil
}

// Protect returns a Fiber middleware handler that validates JWT tokens
// This is the main middleware function that customers will use
func (gr *GuardRail) Protect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": gr.config.ErrorMessages.MissingToken,
			})
		}

		// Extract token from "Bearer <token>" format
		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid authorization header format. Expected: Bearer <token>",
			})
		}

		// Verify the JWT token
		claims, err := gr.verifyJWT(tokenStr)
		if err != nil {
			log.Printf("JWT verification failed: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": gr.config.ErrorMessages.InvalidToken,
			})
		}

		// Extract user information from claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid token claims",
			})
		}

		// Store user info in Fiber context for downstream handlers
		c.Locals("user_id", userID)

		// Store role if available
		if role, ok := claims["role"].(string); ok {
			c.Locals("role", role)
		}

		// Store tenant_id if multi-tenant is enabled
		if gr.config.EnableMultiTenant {
			if tenantID, ok := claims["tenant_id"].(string); ok {
				c.Locals("tenant_id", tenantID)
			}
		}

		// Store all claims for advanced use cases
		c.Locals("claims", claims)

		return c.Next()
	}
}

// ProtectWithRole returns middleware that validates JWT AND checks for specific roles
// Usage: app.Get("/admin", gr.ProtectWithRole("admin"), handler)
func (gr *GuardRail) ProtectWithRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First, run the standard protection
		if err := gr.Protect()(c); err != nil {
			return err
		}

		// Check if RBAC is enabled
		if !gr.config.EnableRBAC {
			return c.Next()
		}

		// Get role from context
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": gr.config.ErrorMessages.Forbidden,
			})
		}

		// Check if user's role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Insufficient permissions. Required role: " + strings.Join(allowedRoles, " or "),
		})
	}
}

// ApplicationKeyMiddleware validates application keys for multi-tenant apps
func (gr *GuardRail) ApplicationKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !gr.config.EnableMultiTenant {
			return c.Next()
		}

		key := c.Query("application_key")
		if key == "" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   true,
				"message": gr.config.ErrorMessages.MissingAppKey,
			})
		}

		// Check cache first if Redis is available
		var tenantID string
		cacheKey := "application_key:" + key

		if gr.redis != nil {
			ctx := context.Background()
			val, err := gr.redis.Get(ctx, cacheKey).Result()
			if err == nil {
				tenantID = val
				c.Locals("tenant_id", tenantID)
				c.Locals("application_key", key)
				return c.Next()
			}
		}

		// Verify from database
		var result struct {
			TenantID string
		}
		err := gr.db.Raw("SELECT tenant_id FROM application_tokens WHERE token = ? AND is_active = true", key).Scan(&result).Error
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": gr.config.ErrorMessages.InvalidAppKey,
			})
		}

		tenantID = result.TenantID
		if tenantID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": gr.config.ErrorMessages.InvalidAppKey,
			})
		}

		// Cache the result if Redis is available
		if gr.redis != nil {
			ctx := context.Background()
			gr.redis.Set(ctx, cacheKey, tenantID, 1*time.Hour)
		}

		c.Locals("tenant_id", tenantID)
		c.Locals("application_key", key)
		return c.Next()
	}
}

// verifyJWT validates a JWT token and returns its claims
func (gr *GuardRail) verifyJWT(tokenStr string) (jwt.MapClaims, error) {
	// Check if token is blacklisted (if Redis is available)
	if gr.redis != nil {
		ctx := context.Background()
		blacklisted, err := gr.redis.Exists(ctx, "blacklist:"+tokenStr).Result()
		if err == nil && blacklisted > 0 {
			return nil, fmt.Errorf("token is blacklisted")
		}

		// Check cache for valid tokens
		val, err := gr.redis.Get(ctx, "token:"+tokenStr).Result()
		if err == nil && val != "" {
			// Parse cached claims (simplified - in production, use proper serialization)
			token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return gr.jwtSecret, nil
			})
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				return claims, nil
			}
		}
	}

	// Parse and validate token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return gr.jwtSecret, nil
	})

	if err != nil {
		// Blacklist invalid tokens if Redis is available
		if gr.redis != nil {
			ctx := context.Background()
			gr.redis.Set(ctx, "blacklist:"+tokenStr, "invalid", 1*time.Hour)
		}
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Cache valid token if Redis is available
		if gr.redis != nil {
			if exp, ok := claims["exp"].(float64); ok {
				expiresAt := time.Unix(int64(exp), 0)
				ttl := time.Until(expiresAt)
				if ttl > 0 {
					ctx := context.Background()
					gr.redis.Set(ctx, "token:"+tokenStr, "valid", ttl)
				}
			}
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token or claims")
}

// GetUserID is a helper function to extract user_id from Fiber context
func GetUserID(c *fiber.Ctx) (string, bool) {
	userID, ok := c.Locals("user_id").(string)
	return userID, ok
}

// GetRole is a helper function to extract role from Fiber context
func GetRole(c *fiber.Ctx) (string, bool) {
	role, ok := c.Locals("role").(string)
	return role, ok
}

// GetTenantID is a helper function to extract tenant_id from Fiber context
func GetTenantID(c *fiber.Ctx) (string, bool) {
	tenantID, ok := c.Locals("tenant_id").(string)
	return tenantID, ok
}

// GetClaims is a helper function to extract all JWT claims from Fiber context
func GetClaims(c *fiber.Ctx) (jwt.MapClaims, bool) {
	claims, ok := c.Locals("claims").(jwt.MapClaims)
	return claims, ok
}
