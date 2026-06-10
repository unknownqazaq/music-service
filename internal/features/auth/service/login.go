package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	log := logger.FromContext(ctx)

	u, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		log.Warn("user authorization failed: user not found", zap.String("email", input.Email))
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password)); err != nil {
		log.Warn("user authorization failed: invalid password", zap.String("email", input.Email), zap.Int64("user_id", u.ID))
		return nil, ErrInvalidCredentials
	}

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

	// Генерируем refresh token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Сохраняем refresh token в БД
	expiresAt := time.Now().Add(s.refreshTTL)
	if err := s.refreshRepo.Save(ctx, u.ID, refreshToken, expiresAt); err != nil {
		log.Warn("failed to save refresh token", zap.Int64("user_id", u.ID), zap.Error(err))
		// Не фатально, возвращаем только access token
		return &LoginResult{AccessToken: accessToken}, nil
	}

	log.Info("user authorized successfully",
		zap.Int64("user_id", u.ID),
		zap.String("email", u.Email),
		zap.String("role", string(u.Role)),
		zap.String("subscription_type", string(u.SubscriptionType)),
	)

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
