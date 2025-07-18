// All the configuration is set here only
package config

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/vviveksharma/auth/cache"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/initsetup"
	"github.com/vviveksharma/auth/internal/controllers"
	"github.com/vviveksharma/auth/internal/services"
	"github.com/vviveksharma/auth/limiter"
	"github.com/vviveksharma/auth/routes"
)

func Init() {
	app := fiber.New()

	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	db.ConnectDB()

	//Connecting to cache
	client := cache.ConnectCache()

	//Create Roles in the db
	initsetup.InitRoles()

	app.Use(limiter.RateLimiter(client))

	userService, err := services.NewUserService()
	if err != nil {
		log.Fatalln("error while starting the user-service: ", err)
	}
	roleService, err := services.NewRoleService()
	if err != nil {
		log.Fatalln("error while starting the role-service: ", err)
	}
	authService, err := services.NewAuthService(client)
	if err != nil {
		log.Fatalln("error while starting the auth-service: ", err)
	}
	tenantService, err := services.NewTenantService()
	if err != nil {
		log.Fatalln("error while starting the tenant-service: ", err)
	}

	handler, err := controllers.NewHandler(userService, roleService, authService, tenantService)
	if err != nil {
		log.Fatalln("error while starting the handler: ", err)
	}

	//Starting the server
	routes.Routes(app, handler, client)
	err = app.Listen(":8080")
	if err != nil {
		log.Fatal("error while starting the server: ", err)
	}

}
