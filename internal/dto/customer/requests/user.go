package requests


type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUserDetailsRequest struct {
	Id string `json:"id"`
}

type AssignRoleRequest struct {
	Role string `json:"role"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type UserVerifyOTPRequest struct {
	OTP             string `json:"otp"`
	Email           string `json:"email"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}