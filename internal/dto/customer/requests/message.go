package requests

type CreateMessageRequest struct {
	Email         string `json:"email"`
	RequestedRole string `json:"requested_role"`
}
