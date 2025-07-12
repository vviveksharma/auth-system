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


##### Currently adding the insecure cockrochDB will be makeing it secure + TLS 
#### Add Response Structure to every error response



```
Category	Endpoint	                Method	Description	Who Can Access?
Auth	
            /auth/register	            POST	Register new user (auto-assign user role)	guest
            /auth/login	                POST	Generate JWT token	guest
            /auth/refresh	            POST	Refresh expired JWT	user, moderator, admin 
User	
            /users/me	                GET	    Get current user’s profile	user+ roles
            /users/me	                PUT	    Update own profile	user 
            /users/:id	                GET	    Get user by ID (admin view)	admin, moderator
            /users/:id/roles	        PUT	    Assign roles to a user (e.g., make premium)	admin
Roles	
            /roles	                    GET	    List all roles (e.g., user, admin)	admin
            /roles	                    POST	Create custom role (e.g., support_agent)	admin
            /roles/:role/permissions	PUT	    Update permissions for a role	admin
            /roles/verify               GET     To check if the role matches with the roleId
Resources	
            /posts	                    GET	    List public posts	guest
            /posts	                    POST	Create a post	user+ roles
            /posts/:id	                DELETE	Delete a post	owner, moderator, admin
Admin	
            /admin/logs	                GET	    View system access logs	admin
            /admin/stats	            GET	    Get RBAC usage statistics	admin
```

### Feature in the Development
- Add refresh xpired token in the cache ( adding the logic for cache )
- Create a features docs so that people could see and use that
- Create a swagger docs
- Route and role tables
- Add the token verification middleware