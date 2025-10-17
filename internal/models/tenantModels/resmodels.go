package tenantmodels

import "github.com/google/uuid"

type TenantListRoleResponseBody struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"string"`
	DisplayName string    `json:"display_name"`
	RoleType    string    `json:"role_type"`
	Status      bool      `json:"status"`
}
