package auth

import (
	"Medods/internal/auth/dto"
	"context"
)


type UseCase interface {
	CreateTokens(ctx context.Context, user_id string, userAgent string, user_ip string) (dto.TokensPair, error)
	RefreshToken(ctx context.Context, refToken string, accessToken, userAgent string, user_ip string) (dto.TokensPair, error)
	GetUserID(ctx context.Context, accessToken string) (string, error)
	DeleteRefreshToken(ctx context.Context, accessToken string) error
}
