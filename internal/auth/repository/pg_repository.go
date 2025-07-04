package repository

import (
	"Medods/internal/auth"
	"Medods/internal/models"
	"context"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func (r *authRepository) CreateRefreshToken(ctx context.Context, refToken *models.RefreshToken) error {
	err := r.db.WithContext(ctx).Create(refToken).Error
	return err
}

func (r *authRepository) DeleteRefreshTokenByUserID(ctx context.Context, user_id string) error {
	err := r.db.WithContext(ctx).Delete(&models.RefreshToken{}, "user_id = ?", user_id).Error
	return err
}

func (r *authRepository) GetRefreshTokenByUserID(ctx context.Context, user_id string) (models.RefreshToken, error) {
	var refToken models.RefreshToken
	err := r.db.WithContext(ctx).First(&refToken, "user_id = ?", user_id).Error
	return refToken, err
}

func NewAuthRepository(db *gorm.DB) auth.Repository {
	return &authRepository{db: db}
}
