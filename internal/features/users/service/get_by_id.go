package service

import (
	"context"
	"music-service/internal/features/users/model"
)

func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}
