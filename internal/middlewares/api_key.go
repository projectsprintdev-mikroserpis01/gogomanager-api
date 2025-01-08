package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/infra/env"
)

func ApiKey() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		apiKey := ctx.Get("x-api-key")
		if apiKey == "" {
			return domain.ErrNoAPIKey
		}

		keySlice := strings.Split(apiKey, " ")
		if len(keySlice) != 2 {
			return domain.ErrInvalidAPIKey
		}

		key := keySlice[1]
		if key != env.AppEnv.ApiKey {
			return domain.ErrInvalidAPIKey
		}

		return ctx.Next()
	}
}
