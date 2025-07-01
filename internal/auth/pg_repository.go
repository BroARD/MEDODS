package auth

import (
	"Medods/internal/models"
	"context"
)

type Repository interface {
	CreateRefreshToken(ctx context.Context, refToken *models.RefreshToken) error
	GetRefreshTokenByUserID(ctx context.Context, user_id string) (models.RefreshToken, error)
	DeleteRefreshTokenByUserID(ctx context.Context, user_id string) error
}