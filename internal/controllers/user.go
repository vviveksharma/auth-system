package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/internal/models"
	responsemodels "github.com/vviveksharma/auth/models"
)

// GetUserDetails retrieves details of the currently authenticated user.
//
// @Summary Get Authenticated User Details
// @Description Returns the details (name, email, roles) of the user currently authenticated via JWT token. Extracted from token claims. Returns 404 if user not found.
// @Tags User
// @Produce json
// @Success 200 {object} responsemodels.ServiceResponse{data=models.UserDetailsResponse} "User details successfully retrieved with name, email, and roles array"
// @Failure 401 {object} responsemodels.UnauthorizedResponse "Unauthorized, invalid or missing JWT authentication"
// @Failure 404 {object} responsemodels.ServiceResponse "User not found in the system"
// @Failure 500 {object} responsemodels.InternalServerErrorResponse "Internal server error"
// @Router /user/me [get]
// @Security ApiKeyAuth
func (h *Handler) GetUserDetails(ctx *fiber.Ctx) error {
	req := &models.GetUserDetailsRequest{}
	userId := ctx.Locals("user_id").(string)
	fmt.Println("the userid: ", userId)
	req.Id = userId
	resp, err := h.UserService.GetUserDetails(ctx.Context(), req)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "",
		Data:    resp,
	})
}

// UpdateUserDetails updates the details of the currently authenticated user.
//
// @Summary Update Authenticated User Details
// @Description Allows the authenticated user to update their profile information (name, email, or password). At least one field must be provided. User ID extracted from JWT token.
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UpdateUserRequest true "Fields to update (email, name, password). At least one field required. All fields are optional pointers."
// @Success 200 {object} responsemodels.ServiceResponse "User details updated successfully"
// @Failure 400 {object} responsemodels.BadRequestResponse "Bad request, no fields provided for update"
// @Failure 401 {object} responsemodels.UnauthorizedResponse "Unauthorized, invalid or missing JWT authentication"
// @Failure 404 {object} responsemodels.ServiceResponse "User not found"
// @Failure 422 {object} responsemodels.StatusUnprocessableEntityResponse "Unprocessable entity, invalid JSON format"
// @Failure 500 {object} responsemodels.InternalServerErrorResponse "Internal server error"
// @Router /user/me [put]
// @Security ApiKeyAuth
func (h *Handler) UpdateUserDetails(ctx *fiber.Ctx) error {
	req := &models.UpdateUserRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return UnprocessableEntity(ctx)
	}
	if req.Email == nil && req.Name == nil && req.Password == nil {
		return BadRequest(ctx, "At least one field (email, name, or password) must be provided for update.")
	}
	userId := ctx.Locals("user_id").(string)
	fmt.Println("the userid: ", userId)
	resp, err := h.UserService.UpdateUserDetails(ctx.Context(), req, userId)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}

// GetUserByIdDetails retrieves user details by user ID.
//
// @Summary Get User Details by ID
// @Description Fetches the details of a user by their unique user ID. This endpoint is typically used by admins or services that need to look up users.
// @Tags User
// @Produce json
// @Param id path string true "Unique identifier of the user to retrieve"
// @Success 200 {object} responsemodels.ServiceResponse "User details successfully retrieved"
// @Failure 401 {object} responsemodels.UnauthorizedResponse "Unauthorized, invalid or missing authentication"
// @Failure 404 {object} responsemodels.ServiceResponse "User not found"
// @Failure 500 {object} responsemodels.InternalServerErrorResponse "Internal server error"
// @Router /user/{id} [get]
// @Security ApiKeyAuth
func (h *Handler) GetUserByIdDetails(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	resp, err := h.UserService.GetUserById(ctx.Context(), userId)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "",
		Data:    resp,
	})
}

// AssignUserRole assigns a role to a user by user ID.
//
// @Summary Assign Role to User
// @Description Assigns a specific role (role name as string) to a user identified by their user ID. Role must exist in the system. Only users with admin privileges can perform this action.
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Unique identifier (UUID) of the user to assign a role"
// @Param request body models.AssignRoleRequest true "Role assignment details with 'role' field (role name string)"
// @Success 200 {object} responsemodels.ServiceResponse "Role assigned successfully"
// @Failure 400 {object} responsemodels.BadRequestResponse "Bad request, 'role' field is required and cannot be empty"
// @Failure 401 {object} responsemodels.UnauthorizedResponse "Unauthorized, invalid or missing authentication"
// @Failure 404 {object} responsemodels.ServiceResponse "User not found"
// @Failure 422 {object} responsemodels.StatusUnprocessableEntityResponse "Unprocessable entity, invalid JSON format"
// @Failure 500 {object} responsemodels.InternalServerErrorResponse "Internal server error"
// @Router /user/{id}/roles [put]
// @Security ApiKeyAuth
func (h *Handler) AssignUserRole(ctx *fiber.Ctx) error {
	var req *models.AssignRoleRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return UnprocessableEntity(ctx)
	}
	userId := ctx.Params("id")
	if req.Role == "" {
		return BadRequest(ctx, "Invalid request: 'roleName' fields is required and cannot be empty.")
	}
	resp, err := h.UserService.AssignUserRole(ctx.Context(), req, userId)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}

// RegisterUser registers a new user under a tenant.
//
// @Summary Register New User
// @Description Registers a new user in the system under the authenticated tenant. User is assigned 'guest' role by default. Email must be unique within tenant. Password is hashed with Argon2.
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UserRequest true "User registration details: name (required), email (required), password (required)"
// @Success 200 {object} responsemodels.ServiceResponse "User registered successfully with default 'guest' role"
// @Failure 400 {object} responsemodels.BadRequestResponse "Bad request, missing required fields (name, email, or password)"
// @Failure 409 {object} responsemodels.ConflictResponse "Conflict, user with this email already exists in tenant"
// @Failure 422 {object} responsemodels.StatusUnprocessableEntityResponse "Unprocessable entity, invalid JSON format"
// @Failure 500 {object} responsemodels.InternalServerErrorResponse "Internal server error"
// @Router /auth/ [post]
// @Security ApiKeyAuth
func (h *Handler) RegisterUser(ctx *fiber.Ctx) error {
	var req *models.UserRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return UnprocessableEntity(ctx)
	}
	if req.Email == "" || req.Name == "" || req.Password == "" {
		return BadRequest(ctx, "Invalid request: 'email' , 'password' and 'name' fields are required and cannot be empty.")
	}
	resp, err := h.UserService.RegisterUser(ctx.Context(), req)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}

// ResetUserPassword initiates the password reset process for a user.
//
// @Summary Reset User Password
// @Description Initiates password reset by creating a reset token valid for 15 minutes. Returns OTP/token in response (should be sent via email in production). User must exist in the authenticated tenant.
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Password reset request with 'email' field"
// @Success 200 {object} responsemodels.ServiceResponse "Password reset token generated successfully. Response contains OTP in message field."
// @Failure 400 {object} responsemodels.BadRequestResponse "Bad request, email field is required or user not found with provided email"
// @Failure 422 {object} responsemodels.StatusUnprocessableEntityResponse "Unprocessable entity, invalid JSON format"
// @Failure 500 {object} responsemodels.InternalServerErrorResponse "Internal server error during token generation"
// @Router /user/resetpassword [post]
// @Security ApiKeyAuth
func (h *Handler) ResetUserPassword(ctx *fiber.Ctx) error {
	var req *models.ResetPasswordRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return UnprocessableEntity(ctx)
	}
	if req.Email == "" {
		return BadRequest(ctx, "Invalid request: 'email' field is required and cannot be empty.")
	}
	resp, err := h.UserService.ResetPassword(ctx.Context(), req)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}

// SetUserPassword sets a new password for the user after OTP verification.
//
// @Summary Set New User Password
// @Description Sets a new password after verifying OTP token. Requires email, OTP (from reset request), new_password, and confirm_password. Passwords must match. OTP expires after 15 minutes.
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UserVerifyOTPRequest true "OTP verification with fields: email, otp, new_password, confirm_password (all required)"
// @Success 200 {object} responsemodels.ServiceResponse "Password updated successfully"
// @Failure 400 {object} responsemodels.BadRequestResponse "Bad request, missing required fields (email, otp, new_password, or confirm_password)"
// @Failure 404 {object} responsemodels.ServiceResponse "User not found with provided email"
// @Failure 409 {object} responsemodels.ConflictResponse "Conflict, passwords don't match or OTP expired/invalid"
// @Failure 422 {object} responsemodels.StatusUnprocessableEntityResponse "Unprocessable entity, invalid JSON format"
// @Failure 500 {object} responsemodels.InternalServerErrorResponse "Internal server error during password update"
// @Router /user/setpassword [put]
// @Security ApiKeyAuth
func (h *Handler) SetUserPassword(ctx *fiber.Ctx) error {
	var req *models.UserVerifyOTPRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return UnprocessableEntity(ctx)
	}
	if req.Email == "" || req.OTP == "" || req.ConfirmPassword == "" || req.NewPassword == "" {
		return BadRequest(ctx, "Invalid request: 'otp', 'confirm_password', 'new_password' and 'email' fields are required and cannot be empty.")
	}
	if req.NewPassword != req.ConfirmPassword {
		log.Printf("Password mismatch for email: %s", req.Email)
		return ctx.Status(fiber.StatusConflict).JSON(responsemodels.ServiceResponse{
			Code:    fiber.StatusConflict,
			Message: "Password confirmation failed: new password and confirmation do not match. Please ensure both fields are identical.",
		})
	}
	resp, err := h.UserService.SetPassword(ctx.Context(), req)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "The password was successfully updated",
		Data:    resp,
	})
}

// DeleteUser handles the deletion of a user by their ID
// @Summary Delete a user
// @Description Permanently deletes a user from the system using their unique identifier (UUID). User must belong to the authenticated tenant. Cannot be undone.
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID format)" format(uuid)
// @Success 200 {object} responsemodels.ServiceResponse "User successfully deleted with confirmation message"
// @Failure 400 {object} responsemodels.BadRequestResponse "Bad request, 'id' path parameter is missing or empty"
// @Failure 401 {object} responsemodels.UnauthorizedResponse "Unauthorized, invalid or missing authentication"
// @Failure 404 {object} responsemodels.ServiceResponse "User not found with provided ID in the authenticated tenant"
// @Failure 500 {object} responsemodels.InternalServerErrorResponse "Internal server error during deletion"
// @Router /user/{id} [delete]
// @Security ApiKeyAuth
func (h *Handler) DeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return BadRequest(ctx, "Invalid request: 'id' in path parameter is required and cannot be empty.")
	}
	resp, err := h.UserService.DeleteUser(ctx.Context(), uuid.MustParse(id))
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			log.Printf("Unexpected error while deleting user with id %s: %v", id, err)
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "The user was successfully deleted",
		Data:    resp,
	})
}

// ListUsers retrieves a paginated list of all users in the system.
//
// @Summary List All Users
// @Description Fetches a list of all users for the authenticated tenant. Returns user details including id, email, name, created_at (RFC3339), and roles array. Client-side pagination applied (page 1, size 5).
// @Tags User
// @Produce json
// @Success 200 {object} responsemodels.ServiceResponse "Users list successfully retrieved with client-side pagination metadata"
// @Failure 401 {object} responsemodels.UnauthorizedResponse "Unauthorized, invalid or missing authentication"
// @Failure 404 {object} responsemodels.ServiceResponse "No users found for the authenticated tenant"
// @Failure 500 {object} responsemodels.InternalServerErrorResponse "Internal server error during user retrieval"
// @Router /users [get]
// @Security ApiKeyAuth
func (h *Handler) ListUsers(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	pageSize := ctx.QueryInt("page_size", 5)
	status := ctx.Query("roleTypeFlag")
	if status == "" {
		status = "enabled"
	} else if status != "enabled" && status != "disabled" {
		return BadRequest(ctx, "the roleTypeFlag could be only default or user")
	}
	resp, err := h.UserService.ListUsers(ctx.Context(), page, pageSize, status)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			log.Printf("Unexpected error while fetching users : %v", err)
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while fetching users: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "The user have fetch successfully",
		Data:    resp,
	})
}

// EnableUser enables a disabled user account by their ID.
//
// @Summary Enable User Account
// @Description Enables a previously disabled user account, allowing them to access the system again. Returns 400 if user is already enabled. User must belong to the authenticated tenant.
// @Tags User
// @Produce json
// @Param id path string true "Unique identifier (UUID) of the user to enable" format(uuid)
// @Success 200 {object} responsemodels.ServiceResponse "User account enabled successfully with confirmation message"
// @Failure 400 {object} responsemodels.BadRequestResponse "Bad request, 'id' is missing/empty OR user is already enabled"
// @Failure 401 {object} responsemodels.UnauthorizedResponse "Unauthorized, invalid or missing authentication"
// @Failure 404 {object} responsemodels.ServiceResponse "User not found with provided ID in the authenticated tenant"
// @Failure 500 {object} responsemodels.InternalServerErrorResponse "Internal server error during status update"
// @Router /users/{id}/enable [put]
// @Security ApiKeyAuth
func (h *Handler) EnableUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return BadRequest(ctx, "Invalid request: 'id' in path parameter is required and cannot be empty.")
	}
	resp, err := h.UserService.EnableUser(ctx.Context(), uuid.MustParse(id))
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			log.Printf("Unexpected error while fetching users : %v", err)
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while fetching users: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "The user was successfully enabled",
		Data:    resp.Message,
	})
}

// DisableUser disables a user account by their ID.
//
// @Summary Disable User Account
// @Description Disables a user account, preventing them from accessing the system. Returns 400 if user is already disabled. User must belong to the authenticated tenant. Existing sessions remain valid until expiry.
// @Tags User
// @Produce json
// @Param id path string true "Unique identifier (UUID) of the user to disable" format(uuid)
// @Success 200 {object} responsemodels.ServiceResponse "User account disabled successfully with confirmation message"
// @Failure 400 {object} responsemodels.BadRequestResponse "Bad request, 'id' is missing/empty OR user is already disabled"
// @Failure 401 {object} responsemodels.UnauthorizedResponse "Unauthorized, invalid or missing authentication"
// @Failure 404 {object} responsemodels.ServiceResponse "User not found with provided ID in the authenticated tenant"
// @Failure 500 {object} responsemodels.InternalServerErrorResponse "Internal server error during status update"
// @Router /users/{id}/disable [put]
// @Security ApiKeyAuth
func (h *Handler) DisableUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return BadRequest(ctx, "Invalid request: 'id' in path parameter is required and cannot be empty.")
	}
	resp, err := h.UserService.DisbaleUser(ctx.Context(), uuid.MustParse(id))
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			log.Printf("Unexpected error while fetching users : %v", err)
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while fetching users: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "The user was successfully disabled",
		Data:    resp.Message,
	})
}
