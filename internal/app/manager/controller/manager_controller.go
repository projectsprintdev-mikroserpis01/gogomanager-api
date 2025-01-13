package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/manager/service"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/middlewares"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/jwt"
)

type managerController struct {
	managerService service.ManagerService
}

func InitManagerController(router fiber.Router, managerService service.ManagerService) {
	controller := managerController{
		managerService: managerService,
	}

	authGroup := router.Group("/v1/auth")
	authGroup.Post("/", controller.handleAuth)

	jwtManager := jwt.JwtManager
	jwt := jwt.Jwt

	middleware := middlewares.NewMiddleware(jwt, jwtManager)

	managerRoute := router.Group("/v1/user")
	managerRoute.Get("/", middleware.RequireAdmin(), controller.GetManagerById)
	managerRoute.Patch("/", middleware.RequireAdmin(), controller.UpdateManagerById)
}

func (mc *managerController) handleAuth(ctx *fiber.Ctx) error {
	var req dto.AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := mc.managerService.Authenticate(ctx.Context(), req)
	if err != nil {
		return err
	}

	status := fiber.StatusOK
	if req.Action == "create" {
		status = fiber.StatusCreated
	}
	return ctx.Status(status).JSON(res)

}

func (mc *managerController) GetManagerById(ctx *fiber.Ctx) error {
	managerID := ctx.Locals("claims").(jwt.ClaimsManager).UserID

	res, err := mc.managerService.GetManagerById(ctx.Context(), managerID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (mc *managerController) UpdateManagerById(ctx *fiber.Ctx) error {
	var requestBody dto.UpdateManagerRequest
	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	managerID := ctx.Locals("claims").(jwt.ClaimsManager).UserID

	_, err := mc.managerService.UpdateManagerById(ctx.Context(), managerID, requestBody)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(requestBody)
}
