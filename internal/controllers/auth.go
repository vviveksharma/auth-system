package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vviveksharma/auth/internal/models"
	dbmodels "github.com/vviveksharma/auth/models"
)

// LoginUser handles user login requests.
//
// @Summary      User Login
// @Description  Authenticates a user and returns a JWT token upon successful login.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.UserLoginRequest  true  "User login credentials"
// @Success      200   {object}  dbmodels.ServiceResponse "JWT token and success message"
// @Failure      400   {object}  dbmodels.ServiceResponse "Invalid credentials or missing fields"
// @Failure      502   {object}  dbmodels.ServiceResponse "Error while parsing the request body"
// @Failure      500   {object}  dbmodels.ServiceResponse "Unexpected server error"
// @Router       /login [post]
func (h *Handler) LoginUser(ctx *fiber.Ctx) error {
	var req models.UserLoginRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("Error in parsing the request Body" + err.Error())
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "error while parsing the requestBody: " + err.Error(),
		}
	}
	if req.Email == "" || req.Password == "" {
		return &dbmodels.ServiceResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Please check your credentials",
		}
	}
	if req.Role == "" {
		return &dbmodels.ServiceResponse{
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

// RefreshToken refreshes the JWT token for an authenticated user.
//
// @Summary      Refresh JWT Token
// @Description  Refreshes and returns a new JWT token for the authenticated user.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dbmodels.ServiceResponse "Refreshed JWT token and success message"
// @Failure      500  {object}  dbmodels.ServiceResponse "Unexpected server error"
// @Router       /refresh-token [post]
func (h *Handler) RefreshToken(ctx *fiber.Ctx) error {
	claims := ctx.Locals("authClaims").(jwt.MapClaims)
	userId := claims["user_id"]
	roleId := claims["role_id"]
	resp, err := h.AuthService.RefreshToken(userId.(string), roleId.(string))
	if err != nil {
		if serviceErr, ok := err.(*dbmodels.ServiceResponse); ok {
			return ctx.Status(serviceErr.Code).JSON(err)
		} else {
			return ctx.JSON(500, "an unexpected error occurred")
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(dbmodels.ServiceResponse{
		Code:    200,
		Message: "The JWT token is refreshed",
		Data:    resp,
	})
}
