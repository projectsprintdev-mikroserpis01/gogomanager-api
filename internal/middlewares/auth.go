package middlewares

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/enums"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/jwt"
)

func (m *Middleware) RequireAuth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		header := ctx.Get("Authorization")
		if header == "" {
			return domain.ErrNoBearerToken
		}

		headerSlice := strings.Split(header, " ")
		if len(headerSlice) != 2 && headerSlice[0] != "Bearer" {
			return domain.ErrInvalidBearerToken
		}

		token := headerSlice[1]
		var claims jwt.Claims
		err := m.jwt.Decode(token, &claims)
		if err != nil {
			return domain.ErrInvalidBearerToken
		}

		notBefore, err := claims.GetNotBefore()
		if err != nil {
			return domain.ErrInvalidBearerToken
		}

		if notBefore.After(time.Now()) {
			return domain.ErrBearerTokenNotActive
		}

		expirationTime, err := claims.GetExpirationTime()
		if err != nil {
			return domain.ErrInvalidBearerToken
		}

		if expirationTime.Before(time.Now()) {
			return domain.ErrExpiredBearerToken
		}

		ctx.Locals("claims", claims)

		return ctx.Next()
	}
}

func (m *Middleware) RequireOneOfRoles(roles ...enums.RoleEnum) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims, ok := ctx.Locals("claims").(jwt.Claims)
		if !ok {
			return domain.ErrInvalidBearerToken
		}

		if claims.RoleName == enums.SuperAdmin.String() {
			return ctx.Next()
		}

		for _, role := range roles {
			if claims.RoleName == role.String() {
				return ctx.Next()
			}
		}

		return domain.ErrRoleCantAccessResource
	}
}

func (m *Middleware) RequireAdmin() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		header := ctx.Get("Authorization")
		if header == "" {
			return domain.ErrNoBearerToken
		}

		headerSlice := strings.Split(header, " ")
		if len(headerSlice) != 2 && headerSlice[0] != "Bearer" {
			return domain.ErrInvalidBearerToken
		}

		token := headerSlice[1]
		var claims jwt.ClaimsManager
		err := m.jwtManager.DecodeManager(token, &claims)
		if err != nil {
			return domain.ErrInvalidBearerToken
		}

		notBefore, err := claims.GetNotBefore()
		if err != nil {
			return domain.ErrInvalidBearerToken
		}

		if notBefore.After(time.Now()) {
			return domain.ErrBearerTokenNotActive
		}

		expirationTime, err := claims.GetExpirationTime()
		if err != nil {
			return domain.ErrInvalidBearerToken
		}

		if expirationTime.Before(time.Now()) {
			return domain.ErrExpiredBearerToken
		}

		ctx.Locals("claims", claims)

		return ctx.Next()
	}
}
