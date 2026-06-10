package service

import (
	"context"

	"go.uber.org/zap"
	"music-service/internal/core/logger"
)

// Logout удаляет refresh token пользователя.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	log := logger.FromContext(ctx)

	err := s.refreshRepo.Delete(ctx, refreshToken)
	if err != nil {
		log.Error("logout: failed to delete refresh token", zap.Error(err))
		return err
	}

	log.Info("user logged out")
	return nil
}
