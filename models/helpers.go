package models

import "fmt"

type ServiceResponse struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *ServiceResponse) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Data: %+v", e.Code, e.Message, e.Data)
}

type StatusUnprocessableEntityResponse struct {
	Code    int    `json:"code" example:"422"`
	Message string `json:"message" example:"Invalid request body format. Please check your JSON syntax and field types."`
}

type UnauthorizedResponse struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"Unauthorized access. Please provide valid authentication credentials."`
}

type BadRequestResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad request. Missing or invalid required fields."`
}

type ConflictResponse struct {
	Code    int    `json:"code" example:"409"`
	Message string `json:"message" example:"Conflict. Resource already exists or operation conflicts with current state."`
}

type InternalServerErrorResponse struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Internal server error occurred."`
}
