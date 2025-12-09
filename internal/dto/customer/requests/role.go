package requests

import "github.com/vviveksharma/auth/internal/dto/customer/shared"

type CreateCustomRole struct {
	Name        string              `json:"name"`
	DisplayName string              `json:"display_name"`
	Description string              `json:"description"`
	Permissions []shared.Permission `json:"Permissions"`
}

type UpdateRolePermissions struct {
	RoleName          string              `json:"role"`
	AddPermisions     []shared.Permission `json:"add_permissions"`
	RemovePermissions []shared.Permission `json:"remove_permissions"`
}
