package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/controllers"
	"github.com/vviveksharma/auth/internal/middlewares"
)

func Routes(app *fiber.App, h *controllers.Handler) {
	app.Get("/health", h.Welcome)
	auth := app.Group("/auth")
	user := app.Group("/user")
	//Routes
	auth.Post("/register", h.CreateUser)
	user.Get("/me", h.GetUserDetails, middlewares.ExtractHeadersMiddleware())
	user.Put("/me", h.UpdateUserDetails, middlewares.ExtractHeadersMiddleware())
	user.Get("/:id", h.GetUserByIdDetails, middlewares.ExtractAdminIdMiddleware())
}
