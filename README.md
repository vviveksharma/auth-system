# GuardRail

JWT auth middleware for Fiber. That's it.

Ok so I kept copy-pasting the same auth code into every project and finally just extracted it into this library. Started as a full auth API but that was overkill - most of the time you just need something that handles tokens and user sessions without being in your way.

(also was getting annoying to maintain lol)

## Features

- JWT tokens (access + refresh)
- role-based route protection
- Multi-tenant support (optional)
- Redis caching so you're not hitting postgres constantly  
- argon2 password hashing
- reasonable defaults, configure if you want

## Installation

```bash
go get github.com/vviveksharma/auth
```

## Quick Start

Basically 5 lines of code:

```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/vviveksharma/auth"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    db.AutoMigrate(&guardrail.User{})

    gr, _ := guardrail.New(guardrail.Config{
        DB:        db,
        JWTSecret: "your-secret-key",
    })

    app := fiber.New()
    app.Get("/protected", gr.Protect(), func(c *fiber.Ctx) error {
        userID, _ := guardrail.GetUserID(c)
        return c.JSON(fiber.Map{"user_id": userID})
    })

    log.Fatal(app.Listen(":3000"))
}
```

Done. Route is protected.

## Full Example

Most apps need signup/login so here's that:

```go
package main

import (
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/vviveksharma/auth"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    db.AutoMigrate(&guardrail.User{})

    gr, err := guardrail.New(guardrail.Config{
        DB:        db,
        JWTSecret: os.Getenv("JWT_SECRET"),
        EnableRBAC: true,
    })
    if err != nil {
        log.Fatal(err)
    }

    authService := gr.NewAuthService()
    app := fiber.New()

    // Registration endpoint
    app.Post("/register", func(c *fiber.Ctx) error {
        var req guardrail.RegisterRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
        }
        response, err := authService.Register(req)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(response)
    })

    // Login endpoint
    app.Post("/login", func(c *fiber.Ctx) error {
        var req guardrail.LoginRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
        }
        response, err := authService.Login(req)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(response)
    })

    // Authenticated user info
    app.Get("/profile", gr.Protect(), func(c *fiber.Ctx) error {
        userID, _ := guardrail.GetUserID(c)
        role, _ := guardrail.GetRole(c)
        return c.JSON(fiber.Map{
            "user_id": userID,
            "role":    role,
        })
    })

    // Admin-only route
    app.Get("/admin", gr.ProtectWithRole("admin"), func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "Admin access granted"})
    })

    log.Fatal(app.Listen(":3000"))
}
```

## Config

### Basic

Just need DB + secret:

```go
gr, err := guardrail.New(guardrail.Config{
    DB:        db,
    JWTSecret: "secret-key",
})
```

### Advanced

Tweak timeouts, add redis, whatever:

```go
gr, err := guardrail.New(guardrail.Config{
    DB:                 db,
    JWTSecret:          "your-secret-key",
    AccessTokenExpiry:  15 * time.Minute,       // Defaults to 15 minutes
    RefreshTokenExpiry: 7 * 24 * time.Hour,     // Defaults to 7 days
    RedisClient:        redisClient,            // Optional, for caching
    EnableRBAC:         true,                   // Defaults to true
    EnableMultiTenant:  false,                  // Defaults to false
    ErrorMessages: guardrail.ErrorMessages{
        Unauthorized: "Custom unauthorized message",
        Forbidden:    "Custom forbidden message",
    },
})
```

## API

### Middleware stuff

#### `gr.Protect()`
Basic JWT check.

```go
app.Get("/protected", gr.Protect(), handler)
```

#### `gr.ProtectWithRole(roles...)`
Check for specific roles.

```go
app.Get("/admin", gr.ProtectWithRole("admin"), handler)
app.Get("/moderator", gr.ProtectWithRole("admin", "moderator"), handler)
```

#### `gr.ApplicationKeyMiddleware()`
Multi-tenant stuff with API keys.

```go
app.Use(gr.ApplicationKeyMiddleware())
```

### Auth Service

#### `authService.Register(req)`
Signup.

```go
response, err := authService.Register(guardrail.RegisterRequest{
    Email:     "user@example.com",
    Password:  "secure-password",
    FirstName: "John",
    LastName:  "Doe",
    Role:      "user", // defaults to "user"
})
```

#### `authService.Login(req)`
Login, get tokens back.

```go
response, err := authService.Login(guardrail.LoginRequest{
    Email:    "user@example.com",
    Password: "secure-password",
})
```

#### `authService.RefreshToken(refreshToken)`
Get new access token when it expires.

```go
response, err := authService.RefreshToken(refreshToken)
```

#### `authService.Logout(token)`
Blacklist token (needs redis).

```go
err := authService.Logout(accessToken)
```

### Helpers

Pull user data from context:

```go
userID, ok := guardrail.GetUserID(c)
role, ok := guardrail.GetRole(c)
tenantID, ok := guardrail.GetTenantID(c)
claims, ok := guardrail.GetClaims(c)
```

## Database

Needs a users table, something like:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    salt TEXT NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user',
    tenant_id UUID,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

or just let GORM handle it:

```go
db.AutoMigrate(&guardrail.User{})
```

## Multi-tenant

If you need tenant isolation:

```go
gr, _ := guardrail.New(guardrail.Config{
    DB:                db,
    JWTSecret:         "secret",
    EnableMultiTenant: true,
})

app.Use(gr.ApplicationKeyMiddleware())

authService.Register(guardrail.RegisterRequest{
    Email:    "user@example.com",
    Password: "password",
    TenantID: "tenant-uuid",
})
```

## Performance

things that helped:

**Redis caching** - like 10-100x faster for token validation vs hitting postgres

**Short access tokens** - 15min works, refresh can be longer (week or whatever)

**Database indexes** - email, tenant_id, is_active. don't skip this

**Connection pooling** - gorm defaults sometimes aren't enough

YMMV but this is what worked for my usecase

## Security

basics:

- Long random JWT secret (32+ chars) - use `openssl rand -base64 32`
- env variables for secrets, never hardcode
- HTTPS in production obviously
- rate limit login endpoints or you'll get brute forced
- token revocation needs redis
- 15min/7day for access/refresh tokens works pretty well

oh and rotate your secrets periodically

## Usage examples

### Register
```bash
curl -X POST http://localhost:3000/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123"
  }'
```

### Access Protected Route
```bash
curl http://localhost:3000/profile \
  -H "Authorization: Bearer <your-access-token>"
```

## Contributing

PRs welcome. open an issue first for big changes tho

## License

MIT

---

originally built this for a side project, kept reusing it, finally made it a proper library. took some refactoring to clean up the globals but whatever, works now.

had plans to make this a paid SaaS thing but honestly just gonna keep it free and open source. plans changed, priorities shifted, you know how it goes.

problems? open an issue or dm me

[github.com/vviveksharma](https://github.com/vviveksharma)