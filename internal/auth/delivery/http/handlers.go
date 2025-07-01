package http

import (
	"Medods/internal/auth"
	"Medods/internal/auth/dto"
	"Medods/pkg/logging"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type authHandlers struct {
	authUC auth.UseCase
	logger logging.Logger
}

// CreateTokens implements auth.Handlers.
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

// GetUserIDByToken implements auth.Handlers.
func (h *authHandlers) GetUserIDByToken() echo.HandlerFunc {
	return func(ctx echo.Context) error {
        authHeader := ctx.Request().Header.Get("Authorization")
        if authHeader == "" {
            return echo.NewHTTPError(http.StatusUnauthorized, "Failed to find token")
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        tokenString = strings.TrimSpace(tokenString)
        if tokenString == "" {
            return echo.NewHTTPError(http.StatusUnauthorized, "Failed to find token")
        }

		userID, err := h.authUC.GetUserID(ctx.Request().Context(), tokenString)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, "Could not find user for this token")
		}
		return ctx.JSON(http.StatusCreated, userID)
	}
}

// RefreshTokens implements auth.Handlers.
func (h *authHandlers) RefreshTokens() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := new(dto.RefreshTokenRequest)
		ctx.Bind(req)
		userAgent := ctx.Request().UserAgent()
		userIP := ctx.RealIP()
		refToken := req.RefreshToken

		authHeader := ctx.Request().Header.Get("Authorization")
        if authHeader == "" {
            return echo.NewHTTPError(http.StatusUnauthorized, "Failed to find token")
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        tokenString = strings.TrimSpace(tokenString)
        if tokenString == "" {
            return echo.NewHTTPError(http.StatusUnauthorized, "Failed to find token")
        }

		pairTokens, err := h.authUC.RefreshToken(ctx.Request().Context(), refToken, tokenString, userAgent, userIP)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, "Could not refresh tokens")
		}
		return ctx.JSON(http.StatusCreated, pairTokens)
	}
}

// UserLogout implements auth.Handlers.
func (h *authHandlers) UserLogout() echo.HandlerFunc {
	return func(ctx echo.Context) error {
        authHeader := ctx.Request().Header.Get("Authorization")
        if authHeader == "" {
            return echo.NewHTTPError(http.StatusUnauthorized, "Failed to find token")
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        tokenString = strings.TrimSpace(tokenString)
        if tokenString == "" {
            return echo.NewHTTPError(http.StatusUnauthorized, "Failed to find token")
        }

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
