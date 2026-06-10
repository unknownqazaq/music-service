package service

import (
	"context"
	users_model "music-service/internal/features/users/model"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"music-service/internal/core/logger"
)

type RegisterInput struct {
	Email    string
	Password string
	Username string
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*users_model.User, error) {
	log := logger.FromContext(ctx)

	bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &users_model.User{
		Email:            input.Email,
		PasswordHash:     string(bytes),
		Username:         input.Username,
		Role:             users_model.RoleUser,
		SubscriptionType: users_model.SubscriptionFree,
	}

	created, err := s.repo.Create(ctx, u)
	if err != nil {
		log.Warn("user registration failed: database error or conflict", zap.String("email", input.Email), zap.Error(err))
		return nil, err
	}

	log.Info("user registered successfully",
		zap.Int64("user_id", created.ID),
		zap.String("email", created.Email),
		zap.String("username", created.Username),
		zap.String("role", string(created.Role)),
	)

	return created, nil
}

