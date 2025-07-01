package http

import (
	"Medods/internal/auth"

	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers) {
	authGroup.POST("/tokens/:user_id", h.CreateTokens())
	authGroup.POST("/tokens/refresh", h.RefreshTokens())
	authGroup.GET("/tokens", h.GetUserIDByToken())
	authGroup.POST("/tokens/logout", h.UserLogout())
}
