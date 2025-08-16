package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/models"
	responsemodels "github.com/vviveksharma/auth/models"
)

func (h *Handler) CreateUser(ctx *fiber.Ctx) error {
	var req *models.UserRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &responsemodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	if req.Email == "" || req.Name == "" || req.Password == "" {
		log.Println("the requestBody: ", req)
		return &responsemodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "either name or type is missing in the request Body",
		}
	}
	fmt.Println("the userdetails from the request ", req)
	resp, err := h.UserService.CreateUser(req)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}

// GetUserDetails retrieves details of the currently authenticated user.
//
// @Summary Get Authenticated User Details
// @Description Returns the details of the user currently authenticated via the API key or token. This endpoint is useful for profile pages or user dashboards.
// @Tags User
// @Produce json
// @Success 200 {object} responsemodels.ServiceResponse "User details successfully retrieved"
// @Failure 401 {object} responsemodels.ServiceResponse "Unauthorized, invalid or missing authentication"
// @Failure 500 {object} responsemodels.ServiceResponse "Internal server error"
// @Router /user/details [get]
// @Security ApiKeyAuth
func (h *Handler) GetUserDetails(ctx *fiber.Ctx) error {
	req := &models.GetUserDetailsRequest{}
	userId := ctx.Locals("userId").(string)
	fmt.Println("the userid: ", userId)
	req.Id = userId
	resp, err := h.UserService.GetUserDetails(req)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
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
// @Description Allows the authenticated user to update their profile information such as name, email, or other editable fields. Requires authentication.
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UpdateUserRequest true "Fields to update for the user profile"
// @Success 200 {object} responsemodels.ServiceResponse "User details updated successfully"
// @Failure 401 {object} responsemodels.ServiceResponse "Unauthorized, invalid or missing authentication"
// @Failure 422 {object} responsemodels.ServiceResponse "Unprocessable entity, invalid input"
// @Failure 500 {object} responsemodels.ServiceResponse "Internal server error"
// @Router /user/details [put]
// @Security ApiKeyAuth
func (h *Handler) UpdateUserDetails(ctx *fiber.Ctx) error {
	req := &models.UpdateUserRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(responsemodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		})
	}
	userId := ctx.Locals("userId").(string)
	fmt.Println("the userid: ", userId)
	resp, err := h.UserService.UpdateUserDetails(req, userId)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: "an unexpected error occurred",
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
// @Failure 401 {object} responsemodels.ServiceResponse "Unauthorized, invalid or missing authentication"
// @Failure 404 {object} responsemodels.ServiceResponse "User not found"
// @Failure 500 {object} responsemodels.ServiceResponse "Internal server error"
// @Router /user/{id} [get]
// @Security ApiKeyAuth
func (h *Handler) GetUserByIdDetails(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	resp, err := h.UserService.GetUserById(userId)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: "an unexpected error occurred",
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
// @Description Assigns a specific role to a user identified by their user ID. Only users with sufficient privileges (e.g., admins) can perform this action.
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Unique identifier of the user to assign a role"
// @Param request body models.AssignRoleRequest true "Role assignment details"
// @Success 200 {object} responsemodels.ServiceResponse "Role assigned successfully"
// @Failure 401 {object} responsemodels.ServiceResponse "Unauthorized, invalid or missing authentication"
// @Failure 404 {object} responsemodels.ServiceResponse "User not found"
// @Failure 422 {object} responsemodels.ServiceResponse "Unprocessable entity, invalid input"
// @Failure 500 {object} responsemodels.ServiceResponse "Internal server error"
// @Router /user/{id}/role [post]
// @Security ApiKeyAuth
func (h *Handler) AssignUserRole(ctx *fiber.Ctx) error {
	var req *models.AssignRoleRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(responsemodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		})
	}
	userId := ctx.Params("id")
	resp, err := h.UserService.AssignUserRole(req, userId)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: "an unexpected error occurred",
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
// @Description Registers a new user in the system under a specific tenant. This endpoint is typically used for onboarding new users. Requires all mandatory fields such as name, email, and password.
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UserRequest true "User registration details including name, email, and password"
// @Success 200 {object} models.UserResponse "User registered successfully"
// @Failure 400 {object} responsemodels.ServiceResponse "Bad request, missing required fields"
// @Failure 409 {object} responsemodels.ServiceResponse "Conflict, user already exists"
// @Failure 422 {object} responsemodels.ServiceResponse "Unprocessable entity, invalid input"
// @Failure 500 {object} responsemodels.ServiceResponse "Internal server error"
// @Router /user/register [post]
func (h *Handler) RegisterUser(ctx *fiber.Ctx) error {
	var req *models.UserRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &responsemodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	if req.Email == "" || req.Name == "" || req.Password == "" {
		log.Println("the requestBody: ", req)
		return &responsemodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "either name or type is missing in the request Body",
		}
	}
	resp, err := h.UserService.RegisterUser(req, ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
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
// @Description Initiates the password reset process for a user by sending a reset link or OTP to the user's email.
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Password reset request details"
// @Success 200 {object} responsemodels.ServiceResponse "Password reset initiated successfully"
// @Failure 422 {object} responsemodels.ServiceResponse "Unprocessable entity, invalid input"
// @Failure 500 {object} responsemodels.ServiceResponse "Internal server error"
// @Router /user/password/reset [post]
func (h *Handler) ResetUserPassword(ctx *fiber.Ctx) error {
	var req *models.ResetPasswordRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(responsemodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		})
	}
	resp, err := h.UserService.ResetPassword(req)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: "an unexpected error occurred",
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
// @Description Sets a new password for the user after verifying the OTP sent to their email. Requires email, OTP, new password, and confirmation.
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UserVerifyOTPRequest true "OTP verification and new password details"
// @Success 200 {object} responsemodels.ServiceResponse "Password updated successfully"
// @Failure 400 {object} responsemodels.ServiceResponse "Bad request, missing required fields"
// @Failure 409 {object} responsemodels.ServiceResponse "Conflict, password confirmation failed"
// @Failure 422 {object} responsemodels.ServiceResponse "Unprocessable entity, invalid input"
// @Failure 500 {object} responsemodels.ServiceResponse "Internal server error"
// @Router /user/password/set [post]
func (h *Handler) SetUserPassword(ctx *fiber.Ctx) error {
	var req *models.UserVerifyOTPRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(responsemodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		})
	}
	if req.Email == "" || req.OTP == "" || req.ConfirmPassword == "" || req.NewPassword == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(responsemodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "One or more required fields (email, OTP, new password, confirm password) are missing in the request. Please ensure all fields are provided.",
		})
	}
	if req.NewPassword != req.ConfirmPassword {
		log.Printf("Password mismatch for email: %s", req.Email)
		return ctx.Status(fiber.StatusConflict).JSON(responsemodels.ServiceResponse{
			Code:    fiber.StatusConflict,
			Message: "Password confirmation failed: new password and confirmation do not match. Please ensure both fields are identical.",
		})
	}
	resp, err := h.UserService.SetPassword(req)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: "an unexpected error occurred",
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "The password was succcessfully updated",
		Data:    resp,
	})
}
