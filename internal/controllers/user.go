package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/models"
	dbmodels "github.com/vviveksharma/auth/models"
)

func (h *Handler) CreateUser(ctx *fiber.Ctx) error {
	var req *models.UserRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	if req.Email == "" || req.Name == "" || req.Password == "" {
		log.Println("the requestBody: ", req)
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "either name or type is missing in the request Body",
		}
	}
	fmt.Println("the userdetails from the request ", req)
	resp, err := h.UserService.CreateUser(req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
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
// @Success 200 {object} dbmodels.ServiceResponse "User details successfully retrieved"
// @Failure 401 {object} dbmodels.ServiceResponse "Unauthorized, invalid or missing authentication"
// @Failure 500 {object} dbmodels.ServiceResponse "Internal server error"
// @Router /user/details [get]
// @Security ApiKeyAuth
func (h *Handler) GetUserDetails(ctx *fiber.Ctx) error {
	req := &models.GetUserDetailsRequest{}
	userId := ctx.Locals("userId").(string)
	fmt.Println("the userid: ", userId)
	req.Id = userId
	resp, err := h.UserService.GetUserDetails(req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
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
// @Success 200 {object} dbmodels.ServiceResponse "User details updated successfully"
// @Failure 401 {object} dbmodels.ServiceResponse "Unauthorized, invalid or missing authentication"
// @Failure 422 {object} dbmodels.ServiceResponse "Unprocessable entity, invalid input"
// @Failure 500 {object} dbmodels.ServiceResponse "Internal server error"
// @Router /user/details [put]
// @Security ApiKeyAuth
func (h *Handler) UpdateUserDetails(ctx *fiber.Ctx) error {
	req := &models.UpdateUserRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		})
	}
	userId := ctx.Locals("userId").(string)
	fmt.Println("the userid: ", userId)
	resp, err := h.UserService.UpdateUserDetails(req, userId)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: "an unexpected error occurred",
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
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
// @Success 200 {object} dbmodels.ServiceResponse "User details successfully retrieved"
// @Failure 401 {object} dbmodels.ServiceResponse "Unauthorized, invalid or missing authentication"
// @Failure 404 {object} dbmodels.ServiceResponse "User not found"
// @Failure 500 {object} dbmodels.ServiceResponse "Internal server error"
// @Router /user/{id} [get]
// @Security ApiKeyAuth
func (h *Handler) GetUserByIdDetails(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	resp, err := h.UserService.GetUserById(userId)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: "an unexpected error occurred",
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
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
// @Success 200 {object} dbmodels.ServiceResponse "Role assigned successfully"
// @Failure 401 {object} dbmodels.ServiceResponse "Unauthorized, invalid or missing authentication"
// @Failure 404 {object} dbmodels.ServiceResponse "User not found"
// @Failure 422 {object} dbmodels.ServiceResponse "Unprocessable entity, invalid input"
// @Failure 500 {object} dbmodels.ServiceResponse "Internal server error"
// @Router /user/{id}/role [post]
// @Security ApiKeyAuth
func (h *Handler) AssignUserRole(ctx *fiber.Ctx) error {
	var req *models.AssignRoleRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		})
	}
	userId := ctx.Params("id")
	resp, err := h.UserService.AssignUserRole(req, userId)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: "an unexpected error occurred",
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
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
// @Failure 400 {object} dbmodels.ServiceResponse "Bad request, missing required fields"
// @Failure 409 {object} dbmodels.ServiceResponse "Conflict, user already exists"
// @Failure 422 {object} dbmodels.ServiceResponse "Unprocessable entity, invalid input"
// @Failure 500 {object} dbmodels.ServiceResponse "Internal server error"
// @Router /user/register [post]
func (h *Handler) RegisterUser(ctx *fiber.Ctx) error {
	var req *models.UserRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	if req.Email == "" || req.Name == "" || req.Password == "" {
		log.Println("the requestBody: ", req)
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "either name or type is missing in the request Body",
		}
	}
	resp, err := h.UserService.RegisterUser(req, ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}
