package middlewares

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/cache"
	"github.com/vviveksharma/auth/db"
	model "github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
)

func TenantMiddleWare() fiber.Handler {
    return func(c *fiber.Ctx) error {
        log.Println("Inside the TenantMiddleWare")
        
        // Extract Authorization header
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":       true,
                "message":     "Missing authorization token",
                "status_code": fiber.StatusUnauthorized,
            })
        }

        // Extract token from "Bearer <token>" format
        tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))  // ✅ Added space
        if tokenStr == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":       true,
                "message":     "Invalid authorization header format. Expected: Bearer <token>",
                "status_code": fiber.StatusUnauthorized,
            })
        }

        // Check cache first
        cacheKey := "token:" + tokenStr
        var tenant_id string
        err := cache.Get(cacheKey, &tenant_id)
        if err == nil && tenant_id != "" {  // ✅ Check tenant_id is not empty
            log.Printf("✅ Cache hit for token: %s, tenant_id: %s", tokenStr, tenant_id)
            c.Locals("token", tokenStr)
            c.Locals("tenant_id", tenant_id)
            return c.Next()
        }

        log.Println("⚠️ Cache miss, verifying token from database...")

        // Verify token from database
        Newtoken, err := repo.NewTokenRepository(db.DB)
        if err != nil {
            log.Printf("❌ Error creating token repository: %v", err)
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error":       true,
                "message":     "error while connecting to db repository",
                "status_code": 500,
            })
        }

        log.Println("Verifying token:", tokenStr)
        var resp bool
        resp, tenant_id, err = Newtoken.VerifyToken(tokenStr) 
        if err != nil {
            log.Printf("❌ Token verification failed: %v", err)
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":       true,
                "message":     fmt.Sprintf("error while verifying token: %v", err),
                "status_code": fiber.StatusUnauthorized,
            })
        }

        // Check if token is valid
        if !resp || tenant_id == "" { 
            log.Println("❌ Invalid token or missing tenant_id")
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":       true,
                "message":     "Invalid or expired token",
                "status_code": fiber.StatusUnauthorized,
            })
        }

        // Token is valid - cache it and continue
        log.Printf("✅ Token verified successfully for tenant: %s", tenant_id)
        cache.Set(cacheKey, tenant_id, 24*time.Hour)
        c.Locals("token", tokenStr)
        c.Locals("tenant_id", tenant_id)
        
        return c.Next() 
    }
}

func VerifyRoleRouteMapping(roleId string, route string, method string) (bool, error) {
	roleRoute, err := repo.NewRouteRoleRepository(db.DB)
	if err != nil {
		return false, fmt.Errorf("error while initializing the roleroute repository: %w", err)
	}

	cacheKey := fmt.Sprintf("role_access:%s:%s:%s", roleId, route, method)
	var hasAccess bool
	err = cache.Get(cacheKey, &hasAccess)
	if err == nil {
		log.Println("CACHE HIT")
		return hasAccess, nil
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
		hasAccess = false
		return false, nil
	}

	hasAccess = flag

	cache.Set(cacheKey, hasAccess, 30*time.Minute)

	fmt.Println("the flag: ", flag)

	return flag, nil
}
