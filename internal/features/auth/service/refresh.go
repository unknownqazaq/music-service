package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
)

// Refresh проверяет refresh_token и возвращает новую пару токенов.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	log := logger.FromContext(ctx)

	userID, err := s.refreshRepo.GetUserIDByToken(ctx, refreshToken)
	if err != nil {
		log.Warn("refresh token invalid or expired")
		return nil, ErrInvalidRefreshToken
	}

	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// Удаляем старый refresh token (ротация)
	_ = s.refreshRepo.Delete(ctx, refreshToken)

	// Создаём новый access token
	claims := &middleware.Claims{
		UserID:           u.ID,
		Email:            u.Email,
		Role:             u.Role,
		SubscriptionType: u.SubscriptionType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	// Создаём новый refresh token
	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(s.refreshTTL)
	if err := s.refreshRepo.Save(ctx, u.ID, newRefreshToken, expiresAt); err != nil {
		log.Error("failed to save new refresh token", zap.Error(err))
	}

	log.Info("token refreshed",
		zap.Int64("user_id", u.ID),
		zap.String("email", u.Email),
	)

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
