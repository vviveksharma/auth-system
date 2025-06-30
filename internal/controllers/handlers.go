package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/services"
	dbmodels "github.com/vviveksharma/auth/models"
)

type Handler struct {
	UserService services.UserService
	RoleService services.RoleService
	AuthService services.AuthService
}

func NewHandler(userService services.UserService, roleService services.RoleService, authservice services.AuthService) (*Handler, error) {
	return &Handler{
		UserService: userService,
		RoleService: roleService,
		AuthService: authservice,
	}, nil
}

func (h *Handler) Welcome(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "Auth system is working",
	})
}
