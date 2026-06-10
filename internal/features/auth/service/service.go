package service

import (
	"context"
	"errors"
	"time"

	users_model "music-service/internal/features/users/model"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

type UserRepository interface {
	Create(ctx context.Context, u *users_model.User) (*users_model.User, error)
	GetByEmail(ctx context.Context, email string) (*users_model.User, error)
	GetByID(ctx context.Context, id int64) (*users_model.User, error)
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	GetUserIDByToken(ctx context.Context, token string) (int64, error)
	Delete(ctx context.Context, token string) error
	DeleteAllByUserID(ctx context.Context, userID int64) error
}

type AuthService struct {
	repo             UserRepository
	refreshRepo      RefreshTokenRepository
	jwtSecret        string
	jwtRefreshSecret string
	jwtExpiration    time.Duration
	refreshTTL       time.Duration
}

func NewAuthService(
	repo UserRepository,
	refreshRepo RefreshTokenRepository,
	jwtSecret string,
	jwtRefreshSecret string,
	jwtExpiration time.Duration,
	refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		repo:             repo,
		refreshRepo:      refreshRepo,
		jwtSecret:        jwtSecret,
		jwtRefreshSecret: jwtRefreshSecret,
		jwtExpiration:    jwtExpiration,
		refreshTTL:       refreshTTL,
	}
}
