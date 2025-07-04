package auth

import (
	"Medods/internal/auth/dto"
	"Medods/internal/models"
	"context"
)


type UseCase interface {
	CreateTokens(ctx context.Context, user_id string, userAgent string, user_ip string) (dto.TokensPair, error)
	RefreshToken(ctx context.Context, refToken string, accessToken, userAgent string, user_ip string) (dto.TokensPair, error)
	DeleteRefreshToken(ctx context.Context, accessToken string) error
	GetTokenByUserID(ctx context.Context, user_id string) (models.RefreshToken, error)
}
