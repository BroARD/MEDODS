package http

import (
	"Medods/internal/auth"
	"Medods/internal/middleware"

	"github.com/labstack/echo/v4"
)


func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/tokens/:user_id", h.CreateTokens())
	authGroup.POST("/refresh", h.RefreshTokens(), mw.AuthRefreshMiddleware)
	authGroup.GET("/tokens", h.GetUserIDByToken(), mw.AuthJWTMiddleware)
	authGroup.POST("/logout", h.UserLogout(), mw.AuthJWTMiddleware)
}
