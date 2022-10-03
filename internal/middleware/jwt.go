package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"entetry/gotest/internal/model"
)

func NewJwtMiddleware(accessTokenKey string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(accessTokenKey),
		Claims:     new(model.Claim),
	})
}
