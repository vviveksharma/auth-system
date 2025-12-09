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
	Name        string       `json:"name"`
	DisplayName string       `json:"display_name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"Permissions"`
}

type UpdateRolePermissions struct {
	RoleName          string       `json:"role"`
	AddPermisions     []Permission `json:"add_permissions"`
	RemovePermissions []Permission `json:"remove_permissions"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type CreateTokenRequest struct {
	Name     string `json:"name"`
	ExpiryAt string `json:"expiry_at"`
}

type ResetTenantPasswordRequest struct {
	Email string `json:"email"`
}

type SetTenantPasswordRequest struct {
	Email              string `json:"email"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type UserVerifyOTPRequest struct {
	OTP             string `json:"otp"`
	Email           string `json:"email"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type CreateMessageRequest struct {
	Email         string `json:"email"`
	RequestedRole string `json:"requested_role"`
}

type ListMessageRequest struct {
	Email         string `json:"email"`
}