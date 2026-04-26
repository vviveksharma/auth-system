package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	guardrail "github.com/vviveksharma/auth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. Setup database (using SQLite for simplicity)
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.AutoMigrate(&guardrail.User{})

	// 2. Initialize GuardRail in ONE LINE
	gr, _ := guardrail.New(guardrail.Config{DB: db, JWTSecret: "my-secret-key"})

	// 3. Create Fiber app
	app := fiber.New()

	// 4. Protect routes with ONE LINE
	app.Get("/protected", gr.Protect(), func(c *fiber.Ctx) error {
		userID, _ := guardrail.GetUserID(c)
		return c.JSON(fiber.Map{"message": "Protected!", "user_id": userID})
	})

	// 5. Start server
	log.Fatal(app.Listen(":3000"))
}
