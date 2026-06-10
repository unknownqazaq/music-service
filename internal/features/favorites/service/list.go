package service

import (
	"context"
	tracks_model "music-service/internal/features/tracks/model"
)

func (s *FavoritesService) List(ctx context.Context, userID int64) ([]tracks_model.Track, error) {
	return s.repo.List(ctx, userID)
}
