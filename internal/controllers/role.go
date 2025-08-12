package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/models"
	dbmodels "github.com/vviveksharma/auth/models"
)

// ListAllRoles godoc
// @Summary      List all roles
// @Description  Retrieves all roles from the system. Optionally, you can filter roles by type using the 'type' query parameter. If not provided, the default type is used.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        type  query     string  false  "Role type to filter (e.g., 'admin', 'user', 'default'). If not provided, defaults to 'default'."
// @Success      200   {object}  dbmodels.ServiceResponse  "Roles fetched successfully. Data contains the list of roles."
// @Failure      500   {object}  dbmodels.ServiceResponse  "Internal server error. This occurs if there is an unexpected error while fetching roles."
// @Router       /roles [get]
func (h *Handler) ListAllRoles(ctx *fiber.Ctx) error {
	flag := ctx.Query("type")
	if flag == "" {
		flag = "default"
	}
	resp, err := h.RoleService.ListAllRoles(flag)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Roles are as follow",
		Data:    resp,
	})
}

// VerifyRole godoc
// @Summary      Verify role
// @Description  Verifies if a role exists by roleId and roleName. Returns 404 if either is missing or not found, 422 if the request body is invalid.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        request  body      models.VerifyRoleRequest  true  "Verify Role Request. Requires both 'roleId' and 'roleName'."
// @Success      200      {object}  dbmodels.ServiceResponse  "Role exists. Data contains verification result."
// @Failure      404      {object}  dbmodels.ServiceResponse  "Role not found or missing required fields. This occurs if either 'roleId' or 'roleName' is missing or invalid."
// @Failure      422      {object}  dbmodels.ServiceResponse  "Unprocessable Entity. This occurs if the request body cannot be parsed."
// @Router       /roles/verify [post]
func (h *Handler) VerifyRole(ctx *fiber.Ctx) error {
	var req *models.VerifyRoleRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	if req.RoleId == "" || req.RoleName == "" {
		return &dbmodels.ServiceResponse{
			Code:    404,
			Message: "either the roleId or rolename is missing",
		}
	}
	resp, err := h.RoleService.VerifyRole(req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "Roles are as follow",
		Data:    resp,
	})
}

// CreateCustomRole godoc
// @Summary      Create custom role
// @Description  Creates a new custom role with specified routes. Requires a unique role name and a list of routes. Returns 422 if the request is invalid, 500 for internal errors.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateCustomRole  true  "Create Custom Role Request. Requires 'roleName' and 'routes' fields."
// @Success      200      {object}  dbmodels.ServiceResponse  "Role created successfully."
// @Failure      400      {object}  dbmodels.ServiceResponse  "Bad Request. This occurs if 'roleName' or 'routes' are missing or invalid."
// @Failure      422      {object}  dbmodels.ServiceResponse  "Unprocessable Entity. This occurs if the request body cannot be parsed."
// @Failure      500      {object}  dbmodels.ServiceResponse  "Internal server error. This occurs if there is an unexpected error while creating the role."
// @Router       /roles/custom [post]
func (h *Handler) CreateCustomRole(ctx *fiber.Ctx) error {
	var req models.CreateCustomRole
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Invalid request payload. Please ensure the request body is properly formatted.",
		})
	}
	if req.RoleName == "" || req.Routes == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. Please ensure the request body is properly formatted.",
		})
	}
	resp, err := h.RoleService.CreateCustomRole(&req)
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

// UpdateRolePermission godoc
// @Summary      Update role permissions
// @Description  Adds or removes permissions from a role. Requires role name and lists of permissions to add or remove. Returns 422 if the request is invalid, 500 for internal errors.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        request  body      models.UpdateRolePermissions  true  "Update Role Permissions Request. Requires 'roleName', 'addPermisions', and 'removePermissions'."
// @Success      200      {object}  dbmodels.ServiceResponse  "Role permissions updated successfully."
// @Failure      400      {object}  dbmodels.ServiceResponse  "Bad Request. This occurs if required fields are missing or invalid."
// @Failure      422      {object}  dbmodels.ServiceResponse  "Unprocessable Entity. This occurs if the request body cannot be parsed."
// @Failure      500      {object}  dbmodels.ServiceResponse  "Internal server error. This occurs if there is an unexpected error while updating permissions."
// @Router       /roles/permissions [put]
func (h *Handler) UpdateRolePermission(ctx *fiber.Ctx) error {
	var req models.UpdateRolePermissions
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Invalid request payload. Please ensure the request body is properly formatted.",
		})
	}
	if req.AddPermisions == nil || req.RemovePermissions == nil || req.RoleName == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. Please ensure the request body is properly formatted.",
		})
	}
	resp, err := h.RoleService.UpdateRolePermission(&req)
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
