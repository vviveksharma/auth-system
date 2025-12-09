package tenantcontrollers

import (
	"github.com/gofiber/fiber/v2"
	tenantservices "github.com/vviveksharma/auth/internal/services/tenant-services"
	dbmodels "github.com/vviveksharma/auth/models"
)

type TenantHandler struct {
	TenantUserService    tenantservices.ITenantUserService
	TenantRoleService    tenantservices.ITenantRoleService
	TenantService        tenantservices.ITenantService
	TenantMessageService tenantservices.ITenantMessageService
}

func NewTenantHandler(tenantUserService tenantservices.ITenantUserService, tenantRoleService tenantservices.ITenantRoleService, tenantService tenantservices.ITenantService, tenantMessageService tenantservices.ITenantMessageService) (*TenantHandler, error) {
	return &TenantHandler{
		TenantUserService:    tenantUserService,
		TenantRoleService:    tenantRoleService,
		TenantService:        tenantService,
		TenantMessageService: tenantMessageService,
	}, nil
}

func (h *TenantHandler) Welcome(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "GuardRail tenant server is up and working",
	})
}
