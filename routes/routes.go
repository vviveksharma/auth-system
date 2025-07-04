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
	role := app.Group("/roles")

	auth.Post("/register", h.CreateUser)
	auth.Post("/login", h.LoginUser)

	user.Get("/me", middlewares.ExtractHeadersMiddleware(), h.GetUserDetails)
	user.Put("/me", middlewares.ExtractHeadersMiddleware(), h.UpdateUserDetails)
	user.Get("/:id", middlewares.ExtractRoleIdMiddleware(), h.GetUserByIdDetails)
	user.Put("/:id/roles", middlewares.ExtractRoleIdMiddleware(), h.AssignUserRole)

	role.Get("/", middlewares.ExtractRoleIdMiddleware(), h.ListAllRoles)
	role.Get("/verify", h.VerifyRole)
}
