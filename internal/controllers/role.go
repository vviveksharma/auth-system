package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/models"
	dbmodels "github.com/vviveksharma/auth/models"
)

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
	var req models.CreateCustomRole
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
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

func (h *Handler) UpdateRolePermission(ctx *fiber.Ctx) error {
	var req models.UpdateRolePermissions
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
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
