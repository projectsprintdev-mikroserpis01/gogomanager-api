package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain"
)

type Response struct {
	Payload interface{} `json:"payload"`
}

func SendResponse(
	ctx *fiber.Ctx,
	code int,
	payload interface{},
) error {
	if code >= 400 {
		if err, ok := payload.(error); ok {
			var errPayload any = err
			if _, ok := err.(domain.SerializableError); !ok {
				errPayload = err.Error()
			}
			payload = fiber.Map{"error": errPayload}
		}
	}

	return ctx.Status(code).JSON(
		Response{
			Payload: payload,
		},
	)
}
