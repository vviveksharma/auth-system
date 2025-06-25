package controllers

import (
	"github.com/gofiber/fiber/v2"
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
