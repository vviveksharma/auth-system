package config

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"github.com/vviveksharma/auth/cache"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/initsetup"
	"github.com/vviveksharma/auth/internal/controllers"
	tenantcontrollers "github.com/vviveksharma/auth/internal/controllers/tenantControllers"
	"github.com/vviveksharma/auth/internal/services"
	tenantservices "github.com/vviveksharma/auth/internal/services/tenant-services"
	"github.com/vviveksharma/auth/queue"
	"github.com/vviveksharma/auth/routes"
)

// SharedResources holds all shared dependencies
type SharedResources struct {
	RedisClient *redis.Client
	QConn       *amqp.Connection
	Queue       amqp.Queue
}

var (
	sharedResources *SharedResources
	once            sync.Once
	initError       error
)

// InitializeSharedResources initializes resources once (thread-safe)
func InitializeSharedResources() (*SharedResources, error) {
	once.Do(func() {
		// Load environment variables
		if err := godotenv.Load("./.env"); err != nil {
			log.Println("⚠️  Warning: Error loading .env file", err)
		}
		db.ConnectDB()
		redisClient := cache.ConnectCache()
		if redisClient == nil {
			initError = fmt.Errorf("failed to connect to Redis: client is nil")
			return
		}

		// Ping Redis with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := redisClient.Ping(ctx).Err(); err != nil {
			initError = fmt.Errorf("failed to ping Redis: %w", err)
			return
		}
		cache.Init(redisClient)

		log.Println("✅ Redis connected successfully")

		// Connecting to the queue
		Iqueue, err := queue.NewQueueRequest()
		if err != nil {
			initError = fmt.Errorf("failed to intialise the queue on the startup: %w", err)
			return
		}
		qConn, err := Iqueue.Connect()
		if err != nil {
			initError = fmt.Errorf("failed to connect to the queue on the startup: %w", err)
			return
		}
		queue, err := Iqueue.DeclareQueue(qConn)
		if err != nil {
			initError = fmt.Errorf("failed to declare the queue on the startup: %w", err)
			return
		}

		// Start consuming message continously
		go func() {
			log.Println("🚀 Starting queue consumer for role requests...")
			if err := Iqueue.ConsumeMessages(qConn, queue); err != nil {
				log.Printf("❌ Queue consumer error: %v", err)
			}
		}()

		// Initialize roles
		initsetup.InitSetup()

		sharedResources = &SharedResources{
			RedisClient: redisClient,
			QConn:       qConn,
			Queue:       queue,
		}
	})

	return sharedResources, initError
}

// CreateAPIServer creates the API server (port 8080)
func CreateAPIServer(resources *SharedResources) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Auth System - API Server",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	userService, err := services.NewUserService()
	if err != nil {
		log.Fatalln("❌ Error while starting the user-service: ", err)
	}
	roleService, err := services.NewRoleService()
	if err != nil {
		log.Fatalln("❌ Error while starting the role-service: ", err)
	}
	authService, err := services.NewAuthService()
	if err != nil {
		log.Fatalln("❌ Error while starting the auth-service: ", err)
	}
	messageService, err := services.NewMessageService(sharedResources.Queue, sharedResources.QConn)
	if err != nil {
		log.Fatalln("❌ Error while starting the message-service: ", err)
	}

	handler, err := controllers.NewHandler(userService, roleService, authService, messageService)
	if err != nil {
		log.Fatalln("❌ Error while starting the handler: ", err)
	}

	routes.RoutesWithNewMiddleware(app, handler, resources.RedisClient)
	return app
}

// CreateUIServer creates the UI server (port 8081)
func CreateUIServer(resources *SharedResources) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Auth System - UI Server",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error":       true,
				"message":     err.Error(),
				"status_code": code,
			})
		},
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	tenantRoleService, err := tenantservices.NewTenantRoleService()
	if err != nil {
		log.Fatalln("❌ Error while starting the tenant-role-service: ", err)
	}
	tenantUserService, err := tenantservices.NewTenantUserService()
	if err != nil {
		log.Fatalln("❌ Error while starting the tenant-user-service: ", err)
	}
	tenantService, err := tenantservices.NewTenantService()
	if err != nil {
		log.Fatalln("❌ Error while starting the tenant-service: ", err)
	}
	tenantMessageService, err := tenantservices.NewTenantMessageService()
	if err != nil {
		log.Fatalln("❌ Error while starting the tenant-service: ", err)
	}

	tenantHandler, err := tenantcontrollers.NewTenantHandler(
		tenantUserService,
		tenantRoleService,
		tenantService,
		tenantMessageService,
	)
	if err != nil {
		log.Fatalln("❌ Error while starting tenant handler: ", err)
	}

	routes.TenantRoutes(app, tenantHandler, resources.RedisClient)
	return app
}

// InitAPIOnly starts only the API server (port 8080)
func InitAPIOnly() {
	log.Println("🚀 Initializing Auth System - API Server Only...")

	resources, err := InitializeSharedResources()
	if err != nil {
		log.Fatalf("❌ Error while initializing shared resources: %v", err)
	}

	log.Println("✅ Shared resources initialized successfully")
	log.Println("📡 Starting API Server on port 8080...")

	app := CreateAPIServer(resources)
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("❌ API Server failed to start: %v", err)
	}
}

// InitUIOnly starts only the UI/Tenant server (port 8081)
func InitUIOnly() {
	log.Println("🚀 Initializing Auth System - UI Server Only...")

	resources, err := InitializeSharedResources()
	if err != nil {
		log.Fatalf("❌ Error while initializing shared resources: %v", err)
	}

	log.Println("✅ Shared resources initialized successfully")
	log.Println("🖥️  Starting UI Server on port 8081...")

	app := CreateUIServer(resources)
	if err := app.Listen(":8081"); err != nil {
		log.Fatalf("❌ UI Server failed to start: %v", err)
	}
}

// Init starts both servers (default behavior for local development)
func Init() {
	log.Println("🚀 Initializing Auth System...")

	resources, err := InitializeSharedResources()
	if err != nil {
		log.Fatalf("❌ Error while initializing shared resources: %v", err)
	}

	log.Println("✅ Shared resources initialized successfully")
	log.Println("🚀 Starting both servers...")

	var wg sync.WaitGroup
	wg.Add(2)

	// Start API Server
	go func() {
		defer wg.Done()
		log.Println("📡 Starting API Server on port 8080...")
		app := CreateAPIServer(resources)
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("❌ API Server failed to start: %v", err)
		}
	}()

	// Start UI Server
	go func() {
		defer wg.Done()
		log.Println("🖥️  Starting UI Server on port 8081...")
		app := CreateUIServer(resources)
		if err := app.Listen(":8081"); err != nil {
			log.Fatalf("❌ UI Server failed to start: %v", err)
		}
	}()

	log.Println("✅ Both servers are starting concurrently...")
	log.Println("   📡 API Server: http://localhost:8080")
	log.Println("   🖥️  UI Server:  http://localhost:8081")

	wg.Wait()
	log.Println("⚠️  Both servers have stopped")
}
