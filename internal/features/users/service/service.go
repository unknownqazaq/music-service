package service

import (
	"context"
	"music-service/internal/features/users/model"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateSubscription(ctx context.Context, id int64, subscriptionType string) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}
