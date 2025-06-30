package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/models"
	dbmodels "github.com/vviveksharma/auth/models"
)

func (h *Handler) LoginUser(ctx *fiber.Ctx) error {
	var req models.UserLoginRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &fiber.Error{
			Code:    fiber.StatusBadGateway,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	if req.Email == "" || req.Password == "" {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Please check your credentials",
		}
	}
	if req.Role == "" {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "role should be provided",
		}
	}
	resp, err := h.AuthService.LoginUser(&req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "The JWT token is as follow",
		Data:    resp,
	})
}
