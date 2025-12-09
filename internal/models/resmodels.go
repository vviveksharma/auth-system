package models

import (
	"time"

	"github.com/google/uuid"
)

// **NEW: Permission-related models for seeding**
type RoutePermission struct {
	Route       string                 `json:"route"`
	Methods     []string               `json:"methods"`
	Params      map[string]string      `json:"params,omitempty"`
	Description string                 `json:"description,omitempty"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
}

type PermissionSet struct {
	RoleId      uuid.UUID         `json:"role_id"`
	Permissions []RoutePermission `json:"permissions"`
	Version     int               `json:"version"`
	UpdatedAt   string            `json:"updated_at"`
}

// **NEW: Role management models**
type ListRolesResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Status   bool      `json:"status"`
	IsSystem bool      `json:"is_system"`
}

type RoleDetailsResponse struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Status      bool              `json:"status"`
	IsSystem    bool              `json:"is_system"`
	Permissions []RoutePermission `json:"permissions"`
	UserCount   int64             `json:"user_count"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type UpdateRolePermissionsRequest struct {
	RoleId      uuid.UUID         `json:"role_id" validate:"required"`
	Permissions []RoutePermission `json:"permissions" validate:"required"`
}

type ListRolePermissionsResponse struct {
	RoleId      uuid.UUID         `json:"role_id"`
	RoleName    string            `json:"role_name"`
	Permissions []RoutePermission `json:"permissions"`
	Version     int               `json:"version"`
	UpdatedAt   string            `json:"updated_at"`
}

type UserResponse struct {
	Message string `json:"message"`
}

type UserDetailsResponse struct {
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Role  []string `json:"role"`
}

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

type UpdateUserResponse struct {
	Message string `json:"message"`
}

type GetUserByIdResponse struct {
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Role  []string `json:"role"`
}

type AssignRoleResponse struct {
	Message string `json:"message"`
}

type ListAllRolesResponse struct {
	Name     string    `json:"name"`
	RoleId   uuid.UUID `json:"role_id"`
	TenantId uuid.UUID `gorm:"type:uuid;not null"`
	RoleType string    `json:"role_type"`
	Status   bool      `json:"status"`
	Routes   []string  `json:"routes"`
}

type UserLoginResponse struct {
	JWT string `json:"jwt"`
}

type VerifyRoleResponse struct {
	Message bool `json:"message"`
}

type CreateTenantResponse struct {
	Message string `json:"message"`
}

type LoginTenantResponse struct {
	Token string `json:"token"`
}

type RevokeTokenResponse struct {
	Message string `json:"message"`
}

type CreateCustomRoleResponse struct {
	Message string `json:"message"`
}

type UpdateRolePermissionsResponse struct {
	Message string `json:"message"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type ListTokensResponse struct {
	Name      string    `json:"name"`
	TokenId   uuid.UUID `json:"token_id"`
	CreateAt  time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expiry_at"`
	Status    bool      `json:"status"`
}

type CreateTokenResponse struct {
	Message string `json:"message"`
}

type ResetPasswordTenantResponse struct {
	Message string `json:"message"`
}

type SetTenantPasswordResponse struct {
	Message string `json:"message"`
}

type ResetUserPasswordResponse struct {
	Message string `json:"message"`
}

type UserVerifyOTPResponse struct {
	Message string `json:"message"`
}

type ListUserTenant struct {
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Role         []string `json:"role"`
	LogginStatus bool     `json:"log_status"`
	CreatedAt    string   `json:"created_at"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type DeleteRoleResponse struct {
	Message string `json:"message"`
}

type LogoutUserResponse struct {
	Message string `json:"struct"`
}

type EnableRoleResponse struct {
	Message string `json:"struct"`
}

type DisableRoleResponse struct {
	Message string `json:"message"`
}

type GetRouteDetailsResponse struct {
	Routes      interface{} `json:"routes"`
	RoutesJSON  string      `json:"routes_json,omitempty"`
	RoleInfo    RoleInfo    `json:"role_info,omitempty"`
	ProcessedAt time.Time   `json:"processed_at"`
}

type ListUsersResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
}

type EnableUserResponse struct {
	Message string `json:"message"`
}

type DisableUserResponse struct {
	Message string `json:"message"`
}

type GetTenantDetails struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organisation string `json:"organisation"`
}

type GetRoleDetailsUser struct {
	UserId string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
}

type DeleteTenantResponse struct {
	Message string `json:"message"`
}

type DashboardTenantResponse struct {
	UsersCount int `json:"user_count"`
	RoleCount  int `json:"role_count"`
	TokenCount int `json:"token_count"`
}

type GetListTokenWithStatus struct {
	Name      string    `json:"name"`
	TokenId   uuid.UUID `json:"token_id"`
	CreateAt  time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expiry_at"`
	Status    bool      `json:"status"`
}

type CreateMessageResponse struct {
	Message string `json:"message"`
}

type GetMessageStatusResponse struct {
	Status string `json:"string"`
}

type ListMessageStatusResponse struct {
	MessageId     string `json:"message_id"`
	Status        string `json:"status"`
	RequestedRole string `json:"requested_role"`
}