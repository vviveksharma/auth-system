package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/models"
	dbmodels "github.com/vviveksharma/auth/models"
)

func (h *Handler) ListAllRoles(ctx *fiber.Ctx) error {
	resp, err := h.RoleService.ListAllRoles()
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

func (h *Handler) VerifyRole(ctx *fiber.Ctx) error {
	var req *models.VerifyRoleRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusBadGateway,
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

func (h *Handler) CreateCustomRole(ctx *fiber.Ctx) error {
	roleName := ctx.Params("role")
	if roleName == "" {
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Role name parameter is missing or empty. Please provide a valid role name in the URL path parameter.",
		}
	}
	resp, err := h.RoleService.CreateCustomRole(roleName)
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
