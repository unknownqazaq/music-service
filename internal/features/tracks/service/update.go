package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"music-service/internal/core/logger"
	"music-service/internal/features/tracks/model"
)

func (s *TrackService) Update(ctx context.Context, t *model.Track) (*model.Track, error) {
	log := logger.FromContext(ctx)

	_ = s.cache.Delete(ctx, fmt.Sprintf("track:%d", t.ID))

	updated, err := s.repo.Update(ctx, t)
	if err != nil {
		return nil, err
	}

	log.Info("admin action: updated track",
		zap.Int64("track_id", updated.ID),
		zap.String("title", updated.Title),
		zap.Int64("artist_id", updated.ArtistID),
	)

	return updated, nil
}

