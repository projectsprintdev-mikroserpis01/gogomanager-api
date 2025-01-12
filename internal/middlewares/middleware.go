package middlewares

import "github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/jwt"

type Middleware struct {
	jwt jwt.JwtInterface
	jwtManager jwt.JwtManagerInterface
}

func NewMiddleware(
	jwt jwt.JwtInterface,
	jwtManager jwt.JwtManagerInterface,
) *Middleware {
	return &Middleware{
		jwt: jwt,
		jwtManager: jwtManager,
	}
}
