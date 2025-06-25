package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vviveksharma/auth/internal/models"
	dbmodels "github.com/vviveksharma/auth/models"
)

func (h *Handler) CreateUser(ctx *fiber.Ctx) error {
	var req *models.UserRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &fiber.Error{
			Code:    fiber.StatusBadGateway,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	if req.Email == "" || req.Name == "" || req.Password == "" {
		log.Println("the requestBody: ", req)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "either name or type is missing in the request Body",
		}
	}
	fmt.Println("the userdetails from the request ", req)
	resp, err := h.UserService.CreateUser(req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}

func (h *Handler) GetUserDetails(ctx *fiber.Ctx) error {
	req := &models.GetUserDetailsRequest{}
	userId := ctx.Locals("userId").(string)
	fmt.Println("the userid: ", userId)
	req.Id = userId
	resp, err := h.UserService.GetUserDetails(req)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "",
		Data:    resp,
	})
}

func (h *Handler) UpdateUserDetails(ctx *fiber.Ctx) error {
	req := &models.UpdateUserRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &fiber.Error{
			Code:    fiber.StatusBadGateway,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	userId := ctx.Locals("userId").(string)
	fmt.Println("the userid: ", userId)
	resp, err := h.UserService.UpdateUserDetails(req, userId)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}

func (h *Handler) GetUserByIdDetails(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	resp, err := h.UserService.GetUserById(userId)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "",
		Data:    resp,
	})
}

func (h *Handler) AssignUserRole(ctx *fiber.Ctx) error {
	var req *models.AssignRoleRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &fiber.Error{
			Code:    fiber.StatusBadGateway,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	userId := ctx.Params("id")
	resp, err := h.UserService.AssignUserRole(req, userId)
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: resp.Message,
	})
}
