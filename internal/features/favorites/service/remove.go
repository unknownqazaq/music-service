package service

import (
	"context"
)

func (s *FavoritesService) Remove(ctx context.Context, userID, trackID int64) error {
	return s.repo.Remove(ctx, userID, trackID)
}
