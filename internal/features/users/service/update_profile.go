package service

import (
	"context"
	"music-service/internal/features/users/model"
)

func (s *UserService) UpdateProfile(ctx context.Context, id int64, email, username *string) (*model.User, error) {
	if email == nil && username == nil {
		return s.repo.GetByID(ctx, id)
	}
	return s.repo.UpdateProfile(ctx, id, email, username)
}
