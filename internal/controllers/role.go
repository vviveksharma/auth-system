package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/models"
	responsemodels "github.com/vviveksharma/auth/models"
)

// ListAllRoles godoc
// @Summary      List all roles
// @Description  Retrieves all roles from the system. Optionally, you can filter roles by type using the 'type' query parameter. If not provided, the default type is used.
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        type  query     string  false  "Role type to filter (e.g., 'admin', 'user', 'default'). If not provided, defaults to 'default'."
// @Success      200   {object}  responsemodels.ServiceResponse  "Roles fetched successfully. Data contains the list of roles."
// @Failure      500   {object}  responsemodels.InternalServerErrorResponse  "Internal server error. This occurs if there is an unexpected error while fetching roles."
// @Router       /roles [get]
func (h *Handler) ListAllRoles(ctx *fiber.Ctx) error {
	flag := ctx.Query("type")
	if flag == "" {
		flag = "default"
	}
	resp, err := h.RoleService.ListAllRoles(flag)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return InternalServerError(ctx, "Internal server error occured: "+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&responsemodels.ServiceResponse{
		Code:    200,
		Message: "Roles are as follow",
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
// @Success      200      {object}  responsemodels.ServiceResponse  "Role exists. Data contains verification result."
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
			return InternalServerError(ctx, "Internal server error occurred."+"error: "+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "Roles are as follow",
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
// @Router       /roles/custom [post]
func (h *Handler) CreateCustomRole(ctx *fiber.Ctx) error {
	var req models.CreateCustomRole
	err := ctx.BodyParser(&req)
	if err != nil {
		return UnprocessableEntity(ctx)
	}
	if req.RoleName == "" || req.Routes == nil {
		return BadRequest(ctx, "Invalid request: 'roleName' and 'routes' fields are required and cannot be empty.")
	}
	resp, err := h.RoleService.CreateCustomRole(&req)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return InternalServerError(ctx, "Internal server error occurred."+"error: "+err.Error())
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
// @Router       /roles/permissions [put]
func (h *Handler) UpdateRolePermission(ctx *fiber.Ctx) error {
	var req models.UpdateRolePermissions
	err := ctx.BodyParser(&req)
	if err != nil {
		return UnprocessableEntity(ctx)
	}
	if req.AddPermisions == nil || req.RemovePermissions == nil || req.RoleName == "" {
		return BadRequest(ctx, "Invalid request: 'add_permissions', 'role_name' and 'remove_permissions' fields are required and cannot be empty.")
	}
	resp, err := h.RoleService.UpdateRolePermission(&req)
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return InternalServerError(ctx, "Internal server error occurred."+"error: "+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}
