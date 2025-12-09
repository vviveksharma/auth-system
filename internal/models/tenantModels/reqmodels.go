package tenantmodels

import (
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/internal/models"
)

type TeanantEditPermissionRequestBody struct {
	Id                    uuid.UUID           `json:"id"`
	UpdateRoleDetails     bool                `json:"update_role_details"`
	RoleInfo              models.RoleInfo     `json:"role_info"`
	UpdateRolePermissions bool                `json:"update_role_permissions"`
	Permissions           []models.Permission `json:"permissions"`
}

type TenantAddRoleRequestBody struct {
	Name        string              `json:"name"`
	DisplayName string              `json:"display_name"`
	Description string              `json:"description"`
	Permissions []models.Permission `json:"permissions"`
}
