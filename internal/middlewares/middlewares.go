package middlewares

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/db"
	model "github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	"github.com/vviveksharma/auth/models"
)

func ExtractHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("userId")

		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "UserId header is required",
			})
		}
		c.Locals("userId", userID)
		fmt.Println("the userId from the middleware: ", userID)
		return c.Next()
	}
}

func ExtractRoleIdMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleId := c.Get("roleId")
		if roleId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "roleId header is required",
			})
		}
		c.Locals("roleId", roleId)
		return c.Next()
	}
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := VerifyJWT(tokenStr)
		if err != nil {
			fmt.Println("the error while verifying the token: ", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		c.Locals("authClaims", claims)
		return c.Next()
	}
}

func TenantMiddleWare() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("Inside the middleware")
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		Newtoken, err := repo.NewTokenRepository(db.DB)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ServiceResponse{
				Code:    500,
				Message: "error while connecting to db repositry",
			})
		}
		log.Println(" the token string :", tokenStr)
		resp, tenant_id, err := Newtoken.VerifyToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("error while verifying token: %v", err),
			})
		}
		if resp {
			c.Locals("token", tokenStr)
			c.Locals("tenant_id", tenant_id)
		}
		return c.Next()
	}
}

func GetTenantFromToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query("application_key")
		if token == "" {
			return &models.ServiceResponse{
				Code:    400,
				Message: "A valid tenant-generated token is required for every request. Please provide the application_key parameter, or register as a tenant if you are the owner.",
			}
		}
		Newtoken, err := repo.NewTokenRepository(db.DB)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ServiceResponse{
				Code:    500,
				Message: "error while connecting to db repositry",
			})
		}
		resp, tenant_id, err := Newtoken.VerifyToken(token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
		if resp {
			c.Locals("tenant_id", tenant_id)
		}
		return nil
	}
}

func VerifyRoleRouteMapping(roleId string, route string, method string) (bool, error) {
	roleRoute, err := repo.NewRouteRoleRepository(db.DB)
	if err != nil {
		return false, fmt.Errorf("error while initializing the roleroute repository: %w", err)
	}

	routes, err := roleRoute.GetRoleRouteMapping(roleId)
	if err != nil {
		return false, fmt.Errorf("error while getting the role route mapping: %w", err)
	}

	roleData, err := model.ConvertDBData(routes.Permissions)
	if err != nil {
		return false, fmt.Errorf("%s", "error while fetchig the role permissions: "+err.Error())
	}

	permi := model.ClassifyPermissionsByMethod(roleData.Permissions)

	per, flag := model.FindMethodWithPatterns(method, route, permi)

	if per == nil {
		return false, nil
	}

	fmt.Println("the flag: ", flag)

	return flag, nil
}
