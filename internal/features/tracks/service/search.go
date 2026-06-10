package service

import (
	"context"
	"music-service/internal/features/tracks/model"
)

func (s *TrackService) Search(ctx context.Context, query string) ([]model.Track, error) {
	return s.repo.Search(ctx, query)
}
