package models

import (
	"time"

	"github.com/google/uuid"
)

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
	Name string `json:"name"`
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
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	LogginStatus bool   `json:"log_status"`
	CreatedAt    string `json:"created_at"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type DeleteRoleResponse struct {
	Message string `json:"message"`
}