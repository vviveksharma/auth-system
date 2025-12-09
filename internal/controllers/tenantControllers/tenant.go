package tenantcontrollers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/internal/models"
	tenantmodels "github.com/vviveksharma/auth/internal/models/tenantModels"
	"github.com/vviveksharma/auth/internal/pagination"
	dbmodels "github.com/vviveksharma/auth/models"
)

func (h *TenantHandler) CreateTenant(ctx *fiber.Ctx) error {
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
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while creating tenant: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    resp,
	})
}

func (h *TenantHandler) LoginTenant(ctx *fiber.Ctx) error {
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
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "",
		Data:    resp,
	})
}

func (h *TenantHandler) ListTokens(ctx *fiber.Ctx) error {
	resp, err := h.TenantService.ListTokens(ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	// Make the tokens response paginated
	page := ctx.Query("page")
	page_size := ctx.Query("page_size")
	pageInt := 1
	pageSizeInt := 5
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageInt = p
		}
	}
	if page_size != "" {
		if ps, err := strconv.Atoi(page_size); err == nil && ps > 0 && ps <= 100 {
			pageSizeInt = ps
		}
	}
	paginatedResponse := pagination.PaginateSlice(resp, pageInt, pageSizeInt)
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "Tokens retrieved successfully.",
		Data:    paginatedResponse,
	})
}

func (h *TenantHandler) RevokeToken(ctx *fiber.Ctx) error {
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
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

func (h *TenantHandler) CreateToken(ctx *fiber.Ctx) error {
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
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

func (h *TenantHandler) ResetPassword(ctx *fiber.Ctx) error {
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
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

func (h *TenantHandler) SetPassword(ctx *fiber.Ctx) error {
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
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

func (h *TenantHandler) GetTenantDetails(ctx *fiber.Ctx) error {
	resp, err := h.TenantService.GetTenantDetails(ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "The Tenant details are as follows",
		Data:    resp,
	})
}

func (h *TenantHandler) DeleteTenant(ctx *fiber.Ctx) error {
	resp, err := h.TenantService.DeleteTenant(ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
		Data:    nil,
	})
}

func (h *TenantHandler) GetDashboardDetails(ctx *fiber.Ctx) error {
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

func (h *TenantHandler) GetTokenDetailsStatus(ctx *fiber.Ctx) error {
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
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageInt = p
		}
	}
	if page_size != "" {
		if ps, err := strconv.Atoi(page_size); err == nil && ps > 0 && ps <= 100 {
			pageSizeInt = ps
		}
	}
	resp, err := h.TenantService.ListTokensWithStatus(ctx.Context(), pageInt, pageSizeInt, flag)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Tokens fetched successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) ListRoles(ctx *fiber.Ctx) error {
	roleType := ctx.Query("roletype")
	if roleType != "" {
		if roleType != "system" && roleType != "custom" {
			return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid request payload. The roletype could be either `system` or `custom`.",
			})
		}
	} else {
		roleType = "all"
	}
	if roleType == "system" {
		roleType = "default"
	}
	status := ctx.Query("status")
	if status != "" {
		if status != "active" && status != "inactive" {
			return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid request payload. The status could be either `active` or `inactive`.",
			})
		}
	} else {
		status = "all"
	}
	page := ctx.Query("page")
	page_size := ctx.Query("page_size")
	pageInt := 1
	pageSizeInt := 5
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageInt = p
		}
	}
	if page_size != "" {
		if ps, err := strconv.Atoi(page_size); err == nil && ps > 0 && ps <= 100 {
			pageSizeInt = ps
		}
	}
	resp, err := h.TenantRoleService.TenantListRoles(ctx.Context(), pageInt, pageSizeInt, roleType, status)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Roles fetched successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) AddRole(ctx *fiber.Ctx) error {
	var req tenantmodels.TenantAddRoleRequestBody
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	if req.Name == "" || req.DisplayName == "" || req.Permissions == nil {
		log.Println("the requestBody: ", req)
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Missing required fields: name, displayname, or permissions. Please ensure all fields are provided.",
		}
	}
	resp, err := h.TenantRoleService.TenantAddRole(ctx.Context(), &req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Roles fetched successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) GetRolePermissions(ctx *fiber.Ctx) error {
	roleId := ctx.Query("roleId")
	if roleId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. The roletype could be either `system` or `custom`.",
		})
	}
	resp, err := h.TenantRoleService.TenantGetRolePermissions(ctx.Context(), uuid.MustParse(roleId))
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Roles permissions fetched successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) EnableRole(ctx *fiber.Ctx) error {
	roleId := ctx.Query("roleId")
	if roleId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. The roletype could be either `system` or `custom`.",
		})
	}
	resp, err := h.TenantRoleService.TenantEnableRole(ctx.Context(), uuid.MustParse(roleId))
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Role enabled successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) DisableRole(ctx *fiber.Ctx) error {
	roleId := ctx.Query("roleId")
	if roleId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. The roletype could be either `system` or `custom`.",
		})
	}
	resp, err := h.TenantRoleService.TenantDisableRole(ctx.Context(), uuid.MustParse(roleId))
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Role diabled successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) DeleteRole(ctx *fiber.Ctx) error {
	roleId := ctx.Query("roleId")
	if roleId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. The roletype could be either `system` or `custom`.",
		})
	}
	resp, err := h.TenantRoleService.TenantDeleteRole(ctx.Context(), uuid.MustParse(roleId))
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.Status(500).JSON(dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An unexpected error occurred while deleting user: %v", err),
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Role deleted successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) EditRolePermissions(ctx *fiber.Ctx) error {
	var req tenantmodels.TeanantEditPermissionRequestBody
	roleId := ctx.Params("id")
	if roleId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "roleId in the params is requried",
		})
	}
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	if req.UpdateRoleDetails && req.RoleInfo == (models.RoleInfo{}) {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "RoleInfo is required when UpdateRoleDetails is true.",
		})
	}

	if req.UpdateRolePermissions && req.Permissions == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Permissions is required when UpdateRolePermissions is true.",
		})
	}

	resp, err := h.TenantRoleService.TenantEditRolePermissions(ctx.Context(), uuid.MustParse(roleId), &req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while editing role permissions: %v", err))
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Role permissions updated successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) ListMessages(ctx *fiber.Ctx) error {
	log.Println("Inside the list messages handlers")
	ctx.Context().Value("tenant")
	status := ctx.Query("status")
	if status != "" {
		if status != "pending" && status != "approved" && status != "rejected" {
			return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid request payload. The status could be either `active` or `inactive`.",
			})
		}
	} else {
		status = "all"
	}
	page := ctx.Query("page")
	page_size := ctx.Query("page_size")
	pageInt := 1
	pageSizeInt := 5
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageInt = p
		}
	}
	if page_size != "" {
		if ps, err := strconv.Atoi(page_size); err == nil && ps > 0 && ps <= 100 {
			pageSizeInt = ps
		}
	}
	resp, err := h.TenantMessageService.ListMessageRequest(ctx.Context(), pageInt, pageSizeInt, status)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while editing role permissions: %v", err))
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "List Message successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) ApproveMessage(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. The messageId is the required parameter",
		})
	}
	resp, err := h.TenantMessageService.ApproveRequest(ctx.Context(), uuid.MustParse(id))
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while editing role permissions: %v", err))
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Message satus updated  successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) RejectMessage(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. The messageId is the required parameter",
		})
	}
	resp, err := h.TenantMessageService.RejectRequest(ctx.Context(), uuid.MustParse(id))
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while editing role permissions: %v", err))
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "Message satus updated successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) ListUsers(ctx *fiber.Ctx) error {
	status := ctx.Query("status")
	if status != "" {
		if status != "active" && status != "inactive" {
			return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid request payload. The status could be either `active` or `inactive`.",
			})
		}
	} else {
		status = "all"
	}
	page := ctx.Query("page")
	page_size := ctx.Query("page_size")
	pageInt := 1
	pageSizeInt := 5
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageInt = p
		}
	}
	if page_size != "" {
		if ps, err := strconv.Atoi(page_size); err == nil && ps > 0 && ps <= 100 {
			pageSizeInt = ps
		}
	}
	resp, err := h.TenantUserService.TenantListUsers(ctx.Context(), pageInt, pageSizeInt, status)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while editing role permissions: %v", err))
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "User listed successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) EnableUser(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. The userId is the required parameter",
		})
	}
	resp, err := h.TenantUserService.TenantEnableUser(ctx.Context(), uuid.MustParse(id))
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while editing role permissions: %v", err))
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "User status updated successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) DisableUser(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. The userId is the required parameter",
		})
	}
	resp, err := h.TenantUserService.TenantDisableUser(ctx.Context(), uuid.MustParse(id))
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while editing role permissions: %v", err))
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "User status updated successfully",
		Data:    resp,
	})
}

func (h *TenantHandler) DeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request payload. The userId is the required parameter",
		})
	}
	resp, err := h.TenantUserService.TenantDeleteUser(ctx.Context(), uuid.MustParse(id))
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, fmt.Sprintf("An unexpected error occurred while editing role permissions: %v", err))
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(&dbmodels.ServiceResponse{
		Code:    200,
		Message: "User status updated successfully",
		Data:    resp,
	})
}
