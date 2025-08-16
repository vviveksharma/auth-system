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

func BadRequest(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.BadRequestResponse{
		Code:    fiber.StatusBadRequest,
		Message: message,
	})
}

func Unauthorized(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(dbmodels.UnauthorizedResponse{
		Code:    fiber.StatusUnauthorized,
		Message: message,
	})
}

func Conflict(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusConflict).JSON(dbmodels.ConflictResponse{
		Code:    fiber.StatusConflict,
		Message: message,
	})
}

func UnprocessableEntity(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(dbmodels.StatusUnprocessableEntityResponse{
		Code:    fiber.StatusUnprocessableEntity,
		Message: "Invalid request payload. Please ensure the request body is properly formatted.",
	})
}

func InternalServerError(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(dbmodels.InternalServerErrorResponse{
		Code:    fiber.StatusInternalServerError,
		Message: "An unexpected error occurred while processing your request.",
	})
}
