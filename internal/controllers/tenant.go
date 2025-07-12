package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/models"
	dbmodels "github.com/vviveksharma/auth/models"
)

func (h *Handler) CreateTenant(ctx *fiber.Ctx) error {
	var req models.CreateTenantRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusBadGateway,
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
			return ctx.JSON(500, "an unexpected error occurred"+err.Error())
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
	resp, err := h.TenantService.LoginTenant(&req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred"+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "",
		Data:    resp,
	})
}
