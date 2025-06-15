package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/controllers"
	"github.com/vviveksharma/auth/internal/middlewares"
)

func Routes(app *fiber.App, h *controllers.Handler) {
	app.Get("/health", h.Welcome)
	auth := app.Group("/auth")
	auth.Post("/register", h.CreateUser)
	user := app.Group("/user", middlewares.ExtractHeadersMiddleware())
	user.Get("/me", h.GetUserDetails)
	user.Put("/me", h.UpdateUserDetails)
}
