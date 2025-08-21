package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	_ "github.com/vviveksharma/auth/docs"
	"github.com/vviveksharma/auth/internal/controllers"
	"github.com/vviveksharma/auth/internal/middlewares"
)

func Routes(app *fiber.App, h *controllers.Handler, client *redis.Client) {
	app.Get("/health", h.Welcome)
	app.Static("/docs", "./docs")
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/docs/swagger.json",
	}))
	auth := app.Group("/auth")
	user := app.Group("/user")
	role := app.Group("/roles")
	tenant := app.Group("/tenant")

	auth.Post("/login", h.LoginUser)
	auth.Put("/refresh", middlewares.JWTMiddleware(), h.RefreshToken)
	// auth.Post("/invite")
	// auth.Put("/resetPassword")

	user.Get("/me", middlewares.ExtractHeadersMiddleware(), h.GetUserDetails)
	user.Post("/", middlewares.GetTenantFromToken(), h.RegisterUser)
	user.Put("/me", middlewares.ExtractHeadersMiddleware(), h.UpdateUserDetails)
	user.Get("/:id", middlewares.ExtractRoleIdMiddleware(), h.GetUserByIdDetails)
	user.Put("/:id/roles", middlewares.ExtractRoleIdMiddleware(), h.AssignUserRole)
	user.Post("/resetpassword", h.ResetUserPassword)
	user.Put("/setpassword", h.SetUserPassword)

	role.Get("/", middlewares.ExtractRoleIdMiddleware(), h.ListAllRoles)
	role.Post("/", h.CreateCustomRole)
	role.Put("/permissions", h.UpdateRolePermission)
	role.Get("/verify", h.VerifyRole)

	tenant.Post("/", h.CreateTenant)
	tenant.Post("/login", h.LoginTenant)
	tenant.Get("/tokens", middlewares.TenantMiddleWare(), h.ListTokens)
	tenant.Put("/tokens/:id", middlewares.TenantMiddleWare(), h.RevokeToken)
	tenant.Post("/register", middlewares.TenantMiddleWare(), h.CreateUser)
	tenant.Post("/tokens", middlewares.TenantMiddleWare(), h.CreateToken)
	tenant.Post("/reset", middlewares.TenantMiddleWare(), h.ResetPassword)
	tenant.Put("/setpassword", middlewares.TenantMiddleWare(), h.SetPassword)
	tenant.Get("/users", middlewares.TenantMiddleWare(), h.ListUsers)
	tenant.Get("/roles")
}
