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

// @Summary Создать пару токенов (access и refresh) по user_id
// @Tags JWTTokens
// @Description Создаёт пару токенов для пользователя с указанным user_id
// @ID create-tokens
// @Param user_id path string true "user_id"
// @Produce json
// @Success 201 {object} dto.TokensPair "Пара токенов успешно создана"
// @Failure 400 {object} dto.ErrorResponse "Неверный запрос (например, отсутствует user_id)"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера при создании токенов"
// @Router /tokens/{user_id} [post]
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

// @Summary Получить user_id по access токену
// @Tags JWTTokens
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "Возвращает user_id, связанный с access токеном"
// @Failure 401 {object} dto.ErrorResponse "Неавторизован: отсутствует или неверный токен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /tokens [get]
func (h *authHandlers) GetUserIDByToken() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		h.logger.Info("Получение UserID по Access токену")
		userID := ctx.Get("user_id").(string)
		return ctx.JSON(http.StatusOK, userID)
	}
}

// @Summary Обновить пару токенов по refresh токену
// @Tags JWTTokens
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dto.RefreshTokenRequest true "Refresh токен"
// @Success 201 {object} dto.TokensPair "Пара токенов успешно обновлена"
// @Failure 400 {object} dto.ErrorResponse "Неверный запрос (например, отсутствует refresh токен)"
// @Failure 401 {object} dto.ErrorResponse "Неавторизован: неверный или просроченный refresh токен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера при обновлении токенов"
// @Router /refresh [post]
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

// @Summary Выход пользователя (удаление refresh токена)
// @Tags JWTTokens
// @Security BearerAuth
// @Produce json
// @Success 201 {string} string "Logout completed!"
// @Failure 401 {object} dto.ErrorResponse "Неавторизован: отсутствует или неверный токен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера при выходе"
// @Router /logout [post]
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
