package routes

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	_ "github.com/vviveksharma/auth/docs"
	"github.com/vviveksharma/auth/internal/controllers"
	tenantcontrollers "github.com/vviveksharma/auth/internal/controllers/tenantControllers"
	"github.com/vviveksharma/auth/internal/middlewares"
	"github.com/vviveksharma/auth/limiter"
)

func RoutesWithNewMiddleware(app *fiber.App, h *controllers.Handler, redisClient *redis.Client) {
	log.Println("📡 Setting up API Server Routes (Port 8080):")

	app.Get("/health", h.Welcome)
	log.Println("   ✅ GET  /health")

	app.Static("/docs", "./docs", fiber.Static{
		CacheDuration: 24 * time.Hour,
		MaxAge:        86400, // 24 hours
	})
	log.Println("   ✅ Static /docs")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:   "/docs/swagger.json",
		Title: "Auth System API",
	}))
	log.Println("   ✅ GET  /swagger/*")

	// Auth routes - mostly public or basic auth
	auth := app.Group("/auth")
	auth.Post("/", middlewares.PublicWithAppKey(), h.RegisterUser)
	log.Println("   ✅ POST /auth/")

	auth.Post("/login", middlewares.PublicWithAppKey(), h.LoginUser)
	log.Println("   ✅ POST /auth/login")

	// Apply middleware chain manually for individual routes
	refreshHandlers := append(middlewares.BasicAuthChain(), h.RefreshToken)
	auth.Put("/refresh", refreshHandlers...)
	log.Println("   ✅ PUT  /auth/refresh")

	logoutHandlers := append(middlewares.BasicAuthChain(), h.LogoutUser)
	auth.Put("/logout", logoutHandlers...)
	log.Println("   ✅ PUT  /auth/logout")

	// User routes - need full authentication and authorization
	user := app.Group("/users")
	user.Use(limiter.UserRateLimiter(redisClient))
	fullAuthHandlers := middlewares.FullAuthChain()
	for _, handler := range fullAuthHandlers {
		user.Use(handler)
	}
	user.Get("/me", h.GetUserDetails)
	log.Println("   ✅ GET  /user/me")

	user.Put("/me", h.UpdateUserDetails)
	log.Println("   ✅ PUT  /user/me")

	user.Get("/:id", h.GetUserByIdDetails)
	log.Println("   ✅ GET  /user/:id")

	user.Get("/", h.ListUsers)
	log.Println("   ✅ GET  /users/")

	user.Put("/:id/roles", h.AssignUserRole)
	log.Println("   ✅ PUT  /user/:id/roles")

	user.Delete("/:id", h.DeleteUser)
	log.Println("   ✅ DELETE /user/:id")

	// Special user routes that don't need authorization (just app key + auth)
	userPublic := app.Group("/user")
	userPublic.Post("/resetpassword", middlewares.PublicWithAppKey(), h.ResetUserPassword)
	log.Println("   ✅ POST /user/resetpassword")

	userPublic.Put("/setpassword", middlewares.PublicWithAppKey(), h.SetUserPassword)
	log.Println("   ✅ PUT  /user/setpassword")

	// Role routes - need full auth chain
	role := app.Group("/roles")
	role.Use(limiter.UserRateLimiter(redisClient))
	roleAuthHandlers := middlewares.FullAuthChain()
	for _, handler := range roleAuthHandlers {
		role.Use(handler)
	}
	role.Get("/", h.ListAllRoles)
	log.Println("   ✅ GET  /roles")

	role.Post("/", h.CreateCustomRole)
	log.Println("   ✅ POST /roles")

	role.Put("/:id/permissions", h.UpdateRolePermission)
	log.Println("   ✅ PUT  /roles/:id/permissions")

	role.Put("/enable/:id", h.EnableRole)
	log.Println("   ✅ PUT  /roles/enable/:id")

	role.Put("/disable/:id", h.DisableRole)
	log.Println("   ✅ PUT  /roles/disable/:id")

	role.Delete("/:id", h.DeleteCustomRole)
	log.Println("   ✅ DELETE /roles/:id")

	role.Get("/:id/permissions", h.GetRolePermissions)
	log.Println("   ✅ GET  /roles/:id/permissions")

	message := app.Group("/request")
	messageAuthHandlers := middlewares.FullAuthChain()
	for _, handler := range messageAuthHandlers {
		message.Use(handler)
	}

	message.Post("/", h.CreateRequest)
	log.Println("   ✅ POST  /request")

	message.Get("/status", h.GetRequestStatus)
	log.Println("   ✅ GET  /request/status")

	message.Get("/", h.GetMessages)
	log.Println("   ✅ GET  /request")

	log.Println("✅ API Server Routes setup complete!")
}

func TenantRoutes(app *fiber.App, h *tenantcontrollers.TenantHandler, redisClient *redis.Client) {
	log.Println("🖥️  Setting up UI Server Routes (Port 8081):")
	app.Get("/health", h.Welcome)
	log.Println("   ✅ GET  /health")
	tenant := app.Group("/tenant")
	tenant.Use(limiter.TenantRateLimiter(redisClient))
	tenant.Post("/", h.CreateTenant)
	log.Println("   ✅ POST /tenant/")
	tenant.Post("/login", h.LoginTenant)
	log.Println("   ✅ POST /tenant/login")
	// Protected tenant routes
	tenantProtected := tenant.Group("/")
	tenantProtected.Use(middlewares.TenantMiddleWare())
	tenantProtected.Get("/tokens", h.ListTokens)
	log.Println("   ✅ GET  /tenant/tokens")
	tenantProtected.Put("/tokens/:id", h.RevokeToken)
	log.Println("   ✅ PUT  /tenant/tokens/:id")
	tenantProtected.Post("/tokens", h.CreateToken)
	log.Println("   ✅ POST /tenant/tokens")
	tenantProtected.Get("/me", h.GetTenantDetails)
	log.Println("   ✅ GET  /tenant/me")
	tenantProtected.Post("/reset", h.ResetPassword)
	log.Println("   ✅ POST /tenant/reset")
	tenantProtected.Put("/setpassword", h.SetPassword)
	log.Println("   ✅ PUT  /tenant/setpassword")
	tenantProtected.Get("/dashboard", h.GetDashboardDetails)
	log.Println("   ✅ GET  /tenant/dashboard")
	tenantProtected.Get("/tokens/status", h.GetTokenDetailsStatus)
	log.Println("   ✅ GET  /tenant/tokens/status")
	tenantProtected.Get("/roles", h.ListRoles)
	log.Println("   ✅ GET  /tenant/roles")
	tenantProtected.Post("/roles", h.AddRole)
	log.Println("   ✅ POST /tenant/roles")
	tenantProtected.Get("/roles/persmissions", h.GetRolePermissions)
	log.Println("   ✅ GET  /tenant/roles/persmissions")
	tenantProtected.Put("/roles/enable", h.EnableRole)
	log.Println("   ✅ PUT  /tenant/roles/enable")
	tenantProtected.Put("/roles/disable", h.DisableRole)
	log.Println("   ✅ PUT  /tenant/roles/disable")
	tenantProtected.Delete("/", h.DeleteTenant)
	log.Println("   ✅ DELETE /tenant/")
	tenantProtected.Delete("/roles", h.DeleteRole)
	log.Println("   ✅ DELETE /tenant/roles")
	tenantProtected.Put("/roles/permissions", h.EditRolePermissions)
	log.Println("   ✅ PUT /tenant/roles/permissions")
	tenantProtected.Get("/messages", h.ListMessages)
	log.Println("   ✅ GET  /messages")
	tenantProtected.Put("/messages/approve", h.ApproveMessage)
	log.Println("   ✅ PUT  /tenant/messages/approve")
	tenantProtected.Put("/messages/reject", h.RejectMessage)
	log.Println("   ✅ PUT  /tenant/messages/reject")
	tenantProtected.Get("/users", h.ListUsers)
	log.Println("   ✅ GET  /tenant/users")
	tenantProtected.Put("/user/enable", h.EnableUser)
	log.Println("   ✅ PUT /tenant/user/enable")
	tenantProtected.Put("/user/disable", h.DisableUser)
	log.Println("   ✅ PUT  /tenant/user/disable")
	tenantProtected.Delete("/user", h.DeleteUser)
	log.Println("   ✅ DELETE  /tenant/user")
	log.Println("✅ UI Server Routes setup complete!")
}
