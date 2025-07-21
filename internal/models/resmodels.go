package models

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
