package orgcontrollers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	orgmodels "github.com/vviveksharma/auth/internal/models/orgModels"
	dbmodels "github.com/vviveksharma/auth/models"
)

func errResp(ctx *fiber.Ctx, err error) error {
	if svcErr, ok := err.(*dbmodels.ServiceResponse); ok {
		return ctx.Status(svcErr.Code).JSON(svcErr)
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(dbmodels.InternalServerErrorResponse{
		Code:    fiber.StatusInternalServerError,
		Message: "An unexpected error occurred while processing your request.",
	})
}

func parseOrgId(ctx *fiber.Ctx) (uuid.UUID, error) {
	orgId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{ // #nosec G104 -- fiber handles response errors
			Code:    fiber.StatusBadRequest,
			Message: "invalid organization id",
		})
	}
	return orgId, err
}

// ListOrgs handles GET /organizations.
func (h *OrgHandler) ListOrgs(ctx *fiber.Ctx) error {
	resp, err := h.OrgService.ListOrgs(ctx.Context())
	if err != nil {
		return errResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// GetOrg handles GET /organizations/:id.
func (h *OrgHandler) GetOrg(ctx *fiber.Ctx) error {
	orgId, err := parseOrgId(ctx)
	if err != nil {
		return nil
	}
	resp, err := h.OrgService.GetOrgById(ctx.Context(), orgId)
	if err != nil {
		return errResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// CreateOrg handles POST /organizations.
func (h *OrgHandler) CreateOrg(ctx *fiber.Ctx) error {
	var req orgmodels.CreateOrgRequestBody
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(dbmodels.StatusUnprocessableEntityResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Invalid request payload. Please ensure the request body is properly formatted.",
		})
	}
	resp, err := h.OrgService.CreateOrg(ctx.Context(), &req)
	if err != nil {
		return errResp(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// SwitchOrg handles POST /organizations/:id/switch.
func (h *OrgHandler) SwitchOrg(ctx *fiber.Ctx) error {
	orgId, err := parseOrgId(ctx)
	if err != nil {
		return nil
	}
	resp, err := h.OrgService.SwitchOrg(ctx.Context(), orgId)
	if err != nil {
		return errResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// UpdateOrg handles PUT /organizations/:id.
func (h *OrgHandler) UpdateOrg(ctx *fiber.Ctx) error {
	orgId, err := parseOrgId(ctx)
	if err != nil {
		return nil
	}
	var req orgmodels.UpdateOrgRequestBody
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(dbmodels.StatusUnprocessableEntityResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Invalid request payload. Please ensure the request body is properly formatted.",
		})
	}
	resp, err := h.OrgService.UpdateOrg(ctx.Context(), orgId, &req)
	if err != nil {
		return errResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// DeleteOrg handles DELETE /organizations/:id.
// Query params: confirm=true (required), confirmation_text=<org name> (optional safety check).
func (h *OrgHandler) DeleteOrg(ctx *fiber.Ctx) error {
	orgId, err := parseOrgId(ctx)
	if err != nil {
		return nil
	}
	if !ctx.QueryBool("confirm", false) {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "confirm_required: set confirm=true to confirm deletion",
		})
	}
	if confirmationText := ctx.Query("confirmation_text"); confirmationText != "" {
		orgDetails, getErr := h.OrgService.GetOrgById(ctx.Context(), orgId)
		if getErr != nil {
			return errResp(ctx, getErr)
		}
		if confirmationText != orgDetails.Name {
			return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
				Code:    fiber.StatusBadRequest,
				Message: "confirmation_mismatch: confirmation text does not match the organization name",
			})
		}
	}
	resp, err := h.OrgService.DeleteOrg(ctx.Context(), orgId)
	if err != nil {
		return errResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// GetOrgStats handles GET /organizations/:id/stats.
// Query params: start_date, end_date (YYYY-MM-DD), group_by (hour|day|week|month).
func (h *OrgHandler) GetOrgStats(ctx *fiber.Ctx) error {
	orgId, err := parseOrgId(ctx)
	if err != nil {
		return nil
	}
	resp, err := h.OrgService.GetOrgStats(
		ctx.Context(), orgId,
		ctx.Query("start_date"),
		ctx.Query("end_date"),
		ctx.Query("group_by", "day"),
	)
	if err != nil {
		return errResp(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
