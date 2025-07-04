package middleware

import (
	"Medods/config"
	"Medods/internal/auth"
	"Medods/pkg/logging"
)

type MiddlewareManager struct {
	authUC auth.UseCase
	cfg *config.Config
	logger logging.Logger
}

func NewMiddlewareManager(authUC auth.UseCase, cfg *config.Config, logger logging.Logger) *MiddlewareManager {
	return &MiddlewareManager{authUC: authUC, cfg: cfg, logger: logger}
}
