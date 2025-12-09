package tenantmodels

import (
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/internal/models"
)

type TenantListRoleResponseBody struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"string"`
	DisplayName string    `json:"display_name"`
	RoleType    string    `json:"role_type"`
	Status      bool      `json:"status"`
}

type TenantGetPermissionsResponseBody struct {
	Id          uuid.UUID           `json:"id"`
	RoleInfo    models.RoleInfo     `json:"role_info"`
	Permissions []models.Permission `json:"permissions"`
}

type TenantDisableRoleResponsBody struct {
	Message string `json:"message"`
}

type TenantEnableRoleResponseBody struct {
	Message string `json:"message"`
}

type TenantEditPermissionRoleResponseBody struct {
	Message string `json:"message"`
}

type TenantAddRoleResponseBody struct {
	Message string `json:"message"`
}

type TenantListUserResponseBody struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Roles     []string  `json:"roles"`
}

type TenantEnableUserResponseBody struct {
	Message string `json:"message"`
}

type TenantDeleteUserResponseBody struct {
	Message string `json:"message"`
}

type TenantDeleteRoleResponseBody struct {
	Message string `json:"message"`
}

type TenantMessageResponseBoy struct {
	MessageId     uuid.UUID `json:"message_id"`
	UserEmail     string    `json:"user_email"`
	CurrentRole   string    `json:"current_role"`
	RequestedRole string    `json:"requested_role"`
	Status        string    `json:"status"`
	RequestAt     string    `json:"request_at"`
}

type TenantApproveMessageResponseBody struct {
	Message string `json:"message"`
}

type TenantRejectMessageResponseBody struct {
	Message string `json:"message"`
}
