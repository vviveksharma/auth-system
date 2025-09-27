package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	_ "github.com/vviveksharma/auth/docs"
	"github.com/vviveksharma/auth/internal/controllers"
	"github.com/vviveksharma/auth/internal/middlewares"
)

// Example of how to use the new middleware architecture
func RoutesWithNewMiddleware(app *fiber.App, h *controllers.Handler, client *redis.Client) {
	app.Get("/health", h.Welcome)
	app.Static("/docs", "./docs", fiber.Static{
		CacheDuration: 24 * time.Hour,
		MaxAge:        86400, // 24 hours
	})

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:   "/docs/swagger.json",
		Title: "Auth System API",
	}))

	// Auth routes - mostly public or basic auth
	auth := app.Group("/auth")
	auth.Post("/", middlewares.PublicWithAppKey(), h.RegisterUser)
	auth.Post("/login", middlewares.PublicWithAppKey(), h.LoginUser)

	// Apply middleware chain manually for individual routes
	refreshHandlers := append(middlewares.BasicAuthChain(), h.RefreshToken)
	auth.Put("/refresh", refreshHandlers...)

	logoutHandlers := append(middlewares.BasicAuthChain(), h.LogoutUser)
	auth.Put("/logout", logoutHandlers...)

	// User routes - need full authentication and authorization
	user := app.Group("/user")
	fullAuthHandlers := middlewares.FullAuthChain()
	for _, handler := range fullAuthHandlers {
		user.Use(handler)
	}
	user.Get("/me", h.GetUserDetails)
	user.Put("/me", h.UpdateUserDetails)
	user.Get("/:id", h.GetUserByIdDetails)
	user.Put("/:id/roles", h.AssignUserRole)
	user.Delete("/:id", h.DeleteUser)

	// Special user routes that don't need authorization (just app key + auth)
	userPublic := app.Group("/user")
	userPublic.Post("/resetpassword", middlewares.PublicWithAppKey(), h.ResetUserPassword)
	userPublic.Put("/setpassword", middlewares.PublicWithAppKey(), h.SetUserPassword)

	// Role routes - need full auth chain
	role := app.Group("/roles")
	roleAuthHandlers := middlewares.FullAuthChain()
	for _, handler := range roleAuthHandlers {
		role.Use(handler)
	}
	role.Get("/", h.ListAllRoles)
	role.Post("/", h.CreateCustomRole)
	role.Put("/:id/permissions", h.UpdateRolePermission)
	role.Get("/verify", h.VerifyRole)
	role.Put("/enable/:id", h.EnableRole)
	role.Put("/disable/:id", h.DisableRole)
	role.Delete("/:id", h.DeleteCustomRole)
	role.Get("/:id/permissions", h.GetRolePermissions)

	// Tenant routes - custom middleware for tenant-specific logic
	tenant := app.Group("/tenant")
	tenant.Post("/", h.CreateTenant)     // Public registration
	tenant.Post("/login", h.LoginTenant) // Public login

	// Protected tenant routes
	tenantProtected := tenant.Group("/")
	tenantProtected.Use(middlewares.TenantMiddleWare()) // Your existing tenant middleware
	tenantProtected.Get("/tokens", h.ListTokens)
	tenantProtected.Put("/tokens/:id", h.RevokeToken)
	tenantProtected.Post("/tokens", h.CreateToken)
	tenantProtected.Post("/reset", h.ResetPassword)
	tenantProtected.Put("/setpassword", h.SetPassword)
	tenantProtected.Get("/dashboard", h.GetDashboardDetails)
}
