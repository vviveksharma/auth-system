``` 
auth-system/
│
├── cmd/                    # Entry point (main.go)
│   └── server/             # Initialize routes, middleware, server
│       └── main.go
│
├── config/                 # Configuration loader (env, DB, Redis)
│   └── config.go
│
├── internal/               # Core logic broken into modules
│   ├── models/             # DB Models (User, Token, etc.)
│   ├── controllers/        # Route Handlers (Register, Login)
│   ├── services/           # Business logic (auth, JWT, sessions)
│   ├── middlewares/        # JWT validation, Role-based guard, Rate limiter
│   └── utils/              # Utilities (hashing, token generation)
│
├── db/                     # DB-related setup
│   ├── migrations/         # SQL files (or use golang-migrate)
│   └── postgres.go         # Connection setup
│
├── redis/                  # Redis client setup
│   └── redis.go
│
├── routes/                 # API route definitions
│   └── routes.go
│
├── .env                    # Environment variables
├── go.mod
├── go.sum
├── Dockerfile              # For containerizing the app
├── docker-compose.yml      # Local setup: app + Redis + PostgreSQL
├── README.md
└── Makefile (optional)     # Build, run, test shortcuts
```


##### Currently adding the insecure cockrochDB need to make it secure + TLS 