package middleware

import (
	"Medods/internal/auth/dto"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)


func (mw *MiddlewareManager) AuthJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if !mw.validateJWTToken(ctx) {
			mw.logger.Info("Access токен не валиден или отсутствует")
			return echo.NewHTTPError(http.StatusUnauthorized, "Access токен не валиден или отсутствует")
		}
		return next(ctx)
	}
}

func (mw *MiddlewareManager) AuthRefreshMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if !mw.validateJWTToken(ctx) {
			mw.logger.Info("Access токен не валиден или отсутствует")
			return echo.NewHTTPError(http.StatusUnauthorized, "Access токен не валиден или отсутствует")
		}

		req := new(dto.RefreshTokenRequest)
		ctx.Bind(req)
		ctx.Set("refToken", req.RefreshToken)

		userAgent := ctx.Request().UserAgent()
		ctx.Set("userAgent", userAgent)

		userIP := ctx.RealIP()
		ctx.Set("userIP", userIP)

		return next(ctx)
	}
}

func (mw *MiddlewareManager) validateJWTToken(ctx echo.Context) bool{
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		mw.logger.Info("AuthHeader пустой")
		return false
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		mw.logger.Info("Отсутствует Access токен")
		return false
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
        return secretKey, nil
	})

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("test")
		userID := claims["sub"].(string)
		if err:= mw.checkSession(userID, ctx.Request().Context()); err != nil {
			return false
		}
		ctx.Set("user_id", userID)
		c := context.WithValue(ctx.Request().Context(), dto.UserIDKey, userID)
		ctx.SetRequest(ctx.Request().WithContext(c))
		return true
	}
	return false
}

func (mw *MiddlewareManager) checkSession(user_id string, ctx context.Context) error{
	_, err := mw.authUC.GetTokenByUserID(ctx, user_id)
	if err != nil {
		return err
	}
	return nil
}