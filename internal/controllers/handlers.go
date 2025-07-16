package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/services"
	dbmodels "github.com/vviveksharma/auth/models"
)

type Handler struct {
	UserService   services.UserService
	RoleService   services.RoleService
	AuthService   services.AuthService
	TenantService services.TenantService
}

func NewHandler(userService services.UserService, roleService services.RoleService, authService services.AuthService, tenantService services.TenantService) (*Handler, error) {
	return &Handler{
		UserService:   userService,
		RoleService:   roleService,
		AuthService:   authService,
		TenantService: tenantService,
	}, nil
}

func (h *Handler) Welcome(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "GuardRail is up and working",
	})
}
