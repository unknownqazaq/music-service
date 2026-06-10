package service

import (
	"context"
	"music-service/internal/features/tracks/model"
)

func (s *TrackService) List(ctx context.Context, limit, offset int) ([]model.Track, error) {
	return s.repo.List(ctx, limit, offset)
}
