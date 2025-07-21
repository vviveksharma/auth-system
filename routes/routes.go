package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/vviveksharma/auth/internal/controllers"
	"github.com/vviveksharma/auth/internal/middlewares"
)

func Routes(app *fiber.App, h *controllers.Handler, client *redis.Client) {
	app.Get("/health", h.Welcome)
	auth := app.Group("/auth")
	user := app.Group("/user")
	role := app.Group("/roles")
	tenant := app.Group("/tenant")

	auth.Post("/register", h.CreateUser)
	auth.Post("/login", h.LoginUser)
	auth.Put("/refresh", middlewares.JWTMiddleware(), h.RefreshToken)

	user.Get("/me", middlewares.ExtractHeadersMiddleware(), h.GetUserDetails)
	user.Put("/me", middlewares.ExtractHeadersMiddleware(), h.UpdateUserDetails)
	user.Get("/:id", middlewares.ExtractRoleIdMiddleware(), h.GetUserByIdDetails)
	user.Put("/:id/roles", middlewares.ExtractRoleIdMiddleware(), h.AssignUserRole)

	role.Get("/", middlewares.ExtractRoleIdMiddleware(), h.ListAllRoles)
	role.Post("/", h.CreateCustomRole)
	role.Put("/permissions", h.UpdateRolePermission)
	role.Get("/verify", h.VerifyRole) //Internal service

	tenant.Post("/", h.CreateTenant)
	tenant.Post("/login", h.LoginTenant)
	tenant.Get("/tokens", middlewares.TenantMiddleWare(), h.ListTokens)
	tenant.Put("/tokens/:id", middlewares.TenantMiddleWare(), h.RevokeToken)
}
