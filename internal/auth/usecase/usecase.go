package usecase

import (
	"Medods/internal/auth"
	"Medods/internal/auth/dto"
	"Medods/internal/models"
	"Medods/pkg/logging"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	LenOfToken               = 16
	TimeOfActionRefreshToken = time.Hour * 24 
	TimeOfActionAccessToken  = time.Minute * 5
)

type authUseCase struct {
	repo   auth.Repository
	logger logging.Logger
}

func (u *authUseCase) GetTokenByUserID(ctx context.Context, user_id string) (models.RefreshToken, error) {
	u.logger.Info("Запуск получения Refresh токена по UserID")
	// Получения Refresh токена по UserID из БД
	return u.repo.GetRefreshTokenByUserID(ctx, user_id)
}

func (u *authUseCase) CreateTokens(ctx context.Context, user_id string, userAgent string, user_ip string) (dto.TokensPair, error) {
	u.logger.Info("Запуск создания пары токенов")

	// Генерация нового Refresh токена
	newRefToken := make([]byte, LenOfToken)
	if _, err := rand.Read(newRefToken); err != nil {
		return dto.TokensPair{}, err
	}

	// Генерация Хеша из Refresh токена
	refTokenHash, err := bcrypt.GenerateFromPassword(newRefToken, bcrypt.DefaultCost)
	if err != nil {
		return dto.TokensPair{}, err
	}

	// Создание модели для БД
	refTokenDB := models.RefreshToken{
		ID:               uuid.NewString(),
		UserID:           user_id,
		RefreshTokenHash: refTokenHash,
		UserAgent:        userAgent,
		IP:               user_ip,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(TimeOfActionRefreshToken),
	}

	// Генерацйия Access токена
	accessToken, err := generateAccessToken(user_id, u.logger)
	if err != nil {
		u.logger.Info("Ошибка при генерации Access токена")
		return dto.TokensPair{}, err
	}

	// Проверка существования токена(удаление по наличии)
	if u.checkTokenExistence(ctx, user_id) {
		u.logger.Info("Обнаружена существующая пара токенов, удаление")
		u.repo.DeleteRefreshTokenByUserID(ctx, user_id)
	}
	
	// Создание нового Refresh Token в БД
	if err = u.repo.CreateRefreshToken(ctx, &refTokenDB); err != nil {
		u.logger.Info("Ошибка при создании рефреш токена в БД")
		return dto.TokensPair{}, err
	}

	// Пара токенов для возврата в json
	tokensPair := dto.TokensPair{
		AccessToken:  accessToken,
		RefreshToken: base64.URLEncoding.EncodeToString(newRefToken),
	}

	return tokensPair, nil
}

func (u *authUseCase) DeleteRefreshToken(ctx context.Context, accessToken string) error {
	u.logger.Info("Запуск удаления Refresh токена")
	// Получения UserID из context
	user_id := ctx.Value(dto.UserIDKey).(string)

	// Удаление Refresh Token из БД
	err := u.repo.DeleteRefreshTokenByUserID(ctx, user_id)
	if err != nil {
		u.logger.Info("DeletRefToken: Ошибка при удалении из БД")
		return err
	}
	return nil
}

func (u *authUseCase) RefreshToken(ctx context.Context, refToken string, accessToken string, userAgent string, user_ip string) (dto.TokensPair, error) {	
	u.logger.Info("Запуск обновления пары токенов")
	
	refTokenBase64, err := base64.URLEncoding.DecodeString(refToken)
	if err != nil {
		u.logger.Info("RefreshToken: Ошибка декодирования Refresh Token")
		return dto.TokensPair{}, err
	}

	// Получение UserID из context
	user_id := ctx.Value(dto.UserIDKey).(string)

	// Получение Refresh Token по UserID
	refTokenDB, err := u.repo.GetRefreshTokenByUserID(ctx, user_id)
	if err != nil {
		u.logger.Info("RefreshToken: Ошибка получения рефреш токена из БД по UserID")
		return dto.TokensPair{}, err
	}
	
	// Проверка валидности Refresh Token
	if !checkValidRefreshToken(refTokenDB) {
		u.logger.Info("RefreshToken: Refresh токен не прошёл проверку на валидность")
		return dto.TokensPair{}, fmt.Errorf("время жизни Refresh токена истекло")
	}

	// Проверка на совпадение хешей токенов(Принадлежит ли Access токен Refresh токену)
	if err := bcrypt.CompareHashAndPassword(refTokenDB.RefreshTokenHash, refTokenBase64); err != nil {
		u.logger.Info("RefreshToken: UserID у Access и Refresh токена не совпадают")
		return dto.TokensPair{}, err
	}

	// ОТПРАВКА WEBHOOKa
	if user_ip != refTokenDB.IP {
		u.logger.Info("RefreshToken: не совпадение user_ip, отправка WebHook")
		wh_url := "http://example.com"
		wh_payload := WebHookPayload{
			UserID: user_id,
			IP: user_ip,
			Event: "Запрос с нового IP",
		}
		
		go sendWebHook(wh_url, wh_payload, u.logger)
	}

	// Удаление старого Refresh Token
	err = u.repo.DeleteRefreshTokenByUserID(ctx, refTokenDB.UserID)
	if err != nil {
		u.logger.Info("RefreshToken: ошибка при удалении Refresh токена из БД")
		return dto.TokensPair{}, err
	}

	// Если браузер не совпадает с сохранённым, то не выдаёт новую пару, но при этом удаляет старую
	if refTokenDB.UserAgent != userAgent {
		u.logger.Info("RefreshToken: не совпадает User-Agent")
		return dto.TokensPair{}, err
	}

	// Запуск создание новой пары
	return u.CreateTokens(ctx, refTokenDB.UserID, userAgent, user_ip)
}

func NewAuthUseCase(repo auth.Repository, logger logging.Logger) auth.UseCase {
	return &authUseCase{repo: repo, logger: logger}
}

func generateAccessToken(user_id string, logger logging.Logger) (string, error) {
	claims := jwt.MapClaims{
		"sub": user_id,
		"exp": TimeOfActionAccessToken,
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	accessToken, err := token.SignedString(secretKey)
	if err != nil {
		logger.Info("GenerateAccessToken: Ошибка при подписании токена")
		return "", err
	}
	return accessToken, nil
}

func checkValidRefreshToken(refToken models.RefreshToken) bool {
	return time.Now().Before(refToken.ExpiresAt)
}

func (u *authUseCase) checkTokenExistence(ctx context.Context, user_id string) bool {
	if _, err := u.GetTokenByUserID(ctx, user_id); err != nil {
		u.logger.Info("RefreshToken: Токен не обнаружен")
		return false
	}
	return true
}
