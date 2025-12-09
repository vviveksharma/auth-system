package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/internal/models"
	responsemodels "github.com/vviveksharma/auth/models"
)

func (h *Handler) CreateRequest(ctx *fiber.Ctx) error {
	var req *models.CreateMessageRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return UnprocessableEntity(ctx)
	}
	if req.Email == "" || req.RequestedRole == "" {
		return BadRequest(ctx, "Invalid request: 'roleName' and 'email' fields are required and cannot be empty.")
	}
	resp, err := h.MessageService.CreateMessage(req, ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return InternalServerError(ctx, "Internal server error occurred: "+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}

func (h *Handler) GetRequestStatus(ctx *fiber.Ctx) error {
	messageId := ctx.Query("id")
	if messageId == "" {
		return BadRequest(ctx, "the id in the query paramter is a required field")
	}
	resp, err := h.MessageService.GetStatus(uuid.MustParse(messageId), ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return InternalServerError(ctx, "Internal server error occurred: "+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: resp.Status,
	})
}

func (h *Handler) GetMessages(ctx *fiber.Ctx) error {
	email := ctx.Query("email")
	if email == "" {
		return BadRequest(ctx, "Invalid request: email' fields are required and cannot be empty.")
	}
	resp, err := h.MessageService.ListMessages(email, ctx.Context())
	if err != nil {
		if serviceErr, ok := err.(*responsemodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return InternalServerError(ctx, "Internal server error occurred: "+err.Error())
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(responsemodels.ServiceResponse{
		Code:    200,
		Message: "The requested roles are as follow",
		Data:    resp,
	})
}
