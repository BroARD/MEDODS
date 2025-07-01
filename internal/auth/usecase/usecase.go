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
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	LenOfToken = 16
	TimeOfActionRefreshToken = time.Hour * 24
	TimeOfActionAccessToken = time.Minute * 5 
	SecretKey = "secret_key_here"
)

type authUseCase struct {
	repo   auth.Repository
	logger logging.Logger
}

// CreateTokens implements auth.UseCase.
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

	refTokenDB :=  models.RefreshToken{
		ID: uuid.NewString(),
		UserID: user_id,
		RefreshTokenHash: refTokenHash,
		UserAgent: userAgent,
		IP: user_ip,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(TimeOfActionRefreshToken),
	}

	u.logger.Info("Создание Refresh токена в БД")
	if err = u.repo.CreateRefreshToken(ctx, &refTokenDB); err != nil {
		u.logger.Info("Ошибка при создании рефреш токена в БД")
		return dto.TokensPair{}, err
	}

	u.logger.Info("Генерация Access токена")
	accessToken, err := generateAccessToken(user_id)
	if err != nil {
		u.logger.Info("Ошибка при генерации Access токена")
		return dto.TokensPair{}, err
	}

	tokensPair := dto.TokensPair{
		AccessToken: accessToken,
		RefreshToken: base64.URLEncoding.EncodeToString(newRefToken),
	}

	return tokensPair, nil
}

// DeleteRefreshToken implements auth.UseCase.
func (u *authUseCase) DeleteRefreshToken(ctx context.Context, accessToken string) error {
	u.logger.Info("Получения пользователя по аксесс токену")
	user_id, err := getUserIDFromToken(accessToken)
	if err != nil {
		return err
	}

	u.logger.Info("Удаление Рефреш токена по UserID")
	err = u.repo.DeleteRefreshTokenByUserID(ctx, user_id)
	if err != nil {
		return err
	}
	return nil
}

// GetUserID implements auth.UseCase.
func (u *authUseCase) GetUserID(ctx context.Context, accessToken string) (string, error) {
	u.logger.Info("Получение UserID из токена")
	user_id, err := getUserIDFromToken(accessToken)
	if err != nil {
		return "", err
	}

	u.logger.Info("Проверка существует ли Refresh Токен с этим userID")
	_, err = u.repo.GetRefreshTokenByUserID(ctx, user_id)
	if err != nil {
		return "", err
	}

	return user_id, nil
}

// RefreshToken implements auth.UseCase.
func (u *authUseCase) RefreshToken(ctx context.Context, refToken string, accessToken string, userAgent string, user_ip string) (dto.TokensPair, error) {
	u.logger.Info("Получение Рефреш токена по Хэшу")
	refTokenBase64, err := base64.URLEncoding.DecodeString(refToken)
	if err != nil {
		return dto.TokensPair{}, err
	}

	u.logger.Info("Получение UserID из токена")
	user_id, err := getUserIDFromToken(accessToken)
	if err != nil {
		return dto.TokensPair{}, err
	}

	u.logger.Info("Получение UserID из Access токена")
	refTokenDB, err := u.repo.GetRefreshTokenByUserID(ctx, user_id)
	if err != nil {
		return dto.TokensPair{}, err
	}
	
	u.logger.Info("Сравнение user id у токенов")
	if err := bcrypt.CompareHashAndPassword(refTokenDB.RefreshTokenHash, refTokenBase64); err != nil{
		return dto.TokensPair{}, err
	}

	u.logger.Info("Удаление старого рефреш токена из БД")
	err = u.repo.DeleteRefreshTokenByUserID(ctx, refTokenDB.UserID)
	if err != nil {
		return dto.TokensPair{}, err
	}
	u.logger.Info("Проверка на совпадения браузера")
	// Если браузер не совпадает с поршлым, то не выдаёт новую пару, но при этом удаляет старую
	if refTokenDB.UserAgent != userAgent {
		return dto.TokensPair{}, err
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
	
	accessToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func getUserIDFromToken(tokenString string) (string, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(SecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(string); ok {
			return userID, nil
		}
	}
	return "", fmt.Errorf("invalid token")
 }
