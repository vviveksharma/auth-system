package requests

import (
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/internal/dto/customer/shared"
)

type TeanantEditPermissionRequestBody struct {
	Id                    uuid.UUID           `json:"id"`
	UpdateRoleDetails     bool                `json:"update_role_details"`
	RoleInfo              shared.RoleInfo     `json:"role_info"`
	UpdateRolePermissions bool                `json:"update_role_permissions"`
	Permissions           []shared.Permission `json:"permissions"`
}

type TenantAddRoleRequestBody struct {
	Name        string              `json:"name"`
	DisplayName string              `json:"display_name"`
	Description string              `json:"description"`
	Permissions []shared.Permission `json:"permissions"`
}

type ResetTenantPasswordRequest struct {
	Email string `json:"email"`
}

type CreateTenantRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Campany  string `json:"campany"`
	Password string `json:"password"`
}

type LoginTenantRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}