package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	guardrail "github.com/vviveksharma/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Setup your database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=myapp port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the User table
	db.AutoMigrate(&guardrail.User{})

	// 2. Initialize GuardRail with your config
	gr, err := guardrail.New(guardrail.Config{
		DB:                db,
		JWTSecret:         os.Getenv("JWT_SECRET"),
		EnableRBAC:        true,
		EnableMultiTenant: false,
	})
	if err != nil {
		log.Fatal("Failed to initialize GuardRail:", err)
	}

	// 3. Create auth service for registration/login
	authService := gr.NewAuthService()

	// 4. Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New())

	// 5. Public routes (no authentication required)
	app.Post("/register", func(c *fiber.Ctx) error {
		var req guardrail.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		response, err := authService.Register(req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(response)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var req guardrail.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		response, err := authService.Login(req)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(response)
	})

	app.Post("/refresh", func(c *fiber.Ctx) error {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		response, err := authService.RefreshToken(req.RefreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(response)
	})

	// 6. Protected routes - Use GuardRail middleware
	// Simple protection (any authenticated user)
	app.Get("/profile", gr.Protect(), func(c *fiber.Ctx) error {
		userID, _ := guardrail.GetUserID(c)
		role, _ := guardrail.GetRole(c)

		return c.JSON(fiber.Map{
			"message": "This is a protected route",
			"user_id": userID,
			"role":    role,
		})
	})

	// Role-based protection (only admin users)
	app.Get("/admin/dashboard", gr.ProtectWithRole("admin"), func(c *fiber.Ctx) error {
		userID, _ := guardrail.GetUserID(c)

		return c.JSON(fiber.Map{
			"message": "Welcome to admin dashboard",
			"user_id": userID,
		})
	})

	// Multiple roles allowed
	app.Get("/moderator/panel", gr.ProtectWithRole("admin", "moderator"), func(c *fiber.Ctx) error {
		userID, _ := guardrail.GetUserID(c)
		role, _ := guardrail.GetRole(c)

		return c.JSON(fiber.Map{
			"message": "Moderator panel",
			"user_id": userID,
			"role":    role,
		})
	})

	// 7. Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
