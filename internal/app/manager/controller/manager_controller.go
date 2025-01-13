package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/manager/service"
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

	jwt := jwt.Jwt

	middleware := middlewares.NewMiddleware(jwt)

	managerRoute := router.Group("/user")
	managerRoute.Get("/", middleware.RequireAuth(), controller.GetManagerById)
	managerRoute.Patch("/", middleware.RequireAuth(), controller.UpdateManagerById)
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
	var requestBody dto.GetCurrentManagerRequest
	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	token := ctx.Get("Authorization")
	var claims jwt.Claims
	if err := jwt.Jwt.Decode(token, &claims); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// TODO:
	var id int
	res, err := mc.managerService.GetManagerById(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
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

	token := ctx.Get("Authorization")
	var claims jwt.Claims
	if err := jwt.Jwt.Decode(token, &claims); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// TODO:
	var id int

	res, err := mc.managerService.UpdateManagerById(ctx.Context(), id, requestBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}
