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
// @Summary      List all roles with pagination
// @Description  Retrieves all roles from the system with pagination support. Filter roles by type using 'roleTypeFlag' query parameter ('custom' or 'default'). Includes role routes in response.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        roleTypeFlag  query     string  false  "Role type to filter: 'custom' or 'default'. If not provided, defaults to 'custom'."
// @Param        page  query     int  false  "Page number (default: 1)"
// @Param        page_size  query     int  false  "Number of items per page (default: 5)"
// @Success      200   {object}  responsemodels.ServiceResponse  "Roles fetched successfully with pagination metadata. Data includes role details and associated routes."
// @Failure      400   {object}  responsemodels.BadRequestResponse  "Bad request. Invalid roleTypeFlag parameter (must be 'custom' or 'default')."
// @Failure      500   {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while fetching roles."
// @Router       /roles [get]
// @Security     ApiKeyAuth
func (h *Handler) ListAllRoles(ctx *fiber.Ctx) error {
	flag := ctx.Query("roleTypeFlag")
	if flag == "" {
		flag = "custom"
	} else if flag != "custom" && flag != "default" {
		return BadRequest(ctx, "the roleTypeFlag could be only default or user")
	}

	// Parse pagination parameters
	page := ctx.QueryInt("page", 1)
	pageSize := ctx.QueryInt("page_size", 5)

	resp, err := h.RoleService.ListRoles(flag, page, pageSize, ctx.Context())
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

// CreateCustomRole godoc
// @Summary      Create custom role
// @Description  Creates a new custom role with specified permissions. Requires name, display_name, description, and permissions array with route, methods, and description for each permission.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateCustomRole  true  "Create Custom Role Request. Requires 'name', 'display_name', 'description', and 'Permissions' array."
// @Success      200      {object}  responsemodels.ServiceResponse  "Role created successfully."
// @Failure      400      {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if required fields (name, display_name, description, or Permissions) are missing or invalid."
// @Failure      409      {object}  responsemodels.ConflictResponse  "Conflict. This occurs if a role with the same name already exists."
// @Failure      422      {object}  responsemodels.StatusUnprocessableEntityResponse  "Unprocessable Entity. This occurs if the request body cannot be parsed."
// @Failure      500      {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while creating the role."
// @Router       /roles/ [post]
// @Security     ApiKeyAuth
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
// @Description  Adds or removes permissions from a role. Requires role name (role field), and lists of permissions to add (add_permissions) or remove (remove_permissions). Each permission includes route, methods array, and description.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Role ID to update permissions for. Must be a valid UUID format."
// @Param        request  body      models.UpdateRolePermissions  true  "Update Role Permissions Request. Requires 'role', 'add_permissions', and 'remove_permissions'."
// @Success      200      {object}  responsemodels.ServiceResponse  "Role permissions updated successfully."
// @Failure      400      {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if required fields are missing or invalid."
// @Failure      404      {object}  responsemodels.ServiceResponse  "Not Found. This occurs if the role with the specified name doesn't exist."
// @Failure      409      {object}  responsemodels.ConflictResponse  "Conflict. This occurs if trying to modify a default/system role."
// @Failure      422      {object}  responsemodels.StatusUnprocessableEntityResponse  "Unprocessable Entity. This occurs if the request body cannot be parsed."
// @Failure      500      {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while updating permissions."
// @Router       /roles/{id}/permissions [put]
// @Security     ApiKeyAuth
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
// @Description  Deletes a custom role by its ID. Only custom roles can be deleted; system/default roles are protected. The role must not be currently assigned to any users.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID to delete. Must be a valid UUID format."
// @Success      200  {object}  responsemodels.ServiceResponse  "Role deleted successfully."
// @Failure      400  {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if the 'id' path parameter is missing, empty, or not a valid UUID format."
// @Failure      404  {object}  responsemodels.ServiceResponse  "Not Found. This occurs if no role exists with the provided ID or if trying to delete a default/system role."
// @Failure      500  {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while deleting the role."
// @Router       /roles/{id} [delete]
// @Security     ApiKeyAuth
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
			return ctx.Status(500).JSON(&responsemodels.ServiceResponse{
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
// @Description  Enables a role by its ID, making it available for assignment to users. Only disabled custom roles can be enabled. System/default roles are always enabled.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID to enable. Must be a valid UUID format."
// @Success      200  {object}  responsemodels.ServiceResponse  "Role enabled successfully."
// @Failure      400  {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if the 'id' path parameter is missing, empty, or not a valid UUID format."
// @Failure      403  {object}  responsemodels.ServiceResponse  "Forbidden. This occurs if trying to modify a system/default role."
// @Failure      404  {object}  responsemodels.ServiceResponse  "Not Found. This occurs if no role exists with the provided ID."
// @Failure      409  {object}  responsemodels.ConflictResponse  "Conflict. This occurs if the role is already enabled."
// @Failure      500  {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while enabling the role."
// @Router       /roles/{id}/enable [put]
// @Security     ApiKeyAuth
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
// @Description  Disables a role by its ID, preventing it from being assigned to new users. Existing users with this role will retain it but new assignments are blocked. System/default roles cannot be disabled.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID to disable. Must be a valid UUID format."
// @Success      200  {object}  responsemodels.ServiceResponse  "Role disabled successfully."
// @Failure      400  {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if the 'id' path parameter is missing, empty, or not a valid UUID format."
// @Failure      403  {object}  responsemodels.ServiceResponse  "Forbidden. This occurs if trying to modify a system/default role."
// @Failure      404  {object}  responsemodels.ServiceResponse  "Not Found. This occurs if no role exists with the provided ID."
// @Failure      409  {object}  responsemodels.ConflictResponse  "Conflict. This occurs if the role is already disabled."
// @Failure      500  {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while disabling the role."
// @Router       /roles/{id}/disable [put]
// @Security     ApiKeyAuth
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
// @Description  Retrieves all routes and permissions associated with a specific role by role ID. Returns detailed permission structure classified by HTTP methods (GET, POST, PUT, DELETE, etc.), route information, role details, and processing timestamp. Supports both custom and system roles.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID to retrieve permissions for. Must be a valid UUID format. Supports both custom and system role IDs."
// @Success      200  {object}  responsemodels.ServiceResponse  "Role permissions retrieved successfully. Data contains Routes (classified by method), RoutesJSON (formatted JSON string), RoleInfo, and ProcessedAt timestamp."
// @Failure      400  {object}  responsemodels.BadRequestResponse  "Bad Request. This occurs if the 'id' path parameter is missing, empty, or not a valid UUID format."
// @Failure      404  {object}  responsemodels.ServiceResponse  "Not Found. This occurs if no role exists with the provided ID."
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
