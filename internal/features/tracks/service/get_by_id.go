package service

import (
	"context"
	"fmt"
	"time"
	"music-service/internal/features/tracks/model"
)

func (s *TrackService) GetByID(ctx context.Context, id int64) (*model.Track, error) {
	cacheKey := fmt.Sprintf("track:%d", id)
	var t model.Track
	found, err := s.cache.Get(ctx, cacheKey, &t)
	if err == nil && found {
		return &t, nil
	}

	dbTrack, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, dbTrack, 10*time.Minute)

	return dbTrack, nil
}
