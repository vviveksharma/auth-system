package projectcontrollers

import (
	"github.com/gofiber/fiber/v2"
	projectservice "github.com/vviveksharma/auth/internal/services/project-service"
	dbmodels "github.com/vviveksharma/auth/models"
)

type ProjectHandler struct {
	ProjectService projectservice.IProjectServiceInterface
}

func NewProjectHandler(projectService projectservice.IProjectServiceInterface) (*ProjectHandler, error) {
	return &ProjectHandler{
		ProjectService: projectService,
	}, nil
}

func (h *ProjectHandler) Welcome(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Project service is up and running ",
	})
}