package service

import (
	"context"

	"go.uber.org/zap"
	"music-service/internal/core/logger"
	"music-service/internal/features/tracks/model"
)

func (s *TrackService) Create(ctx context.Context, t *model.Track) (*model.Track, error) {
	log := logger.FromContext(ctx)

	created, err := s.repo.Create(ctx, t)
	if err != nil {
		return nil, err
	}

	log.Info("admin action: created track",
		zap.Int64("track_id", created.ID),
		zap.String("title", created.Title),
		zap.Int64("artist_id", created.ArtistID),
	)

	return created, nil
}

