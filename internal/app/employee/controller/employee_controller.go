package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/middlewares"
)

type employeeController struct {
	employeeService contracts.EmployeeService
	middleware      *middlewares.Middleware
}

func InitNewController(
	router fiber.Router,
	employeeService contracts.EmployeeService,
	middleware *middlewares.Middleware,
) {
	controller := &employeeController{
		employeeService: employeeService,
		middleware:      middleware,
	}

	route := router.Group("/v1/employee")

	route.Post("/", middleware.RequireAdmin(), controller.Create)
	route.Get("/", middleware.RequireAdmin(), controller.Get)
	route.Patch("/:identityNumber", middleware.RequireAdmin(), controller.Update)
	route.Patch("/", middleware.RequireAdmin(), func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Identity Number is required",
		})
	})
	route.Delete("/:identityNumber", middleware.RequireAdmin(), controller.Delete)
	route.Delete("/", middleware.RequireAdmin(), func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Identity Number is required",
		})
	})
}

func (c *employeeController) Create(ctx *fiber.Ctx) error {
	var req dto.EmployeeCreateReq
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := c.employeeService.Create(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}

func (c *employeeController) Get(ctx *fiber.Ctx) error {
	identityNumber := ctx.Query("identityNumber", "")
	name := ctx.Query("name", "")
	gender := ctx.Query("gender", "")
	departmentID, err := strconv.Atoi(ctx.Query("departmentId", "0"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid department ID query parameter",
		})
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "5"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit query parameter",
		})
	}

	offset, err := strconv.Atoi(ctx.Query("offset", "0"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid offset query parameter",
		})
	}

	res, err := c.employeeService.Find(ctx.Context(), identityNumber, name, gender, departmentID, limit, offset)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *employeeController) Update(ctx *fiber.Ctx) error {
	var req dto.EmployeeUpdateReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	identityNumber := ctx.Params("identityNumber")
	if identityNumber == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invalid identity number",
		})
	}

	res, err := c.employeeService.Update(ctx.Context(), req, identityNumber)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *employeeController) Delete(ctx *fiber.Ctx) error {
	identityNumber := ctx.Params("identityNumber")
	if identityNumber == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invalid identity number",
		})
	}

	err := c.employeeService.Delete(ctx.Context(), identityNumber)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Employee deleted successfully",
	})
}
