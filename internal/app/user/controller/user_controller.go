package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/helpers/http/response"
)

type userController struct {
	userService contracts.UserService
}

func InitNewController(router fiber.Router, userService contracts.UserService) {
	controller := &userController{
		userService: userService,
	}

	userRoute := router.Group("/users")
	userRoute.Get("/", controller.getUsers)
	userRoute.Get("/stats", controller.getUsersStats)
	userRoute.Get("/:id", controller.getUserByID)
	userRoute.Post("/", controller.createUser)
	userRoute.Put("/:id", controller.updateUser)
	userRoute.Patch("/:id/soft-delete", controller.softDeleteUser)
	userRoute.Patch("/:id/restore", controller.restoreUser)
	userRoute.Delete("/:id", controller.deleteUser)
}

func (uc *userController) getUsers(ctx *fiber.Ctx) error {
	var query dto.GetUsersQuery
	if err := ctx.QueryParser(&query); err != nil {
		return err
	}

	res, err := uc.userService.GetUsers(ctx.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (uc *userController) getUsersStats(ctx *fiber.Ctx) error {
	res, err := uc.userService.GetUsersStats(ctx.Context())
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (uc *userController) getUserByID(ctx *fiber.Ctx) error {
	var req dto.GetUserByIDRequest
	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	res, err := uc.userService.GetUserByID(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (uc *userController) createUser(ctx *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	res, err := uc.userService.CreateUser(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, res)
}

func (uc *userController) updateUser(ctx *fiber.Ctx) error {
	var req dto.UpdateUserRequest
	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	res, err := uc.userService.UpdateUser(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (uc *userController) softDeleteUser(ctx *fiber.Ctx) error {
	var req dto.SoftDeleteUserRequest
	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	res, err := uc.userService.SoftDeleteUser(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (uc *userController) restoreUser(ctx *fiber.Ctx) error {
	var req dto.RestoreUserRequest
	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	res, err := uc.userService.RestoreUser(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (uc *userController) deleteUser(ctx *fiber.Ctx) error {
	var req dto.DeleteUserRequest
	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	res, err := uc.userService.DeleteUser(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}
