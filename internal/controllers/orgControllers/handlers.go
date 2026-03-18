package orgcontrollers

import (
	"github.com/gofiber/fiber/v2"
	orgservices "github.com/vviveksharma/auth/internal/services/org-services"
	dbmodels "github.com/vviveksharma/auth/models"
)

type OrgHandler struct {
	OrgService orgservices.IOrgServiceInterface
}

func NewOrgHandler(orgService orgservices.IOrgServiceInterface) (*OrgHandler, error) {
	return &OrgHandler{
		OrgService: orgService,
	}, nil
}

func (h *OrgHandler) Welcome(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Org service is up and running",
	})
}
