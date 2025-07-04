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
	u.logger.Info("Получение токена по UserID")
	return u.repo.GetRefreshTokenByUserID(ctx, user_id)
}

func (u *authUseCase) CreateTokens(ctx context.Context, user_id string, userAgent string, user_ip string) (dto.TokensPair, error) {
	newRefToken := make([]byte, LenOfToken)
	u.logger.Info("Генерация Refresh токена")
	if _, err := rand.Read(newRefToken); err != nil {
		return dto.TokensPair{}, err
	}

	u.logger.Info("Генерация хэша Refresh токена")
	refTokenHash, err := bcrypt.GenerateFromPassword(newRefToken, bcrypt.DefaultCost)
	if err != nil {
		return dto.TokensPair{}, err
	}

	refTokenDB := models.RefreshToken{
		ID:               uuid.NewString(),
		UserID:           user_id,
		RefreshTokenHash: refTokenHash,
		UserAgent:        userAgent,
		IP:               user_ip,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(TimeOfActionRefreshToken),
	}

	u.logger.Info("Генерация Access токена")
	accessToken, err := generateAccessToken(user_id)
	if err != nil {
		u.logger.Info("Ошибка при генерации Access токена")
		return dto.TokensPair{}, err
	}

	if u.checkTokenExistence(ctx, user_id) {
		u.repo.DeleteRefreshTokenByUserID(ctx, user_id)
	}
	
	u.logger.Info("Создание Refresh токена в БД")
	if err = u.repo.CreateRefreshToken(ctx, &refTokenDB); err != nil {
		u.logger.Info("Ошибка при создании рефреш токена в БД")
		return dto.TokensPair{}, err
	}

	tokensPair := dto.TokensPair{
		AccessToken:  accessToken,
		RefreshToken: base64.URLEncoding.EncodeToString(newRefToken),
	}

	return tokensPair, nil
}

func (u *authUseCase) DeleteRefreshToken(ctx context.Context, accessToken string) error {
	u.logger.Info("Получения пользователя по аксесс токену")
	user_id := ctx.Value(dto.UserIDKey).(string)

	u.logger.Info("Удаление Рефреш токена по UserID")
	err := u.repo.DeleteRefreshTokenByUserID(ctx, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (u *authUseCase) RefreshToken(ctx context.Context, refToken string, accessToken string, userAgent string, user_ip string) (dto.TokensPair, error) {
	u.logger.Info("Получение Рефреш токена по Хэшу")
	refTokenBase64, err := base64.URLEncoding.DecodeString(refToken)
	if err != nil {
		return dto.TokensPair{}, err
	}

	u.logger.Info("Получение UserID из токена")
	user_id := ctx.Value(dto.UserIDKey).(string)

	u.logger.Info("Получение UserID из Refresh токена")
	refTokenDB, err := u.repo.GetRefreshTokenByUserID(ctx, user_id)
	if err != nil {
		return dto.TokensPair{}, err
	}
	u.logger.Info("Проверка валидности Refresh токена")
	if !checkValidRefreshToken(refTokenDB) {
		return dto.TokensPair{}, fmt.Errorf("время жизни Refresh токена истекло")
	}

	u.logger.Info("Сравнение user id у токенов")
	if err := bcrypt.CompareHashAndPassword(refTokenDB.RefreshTokenHash, refTokenBase64); err != nil {
		return dto.TokensPair{}, err
	}

	u.logger.Info("Удаление старого рефреш токена из БД")
	err = u.repo.DeleteRefreshTokenByUserID(ctx, refTokenDB.UserID)
	if err != nil {
		return dto.TokensPair{}, err
	}
	u.logger.Info("Проверка на совпадения браузера")
	// Если браузер не совпадает с сохранённым, то не выдаёт новую пару, но при этом удаляет старую
	if refTokenDB.UserAgent != userAgent {
		return dto.TokensPair{}, err
	}

	if refTokenDB.IP != user_ip {

	}

	u.logger.Info("Запуск создания новой пары")
	return u.CreateTokens(ctx, refTokenDB.UserID, userAgent, user_ip)
}

func NewAuthUseCase(repo auth.Repository, logger logging.Logger) auth.UseCase {
	return &authUseCase{repo: repo, logger: logger}
}

func generateAccessToken(user_id string) (string, error) {
	claims := jwt.MapClaims{
		"sub": user_id,
		"exp": TimeOfActionAccessToken,
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	accessToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func checkValidRefreshToken(refToken models.RefreshToken) bool {
	return time.Now().Before(refToken.ExpiresAt)
}

func (u *authUseCase) checkTokenExistence(ctx context.Context, user_id string) bool {
	if _, err := u.GetTokenByUserID(ctx, user_id); err != nil {
		return false
	}
	return true
}
