package middlewares

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/repo"
	"github.com/vviveksharma/auth/models"
)

// ApplicationKeyMiddleware validates the application key for tenant verification
func ApplicationKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Query("application_key")
		if key == "" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ServiceResponse{
				Code:    fiber.StatusUnprocessableEntity,
				Message: "Application key is required",
			})
		}

		// Verify the application key and get tenant info
		tokenRepo, err := repo.NewTokenRepository(db.DB)
		if err != nil {
			log.Printf("Failed to initialize token repository: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ServiceResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Internal server error",
			})
		}

		isValid, tenantID, err := tokenRepo.VerifyApplicationToken(key)
		if err != nil {
			log.Printf("Application key verification failed: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(models.ServiceResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid application key",
			})
		}

		if !isValid {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ServiceResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid application key",
			})
		}
		// Store tenant info for downstream middleware/handlers
		c.Locals("tenant_id", tenantID)
		c.Locals("application_key", key)

		return c.Next()
	}
}

// AuthenticationMiddleware handles JWT token verification
func AuthenticationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ServiceResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Authorization header is required",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := VerifyJWT(tokenStr)
		if err != nil {
			log.Printf("JWT verification failed: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(models.ServiceResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid or expired token. Please login again",
			})
		}

		// Store auth claims for downstream middleware/handlers
		c.Locals("authClaims", claims)
		c.Locals("token", tokenStr)

		return c.Next()
	}
}

// AuthorizationMiddleware handles role-based access control
func AuthorizationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("authClaims").(jwt.MapClaims)
		if !ok {
			log.Printf("No auth claims found in context")
			return c.Status(fiber.StatusUnauthorized).JSON(models.ServiceResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Authentication required",
			})
		}

		roleID, ok := claims["role_id"].(string)
		if !ok {
			log.Printf("Invalid role_id claim type in JWT")
			return c.Status(fiber.StatusUnauthorized).JSON(models.ServiceResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid token claims",
			})
		}

		

		hasAccess, err := VerifyRoleRouteMapping(roleID, c.Path(), c.Method())
		if err != nil {
			log.Printf("Role route verification failed for role %s on path %s: %v", roleID, c.Path(), err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ServiceResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Authorization check failed",
			})
		}

		if !hasAccess {
			log.Printf("Access denied for role %s on path %s", roleID, c.Path())
			return c.Status(fiber.StatusForbidden).JSON(models.ServiceResponse{
				Code:    fiber.StatusForbidden,
				Message: "Insufficient permissions to access this resource",
			})
		}

		// Store role info for downstream handlers
		c.Locals("role_id", roleID)

		return c.Next()
	}
}

// FullAuthChain combines all authentication and authorization steps
// Use this for routes that need complete auth flow
func FullAuthChain() []fiber.Handler {
	return []fiber.Handler{
		ApplicationKeyMiddleware(),
		AuthenticationMiddleware(),
		AuthorizationMiddleware(),
	}
}

// BasicAuthChain for routes that only need app key and JWT verification
func BasicAuthChain() []fiber.Handler {
	return []fiber.Handler{
		ApplicationKeyMiddleware(),
		AuthenticationMiddleware(),
	}
}

// PublicWithAppKey for routes that only need application key validation
func PublicWithAppKey() fiber.Handler {
	return ApplicationKeyMiddleware()
}
