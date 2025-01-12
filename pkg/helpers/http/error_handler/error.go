package errorhandler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/helpers/http/response"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/validator"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		return response.SendResponse(c, fiber.StatusBadRequest, valErr)
	}

	var reqErr *domain.RequestError
	if errors.As(err, &reqErr) {
		return response.SendResponse(c, reqErr.StatusCode, reqErr)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return response.SendResponse(c, fiberErr.Code, fiberErr)
	}

	return response.SendResponse(c, fiber.StatusInternalServerError, err)
}
