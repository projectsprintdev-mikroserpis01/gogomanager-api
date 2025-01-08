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
		return response.SendResponse(c, fiber.StatusUnprocessableEntity, valErr)
	}

	var reqErr *domain.RequestError
	if errors.As(err, &reqErr) {
		return response.SendResponse(c, reqErr.StatusCode, reqErr)
	}

	return response.SendResponse(c, fiber.StatusInternalServerError, err)
}
