package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/middlewares"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/jwt"
)

type departmentController struct {
	service    contracts.DepartmentService
	middleware *middlewares.Middleware
}

func InitNewController(
	router fiber.Router,
	departmentService contracts.DepartmentService,
	middleware *middlewares.Middleware,
) {
	controller := &departmentController{
		service:    departmentService,
		middleware: middleware,
	}

	route := router.Group("/v1")

	route.Post("/department", middleware.RequireAdmin(), controller.Create)
	route.Get("/department", middleware.RequireAdmin(), controller.Get)
	route.Patch("/department/:departmentid", middleware.RequireAdmin(), controller.Update)
	route.Delete("/department/:departmentid", middleware.RequireAdmin(), controller.Delete)
}

func (c *departmentController) Create(ctx *fiber.Ctx) error {
	var requestBody struct {
		ManagerID int    `json:"managerId"`
		Name      string `json:"name"`
	} // manager id need to get from token

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	managerID := ctx.Locals("claims").(jwt.ClaimsManager).UserID
	requestBody.ManagerID = managerID

	departmentRes, err := c.service.Create(ctx.Context(), requestBody.ManagerID, requestBody.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(departmentRes)
}
func (c *departmentController) Update(ctx *fiber.Ctx) error {
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	departmentID := ctx.Params("departmentid")
	id, err := strconv.Atoi(departmentID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid department ID",
		})
	}

	departmentRes, err := c.service.Update(ctx.Context(), id, requestBody.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(departmentRes)
}
func (c *departmentController) Get(ctx *fiber.Ctx) error {
	limit, err := strconv.Atoi(ctx.Query("limit", "0"))
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

	name := ctx.Query("name", "")

	var departmentRes []dto.DepartmentRes

	if limit == 0 && offset == 0 && name == "" {
		departments, err := c.service.FindAll(ctx.Context(), 0, 0)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		for _, dept := range departments {
			departmentRes = append(departmentRes, dto.DepartmentRes{
				ID:   dept.ID,
				Name: dept.Name,
			})
		}
	} else if name != "" {
		departments, err := c.service.FindByName(ctx.Context(), limit, offset, name)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		for _, dept := range departments {
			departmentRes = append(departmentRes, dto.DepartmentRes{
				ID:   dept.ID,
				Name: dept.Name,
			})
		}
	} else if limit > 0 && offset > 0 {
		departments, err := c.service.FindAll(ctx.Context(), limit, offset)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		for _, dept := range departments {
			departmentRes = append(departmentRes, dto.DepartmentRes{
				ID:   dept.ID,
				Name: dept.Name,
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"departments": departmentRes,
	})
}

func (c *departmentController) Delete(ctx *fiber.Ctx) error {
	departmentID := ctx.Params("departmentid")
	id, err := strconv.Atoi(departmentID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid department ID",
		})
	}

	err = c.service.Delete(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Department deleted successfully",
	})
}
