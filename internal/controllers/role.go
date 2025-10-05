package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/internal/models"
	responsemodels "github.com/vviveksharma/auth/models"
)

// ListAllRoles godoc
// @Summary      List all roles
// @Description  Retrieves all roles from the system. Optionally, you can filter roles by type using the 'type' query parameter. If not provided, the default type is used.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        roleTypeFlag  query     string  false  "Role type to filter (e.g., 'admin', 'user', 'default'). If not provided, defaults to 'default'."
// @Success      200   {object}  responsemodels.ServiceResponse  "Roles fetched successfully. Data contains the list of roles."
// @Failure      500   {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while fetching roles."
// @Router       /roles [get]
func (h *Handler) ListAllRoles(ctx *fiber.Ctx) error {
	flag := ctx.Query("roleTypeFlag")
	if flag == "" {
		flag = "custom"
	} else if flag != "custom" && flag != "default" {
		return BadRequest(ctx, "the roleTypeFlag could be only default or user")
	}
	resp, err := h.RoleService.ListRoles(flag, 1, 5, ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return InternalServerError(ctx, "Internal server error occurred: "+err.Error())
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(&responsemodels.ServiceResponse{
		Code:    200,
		Message: "Roles fetched successfully",
		Data:    resp,
	})
}

// VerifyRole godoc
// @Summary      Verify role
// @Description  Verifies if a role exists by roleId and roleName. Returns 404 if either is missing or not found, 422 if the request body is invalid.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        request  body      models.VerifyRoleRequest  true  "Verify Role Request. Requires both 'roleId' and 'roleName'."
// @Success      200      {object}  responsemodels.ServiceResponse  "Role verified successfully."
// @Failure      400      {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if either 'roleId' or 'roleName' is missing or invalid."
// @Failure      422      {object}  responsemodels.StatusUnprocessableEntityResponse  "Unprocessable Entity. This occurs if the request body cannot be parsed."
// @Failure      500      {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while verifying the role."
// @Router       /roles/verify [post]
func (h *Handler) VerifyRole(ctx *fiber.Ctx) error {
	var req *models.VerifyRoleRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return UnprocessableEntity(ctx)
	}
	if req.RoleId == "" || req.RoleName == "" {
		return BadRequest(ctx, "Invalid request: 'role_name' and 'role_id' fields are required and cannot be empty.")
	}
	resp, err := h.RoleService.VerifyRole(req)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return InternalServerError(ctx, "Internal server error occurred: "+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "Role verified successfully",
		Data:    resp,
	})
}

// CreateCustomRole godoc
// @Summary      Create custom role
// @Description  Creates a new custom role with specified routes. Requires a unique role name and a list of routes. Returns 422 if the request is invalid, 500 for internal errors.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateCustomRole  true  "Create Custom Role Request. Requires 'roleName' and 'routes' fields."
// @Success      200      {object}  responsemodels.ServiceResponse  "Role created successfully."
// @Failure      400      {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if 'roleName' or 'routes' are missing or invalid."
// @Failure      422      {object}  responsemodels.StatusUnprocessableEntityResponse  "Unprocessable Entity. This occurs if the request body cannot be parsed."
// @Failure      500      {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while creating the role."
// @Router       /roles/ [post]
func (h *Handler) CreateCustomRole(ctx *fiber.Ctx) error {
	var req models.CreateCustomRole
	err := ctx.BodyParser(&req)
	if err != nil {
		return UnprocessableEntity(ctx)
	}
	if req.Name == "" || req.Description == "" || req.Permissions == nil || req.DisplayName == "" {
		return BadRequest(ctx, "Invalid request: 'roleName' and 'routes' fields are required and cannot be empty.")
	}
	resp, err := h.RoleService.CreateCustomRole(&req, ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return InternalServerError(ctx, "Internal server error occurred: "+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}

// UpdateRolePermission godoc
// @Summary      Update role permissions
// @Description  Adds or removes permissions from a role. Requires role name and lists of permissions to add or remove. Returns 422 if the request is invalid, 500 for internal errors.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        request  body      models.UpdateRolePermissions  true  "Update Role Permissions Request. Requires 'roleName', 'addPermisions', and 'removePermissions'."
// @Success      200      {object}  responsemodels.ServiceResponse  "Role permissions updated successfully."
// @Failure      400      {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if required fields are missing or invalid."
// @Failure      422      {object}  responsemodels.StatusUnprocessableEntityResponse  "Unprocessable Entity. This occurs if the request body cannot be parsed."
// @Failure      500      {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while updating permissions."
// @Router       /roles/:id/permissions [put]
func (h *Handler) UpdateRolePermission(ctx *fiber.Ctx) error {
	var req models.UpdateRolePermissions
	err := ctx.BodyParser(&req)
	if err != nil {
		return UnprocessableEntity(ctx)
	}
	if req.AddPermisions == nil || req.RemovePermissions == nil || req.RoleName == "" {
		return BadRequest(ctx, "Invalid request: 'add_permissions', 'role_name' and 'remove_permissions' fields are required and cannot be empty.")
	}
	id := ctx.Params("id")
	resp, err := h.RoleService.UpdateRolePermission(&req, uuid.MustParse(id), ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return InternalServerError(ctx, "Internal server error occurred: "+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}

// DeleteCustomRole godoc
// @Summary      Delete custom role
// @Description  Deletes a custom role by its ID. Only custom roles can be deleted; system roles are protected. Any users currently assigned this role will need to be reassigned before deletion.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID to delete. Must be a valid UUID format."
// @Success      200  {object}  responsemodels.ServiceResponse  "Role deleted successfully."
// @Failure      400  {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if the 'id' path parameter is missing, empty, or not a valid UUID format."
// @Failure      409  {object}  responsemodels.ConflictResponse  "Conflict. This occurs if the role is currently assigned to users and cannot be deleted."
// @Failure      500  {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while deleting the role."
// @Router       /roles/{id} [delete]
func (h *Handler) DeleteCustomRole(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return BadRequest(ctx, "Invalid request: 'id' in path parameter is required and cannot be empty.")
	}
	resp, err := h.RoleService.DeleteRole(uuid.MustParse(id), ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			log.Printf("Unexpected error while deleting role with id %s: %v", id, err)
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting role: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "The role was successfully deleted",
		Data:    resp,
	})
}

// EnableRole godoc
// @Summary      Enable role
// @Description  Enables a role by its ID, making it available for assignment to users. Only disabled roles can be enabled. System roles are always enabled by default.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID to enable. Must be a valid UUID format."
// @Success      200  {object}  responsemodels.ServiceResponse  "Role enabled successfully."
// @Failure      400  {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if the 'id' path parameter is missing, empty, or not a valid UUID format."
// @Failure      409  {object}  responsemodels.ConflictResponse  "Conflict. This occurs if the role is already enabled."
// @Failure      500  {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while enabling the role."
// @Router       /roles/{id}/enable [put]
func (h *Handler) EnableRole(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return BadRequest(ctx, "Invalid request: 'id' in path parameter is required and cannot be empty.")
	}
	resp, err := h.RoleService.EnableRole(uuid.MustParse(id), ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			log.Printf("Unexpected error while enabling the role with id %s: %v", id, err)
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while enabling the role: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

// DisableRole godoc
// @Summary      Disable role
// @Description  Disables a role by its ID, preventing it from being assigned to new users. Existing users with this role will retain it but new assignments are blocked. System roles cannot be disabled.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID to disable. Must be a valid UUID format."
// @Success      200  {object}  responsemodels.ServiceResponse  "Role disabled successfully."
// @Failure      400  {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if the 'id' path parameter is missing, empty, or not a valid UUID format."
// @Failure      409  {object}  responsemodels.ConflictResponse  "Conflict. This occurs if the role is already disabled."
// @Failure      500  {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while disabling the role."
// @Router       /roles/{id}/disable [put]
func (h *Handler) DisableRole(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return BadRequest(ctx, "Invalid request: 'id' in path parameter is required and cannot be empty.")
	}
	resp, err := h.RoleService.DisableRole(uuid.MustParse(id), ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			log.Printf("Unexpected error while disabling the role with id %s: %v", id, err)
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while disabling the role: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

// GetRolePermissions godoc
// @Summary      Get role permissions
// @Description  Retrieves all routes and permissions associated with a specific role by role ID. This is useful for role management and permission auditing. Returns detailed permission structure including HTTP methods and route information.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID to retrieve permissions for. Must be a valid UUID format."
// @Success      200  {object}  responsemodels.ServiceResponse  "Role permissions retrieved successfully. Data contains the detailed permission structure."
// @Failure      400  {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if the 'id' path parameter is missing, empty, or not a valid UUID format."
// @Failure      500  {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while retrieving role permissions."
// @Router       /roles/{id}/permissions [get]
// @Security     ApiKeyAuth
func (h *Handler) GetRolePermissions(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return BadRequest(ctx, "Invalid request: 'id' in path parameter is required and cannot be empty.")
	}
	resp, err := h.RoleService.GetRouteDetails(uuid.MustParse(id), ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(serviceErr)
		} else {
			log.Printf("Unexpected error while disabling the role with id %s: %v", id, err)
			return ctx.Status(500).JSON(responsemodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while disabling the role: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "Role details retrieved successfully",
		Data:    resp,
	})
}
