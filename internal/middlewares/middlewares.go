package middlewares

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ExtractHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("userId")

		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "UserId header is required",
			})
		}
		c.Locals("userId", userID)
		fmt.Println("the userId from the middleware: ", userID)
		return c.Next()
	}
}

func ExtractRoleIdMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleId := c.Get("roleId")
		if roleId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "roleId header is required",
			})
		}
		c.Locals("roleId", roleId)
		return c.Next()
	}
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := VerifyJWT(tokenStr)
		if err != nil {
			fmt.Println("the error while verifying the token: ", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		c.Locals("authClaims", claims)
		return c.Next()
	}
}
