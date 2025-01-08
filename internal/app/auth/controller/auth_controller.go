package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/helpers/http/response"
)

type authController struct {
	authService contracts.AuthService
}

func InitAuthController(router fiber.Router, authService contracts.AuthService) {
	controller := authController{
		authService: authService,
	}

	authGroup := router.Group("/auth")
	authGroup.Post("/register", controller.registerUser)
	authGroup.Post("/login", controller.loginUser)
}

func (ac *authController) registerUser(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	res, err := ac.authService.RegisterUser(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (ac *authController) loginUser(ctx *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	res, err := ac.authService.LoginUser(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}
