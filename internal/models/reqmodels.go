package models

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUserDetailsRequest struct {
	Id string `json:"id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type AssignRoleRequest struct {
	Role string `json:"role"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type VerifyRoleRequest struct {
	RoleName string `json:"role_name"`
	RoleId   string `json:"role_id"`
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

type CreateCustomRole struct {
	RoleName string   `json:"role"`
	Routes   []string `json:"routes"`
}

type UpdateRolePermissions struct {
	RoleName          string   `json:"role"`
	AddPermisions     []string `json:"add_permissions"`
	RemovePermissions []string `json:"remove_permissions"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}
