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
