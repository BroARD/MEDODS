package http

import (
	"Medods/internal/auth"
	"Medods/pkg/logging"
	"net/http"

	"github.com/labstack/echo/v4"
)

type authHandlers struct {
	authUC auth.UseCase
	logger logging.Logger
}

func (h *authHandlers) CreateTokens() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user_id := ctx.Param("user_id")
		userAgent := ctx.Request().UserAgent()
		userIP := ctx.RealIP()
		tokensPair, err := h.authUC.CreateTokens(ctx.Request().Context(), user_id, userAgent, userIP)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, "Could not create pair of tokens")
		}
		return ctx.JSON(http.StatusCreated, tokensPair)
	}
}

func (h *authHandlers) GetUserIDByToken() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get("user_id").(string)
		return ctx.JSON(http.StatusOK, userID)
	}
}

func (h *authHandlers) RefreshTokens() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userAgent := ctx.Get("userAgent").(string)
		userIP := ctx.Get("userIP").(string)
		refToken := ctx.Get("refToken").(string)
		tokenString := ctx.Get("user_id").(string)
		pairTokens, err := h.authUC.RefreshToken(ctx.Request().Context(), refToken, tokenString, userAgent, userIP)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, "Could not refresh tokens")
		}
		return ctx.JSON(http.StatusCreated, pairTokens)
	}
}

func (h *authHandlers) UserLogout() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tokenString := ctx.Get("user_id").(string)
		err := h.authUC.DeleteRefreshToken(ctx.Request().Context(), tokenString)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, "Could not logout with this token")
		}
		return ctx.JSON(http.StatusCreated, "Logout complited!")
	}
}

func NewAuthHandlers(authUC auth.UseCase, logger logging.Logger) auth.Handlers {
	return &authHandlers{authUC: authUC, logger: logger}
}
