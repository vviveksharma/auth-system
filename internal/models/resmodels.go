package models

type UserResponse struct {
	Message string `json:"message"`
}

type UserDetailsResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginResponse struct {
	Jwt string `json:"jwt"`
}
