package middleware

import (
	"entetry/gotest/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// todo implement fingertprint check
func NewJwtMiddleware(accessTokenKey string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(accessTokenKey),
		Claims:     new(model.Claim),
	})
}
