package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/controllers"
)

func Routes(app *fiber.App, h *controllers.Handler) {
	app.Get("/welcome", h.Welcome)
	user := app.Group("/user")
	user.Post("/", h.CreateUser)
}
