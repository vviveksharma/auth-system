package middlewares

import "github.com/gofiber/fiber/v2"

func ExtractHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("userId")

		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "UserId header is required",
			})
		}
		c.Locals("userId", userID)
		return c.Next()
	}
}

func ExtractAdminIdMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleId := c.Get("roleId")
		if roleId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "UserId header is required",
			})
		}
		c.Locals("roleId", roleId)
		return c.Next()
	}
}
