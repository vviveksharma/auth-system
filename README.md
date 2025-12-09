The refactoring for the new version of the software altogether
backend/
│
├── cmd/
│   └── server/
│       └── main.go                      # Application entrypoint
│
├── internal/
│   │
│   ├── domain/                          # Core domain layer
│   │   │
│   │   ├── entities/                    # Database models (shared)
│   │   │   ├── user.go
│   │   │   ├── tenant.go
│   │   │   ├── role.go
│   │   │   ├── route_role.go
│   │   │   ├── token.go
│   │   │   ├── message.go
│   │   │   ├── login.go
│   │   │   └── reset_token.go
│   │   │
│   │   ├── dto/                         # Data Transfer Objects
│   │   │   │
│   │   │   ├── customer/                # Customer API contracts
│   │   │   │   ├── requests/
│   │   │   │   │   ├── auth.go         # RegisterUserRequest, LoginRequest
│   │   │   │   │   ├── user.go         # UpdateProfileRequest
│   │   │   │   │   ├── role.go         # RequestRoleRequest
│   │   │   │   │   └── message.go      # CreateMessageRequest
│   │   │   │   │
│   │   │   │   └── responses/
│   │   │   │       ├── auth.go         # LoginResponse, RegisterResponse
│   │   │   │       ├── user.go         # UserResponse, UserListResponse
│   │   │   │       ├── role.go         # RoleResponse, RoleListResponse
│   │   │   │       └── message.go      # MessageResponse, MessageStatusResponse
│   │   │   │
│   │   │   └── tenant/                  # Tenant/Admin API contracts
│   │   │       ├── requests/
│   │   │       │   ├── tenant.go       # CreateTenantRequest, LoginTenantRequest
│   │   │       │   ├── user_admin.go   # DisableUserRequest, AssignRoleRequest
│   │   │       │   ├── role_admin.go   # CreateRoleRequest, UpdatePermissionsRequest
│   │   │       │   ├── token.go        # CreateTokenRequest, RevokeTokenRequest
│   │   │       │   ├── message.go      # ApproveMessageRequest, RejectMessageRequest
│   │   │       │   └── dashboard.go    # DashboardFiltersRequest
│   │   │       │
│   │   │       └── responses/
│   │   │           ├── tenant.go       # TenantResponse, TenantDetailsResponse
│   │   │           ├── user_admin.go   # UserAdminResponse, UserListAdminResponse
│   │   │           ├── role_admin.go   # RoleAdminResponse, PermissionsResponse
│   │   │           ├── token.go        # TokenResponse, TokenListResponse
│   │   │           ├── message.go      # MessageAdminResponse
│   │   │           └── dashboard.go    # DashboardStatsResponse
│   │   │
│   │   └── errors/                      # Custom error types
│   │       ├── errors.go               # AppError, ValidationError, etc.
│   │       └── codes.go                # Error codes
│   │
│   ├── repository/                      # Data access layer
│   │   │
│   │   ├── customer/                    # Customer-context repositories
│   │   │   ├── user_repository.go
│   │   │   ├── role_repository.go
│   │   │   ├── message_repository.go
│   │   │   ├── login_repository.go
│   │   │   ├── token_repository.go
│   │   │   └── interfaces.go
│   │   │
│   │   ├── tenant/                      # Tenant-context repositories
│   │   │   ├── tenant_repository.go
│   │   │   ├── user_repository.go
│   │   │   ├── role_repository.go
│   │   │   ├── message_repository.go
│   │   │   ├── token_repository.go
│   │   │   ├── tenant_login_repository.go
│   │   │   ├── dashboard_repository.go
│   │   │   └── interfaces.go
│   │   │
│   │   └── shared/                      # Shared repository utilities
│   │       ├── base_repository.go
│   │       ├── reset_token_repository.go
│   │       └── route_role_repository.go
│   │
│   ├── service/                         # Business logic layer
│   │   │
│   │   ├── customer/                    # Customer business logic
│   │   │   ├── auth_service.go
│   │   │   ├── user_service.go
│   │   │   ├── role_service.go
│   │   │   ├── message_service.go
│   │   │   └── interfaces.go
│   │   │
│   │   ├── tenant/                      # Tenant business logic
│   │   │   ├── tenant_service.go
│   │   │   ├── user_admin_service.go
│   │   │   ├── role_admin_service.go
│   │   │   ├── token_service.go
│   │   │   ├── message_admin_service.go
│   │   │   ├── dashboard_service.go
│   │   │   └── interfaces.go
│   │   │
│   │   └── shared/                      # Shared services
│   │       ├── mail_service.go
│   │       ├── cache_service.go
│   │       └── password_service.go
│   │
│   ├── api/                             # HTTP API layer
│   │   │
│   │   ├── customer/                    # Customer-facing API (Port 8080)
│   │   │   ├── handlers/
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── user_handler.go
│   │   │   │   ├── role_handler.go
│   │   │   │   ├── message_handler.go
│   │   │   │   └── common.go          # Shared response helpers
│   │   │   │
│   │   │   ├── middleware/
│   │   │   │   ├── app_key.go
│   │   │   │   ├── jwt_auth.go
│   │   │   │   ├── permissions.go
│   │   │   │   ├── ratelimit.go
│   │   │   │   └── chain.go           # Middleware chains
│   │   │   │
│   │   │   └── routes.go              # Customer API routes
│   │   │
│   │   └── tenant/                      # Tenant-facing API (Port 8081)
│   │       ├── handlers/
│   │       │   ├── tenant_handler.go
│   │       │   ├── user_admin_handler.go
│   │       │   ├── role_admin_handler.go
│   │       │   ├── token_handler.go
│   │       │   ├── message_admin_handler.go
│   │       │   ├── dashboard_handler.go
│   │       │   └── common.go          # Shared response helpers
│   │       │
│   │       ├── middleware/
│   │       │   ├── tenant_auth.go
│   │       │   ├── ratelimit.go
│   │       │   └── chain.go
│   │       │
│   │       └── routes.go              # Tenant API routes
│   │
│   ├── bootstrap/                       # Application initialization
│   │   ├── app.go                      # Main app bootstrap
│   │   ├── database.go                 # Database setup & migrations
│   │   ├── cache.go                    # Redis setup
│   │   ├── queue.go                    # RabbitMQ setup
│   │   ├── server.go                   # HTTP servers setup
│   │   └── dependencies.go             # Dependency injection
│   │
│   └── pkg/                            # Internal shared packages
│       │
│       ├── cache/                      # Cache utilities
│       │   ├── cache.go
│       │   └── redis.go
│       │
│       ├── queue/                      # Message queue utilities
│       │   ├── queue.go
│       │   ├── consumer.go
│       │   └── producer.go
│       │
│       ├── mailer/                     # Email service
│       │   ├── mailer.go
│       │   └── smtp.go
│       │
│       ├── logger/                     # Logging utilities
│       │   └── logger.go
│       │
│       ├── validator/                  # Validation utilities
│       │   └── validator.go
│       │
│       ├── crypto/                     # Cryptography utilities
│       │   ├── hash.go
│       │   └── jwt.go
│       │
│       └── pagination/                 # Pagination utilities
│           ├── pagination.go
│           └── response.go
│
├── config/                             # Configuration files
│   ├── permissions/                    # Role permission JSON files
│   │   ├── admin.json
│   │   ├── user.json
│   │   ├── moderator.json
│   │   └── guest.json
│   │
│   └── database/                       # Database migrations
│       └── migrations/
│           ├── 001_create_users_table.sql
│           ├── 002_create_roles_table.sql
│           └── ...
│
├── scripts/                            # Utility scripts
│   ├── migrate.sh                      # Run database migrations
│   ├── seed.sh                         # Seed initial data
│   └── test.sh                         # Run tests
│
├── test/                               # Tests
│   ├── integration/
│   │   ├── customer_api_test.go
│   │   ├── tenant_api_test.go
│   │   └── fixtures/
│   │       └── test_data.sql
│   │
│   └── unit/
│       ├── service/
│       │   └── user_service_test.go
│       │
│       └── repository/
│           └── user_repository_test.go
│
├── docs/                               # API documentation
│   ├── swagger/
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   │
│   ├── API_CUSTOMER.md                # Customer API docs
│   └── API_TENANT.md                  # Tenant API docs
│
├── .env.example                        # Environment variables template
├── .gitignore
├── Dockerfile                          # Multi-stage build
├── Makefile                            # Build automation
├── go.mod
├── go.sum
└── README.md