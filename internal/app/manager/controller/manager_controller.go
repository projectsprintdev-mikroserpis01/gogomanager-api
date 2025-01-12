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
