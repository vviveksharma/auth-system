// All the configuration is set here only
package config

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/initsetup"
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

	//Create Roles in the db
	initsetup.InitRoles()

	userService, err := services.NewUserService()
	if err != nil {
		log.Fatalln("error while starting the user-service: ", err)
	}
	roleService, err := services.NewRoleService()
	if err != nil {
		log.Fatalln("error while starting the role-service: ", err)
	}
	authService, err := services.NewAuthService()
	if err != nil {
		log.Fatalln("error while starting the auth-service: ", err)
	}

	handler, err := controllers.NewHandler(userService, roleService, authService)
	if err != nil {
		log.Fatalln("error while starting the handler: ", err)
	}

	//Starting the server
	routes.Routes(app, handler)
	err = app.Listen(":8080")
	if err != nil {
		log.Fatal("error while starting the server: ", err)
	}

}
