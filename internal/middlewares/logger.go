package middlewares

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/log"
)

func LoggerConfig() fiber.Handler {
	logger := log.GetLogger()
	config := fiberzerolog.Config{
		Logger:          logger,
		FieldsSnakeCase: true,
		Fields: []string{
			"referer",
			"ip",
			"host",
			"url",
			"ua",
			"latency",
			"status",
			"method",
		},
		Messages: []string{
			"[LoggerMiddleware.LoggerConfig] Server error",
			"[LoggerMiddleware.LoggerConfig] Client error",
			"[LoggerMiddleware.LoggerConfig] Success",
		},
	}

	return fiberzerolog.New(config)
}
