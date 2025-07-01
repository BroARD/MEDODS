package auth

import "github.com/labstack/echo/v4"

type Handlers interface {
	CreateTokens() echo.HandlerFunc
	RefreshTokens() echo.HandlerFunc
	GetUserIDByToken() echo.HandlerFunc
	UserLogout() echo.HandlerFunc
}
