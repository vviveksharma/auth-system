package projectcontrollers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	projectmodels "github.com/vviveksharma/auth/internal/models/projectModels"
	dbmodels "github.com/vviveksharma/auth/models"
)

func projectErrResp(ctx *fiber.Ctx, err error) error {
	if svcErr, ok := err.(*dbmodels.ServiceResponse); ok {
		return ctx.Status(svcErr.Code).JSON(svcErr)
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(dbmodels.InternalServerErrorResponse{
		Code:    fiber.StatusInternalServerError,
		Message: "An unexpected error occurred while processing your request.",
	})
}

func parseProjectId(ctx *fiber.Ctx) (uuid.UUID, error) {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{ // #nosec G104 -- fiber handles response errors
			Code:    fiber.StatusBadRequest,
			Message: "invalid project id",
		})
	}
	return id, err
}

func parseOrgIdParam(ctx *fiber.Ctx) (uuid.UUID, error) {
	id, err := uuid.Parse(ctx.Params("orgId"))
	if err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{ // #nosec G104 -- fiber handles response errors
			Code:    fiber.StatusBadRequest,
			Message: "invalid organization id",
		})
	}
	return id, err
}

// ListProjects handles GET /organizations/:orgId/projects
func (h *ProjectHandler) ListProjects(ctx *fiber.Ctx) error {
	orgId, err := parseOrgIdParam(ctx)
	if err != nil {
		return nil
	}
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)

	resp, err := h.ProjectService.ListProjects(ctx.Context(), orgId, page, limit)
	if err != nil {
		return projectErrResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// CreateProject handles POST /organizations/:orgId/projects
func (h *ProjectHandler) CreateProject(ctx *fiber.Ctx) error {
	orgId, err := parseOrgIdParam(ctx)
	if err != nil {
		return nil
	}
	var req projectmodels.CreateProjectRequestBody
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Invalid request payload.",
		})
	}
	resp, err := h.ProjectService.CreateProject(ctx.Context(), orgId, &req)
	if err != nil {
		return projectErrResp(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// GetProjectDetail handles GET /projects/:id/details
func (h *ProjectHandler) GetProjectDetail(ctx *fiber.Ctx) error {
	projectId, err := parseProjectId(ctx)
	if err != nil {
		return nil
	}

	dateStr := ctx.Query("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "invalid date format, expected YYYY-MM-DD",
		})
	}

	resp, err := h.ProjectService.GetProjectDetails(ctx.Context(), projectId, date)
	if err != nil {
		return projectErrResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// GetProvidersBreakdown handles GET /projects/:id/providers-breakdown
func (h *ProjectHandler) GetProvidersBreakdown(ctx *fiber.Ctx) error {
	projectId, err := parseProjectId(ctx)
	if err != nil {
		return nil
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	startStr := ctx.Query("start_date", startOfMonth.Format("2006-01-02"))
	endStr := ctx.Query("end_date", now.Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "invalid start_date format, expected YYYY-MM-DD",
		})
	}
	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "invalid end_date format, expected YYYY-MM-DD",
		})
	}

	resp, err := h.ProjectService.GetProvidersBreakdown(ctx.Context(), projectId, startDate, endDate)
	if err != nil {
		return projectErrResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// UpdateProject handles PUT /projects/:id
func (h *ProjectHandler) UpdateProject(ctx *fiber.Ctx) error {
	projectId, err := parseProjectId(ctx)
	if err != nil {
		return nil
	}
	var req projectmodels.UpdateProjectRequestBody
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Invalid request payload.",
		})
	}
	resp, err := h.ProjectService.UpdateProject(ctx.Context(), projectId, &req)
	if err != nil {
		return projectErrResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// DeleteProject handles DELETE /projects/:id
func (h *ProjectHandler) DeleteProject(ctx *fiber.Ctx) error {
	projectId, err := parseProjectId(ctx)
	if err != nil {
		return nil
	}

	confirm := ctx.QueryBool("confirm", false)
	if !confirm {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "confirm=true is required to delete a project",
		})
	}

	resp, err := h.ProjectService.DeleteProject(ctx.Context(), projectId)
	if err != nil {
		return projectErrResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// GetProjectErrors handles GET /projects/:id/errors
func (h *ProjectHandler) GetProjectErrors(ctx *fiber.Ctx) error {
	projectId, err := parseProjectId(ctx)
	if err != nil {
		return nil
	}

	dateStr := ctx.Query("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "invalid date format, expected YYYY-MM-DD",
		})
	}
	limit := ctx.QueryInt("limit", 50)

	resp, err := h.ProjectService.GetProjectErrors(ctx.Context(), projectId, date, limit)
	if err != nil {
		return projectErrResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
