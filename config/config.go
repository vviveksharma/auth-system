// All the configuration is set here only
package config

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/controllers"
	"github.com/vviveksharma/auth/internal/services"
	"github.com/vviveksharma/auth/routes"
)

func Init() {
	app := fiber.New()

	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	db.ConnectDB()

	userService, err := services.NewUserService()
	if err != nil {
		log.Println("error while starting the user-service: ", err)
	}
	handler, err := controllers.NewHandler(userService)
	if err != nil {
		log.Println("error while starting the handler: ", err)
	}
	routes.Routes(app, handler)
	err = app.Listen(":8080")
	if err != nil {
		log.Fatal("error while starting the server: ", err)
	}

}
