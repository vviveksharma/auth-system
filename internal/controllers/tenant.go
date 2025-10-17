package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/pagination"
	dbmodels "github.com/vviveksharma/auth/models"
)

func (h *Handler) CreateTenant(ctx *fiber.Ctx) error {
	var req models.CreateTenantRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	if req.Email == "" || req.Name == "" || req.Password == "" || req.Campany == "" {
		log.Println("the requestBody: ", req)
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Missing required fields: name, email, password, or company. Please ensure all fields are provided.",
		}
	}
	resp, err := h.TenantService.CreateTenant(&req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while deleting user: %v", err)+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}

func (h *Handler) LoginTenant(ctx *fiber.Ctx) error {
	var req models.LoginTenantRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Printf("Failed to parse login request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. Please ensure the request body is properly formatted.",
		})
	}
	if req.Email == "" || req.Password == "" {
		log.Printf("Login attempt with missing credentials: %+v", req)
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Email and password are required fields. Please provide both to proceed.",
		})
	}
	resp, err := h.TenantService.LoginTenant(&req, ctx.IP())
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while deleting user: %v", err)+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "",
		Data:    resp,
	})
}

func (h *Handler) ListTokens(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(string)
	resp, err := h.TenantService.ListTokens(ctx.Context(), token)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while deleting user: %v", err)+err.Error())
		}
	}
	// Make the tokens response paginated
	page := ctx.Query("page")
	page_size := ctx.Query("page_size")
	pageInt := 1
	pageSizeInt := 5
	if page != "" {
		pageInt, _ = strconv.Atoi(page)
	}
	if page_size != "" {
		pageSizeInt, _ = strconv.Atoi(page_size)
	}
	paginatedResponse := pagination.PaginateSlice(resp, pageInt, pageSizeInt)
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "Tokens retrieved successfully.",
		Data:    paginatedResponse,
	})
}

func (h *Handler) RevokeToken(ctx *fiber.Ctx) error {
	token := ctx.Params("id")
	if token == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    400,
			Message: "token that needed to be revoked should not be empty",
		})
	}
	resp, err := h.TenantService.RevokeToken(ctx.Context(), token)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while deleting user: %v", err)+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

func (h *Handler) CreateToken(ctx *fiber.Ctx) error {
	var req models.CreateTokenRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Printf("Failed to parse login request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. Please ensure the request body is properly formatted.",
		})
	}
	if req.ExpiryAt == "" || req.Name == "" {
		log.Printf("Create token attempt failed: %+v", req)
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Expiry at and name are required fields. Please provide both to proceed.",
		})
	}
	resp, err := h.TenantService.CreateToken(ctx.Context(), &req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while deleting user: %v", err)+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

func (h *Handler) ResetPassword(ctx *fiber.Ctx) error {
	var req models.ResetTenantPasswordRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Printf("Failed to parse login request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. Please ensure the request body is properly formatted.",
		})
	}
	if req.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Email are required fields. Please provide both to proceed.",
		})
	}
	resp, err := h.TenantService.ResetPassword(ctx.Context(), &req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while deleting user: %v", err)+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

func (h *Handler) SetPassword(ctx *fiber.Ctx) error {
	var req models.SetTenantPasswordRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Printf("Failed to parse login request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. Please ensure the request body is properly formatted.",
		})
	}
	if req.NewPassword == "" || req.ConfirmNewPassword == "" {
		log.Printf("Create token attempt failed: %+v", req)
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "NewPassword and ConfirmNewPassword are required fields. Please provide both to proceed.",
		})
	}
	resp, err := h.TenantService.SetPassword(ctx.Context(), &req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while deleting user: %v", err)+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

func (h *Handler) GetTenantDetails(ctx *fiber.Ctx) error {
	resp, err := h.TenantService.GetTenantDetails(ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while deleting user: %v", err)+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "The Tenant details are as follows",
		Data:    resp,
	})
}

func (h *Handler) DeleteTenant(ctx *fiber.Ctx) error {
	resp, err := h.TenantService.DeleteTenant(ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while deleting user: %v", err)+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

func (h *Handler) GetDashboardDetails(ctx *fiber.Ctx) error {
	resp, err := h.TenantService.GetDashboardDetails(ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while getting the tenant: %v", err)+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "The Tenant details",
		Data:    resp,
	})
}

func (h *Handler) GetTokenDetailsStatus(ctx *fiber.Ctx) error {
	flag := ctx.Query("status")
	if flag != "active" && flag != "revoked" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. The status could be either `active` or `revoked`.",
		})
	}
	page := ctx.Query("page")
	page_size := ctx.Query("page_size")
	pageInt := 1
	pageSizeInt := 5
	if page != "" {
		pageInt, _ = strconv.Atoi(page)
	}
	if page_size != "" {
		pageSizeInt, _ = strconv.Atoi(page_size)
	}
	resp, err := h.TenantService.ListTokensWithStatus(ctx.Context(), pageInt, pageSizeInt, flag)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while deleting user: %v", err)+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Roles fetched successfully",
		Data:    resp,
	})
}

func (h *Handler) ListUserTenant(ctx *fiber.Ctx) error {
	status := ctx.Query("status")

	if status == "" {
		status = "all"
	} else {
		if status != "active" && status != "inactive" && status != "all" {
			return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid request in the query parameter. The status must be 'active', 'inactive', or 'all'.",
			})
		}
	}

	// Parse pagination parameters
	page := ctx.Query("page")
	page_size := ctx.Query("page_size")
	pageInt := 1
	pageSizeInt := 5
	if page != "" {
		pageInt, _ = strconv.Atoi(page)
	}
	if page_size != "" {
		pageSizeInt, _ = strconv.Atoi(page_size)
	}

	resp, err := h.TenantService.ListUsers(ctx.Context(), pageInt, pageSizeInt, status)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while deleting user: %v", err)+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "User fetched successfully",
		Data:    resp,
	})
}
