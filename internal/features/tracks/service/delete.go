package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"music-service/internal/core/logger"
)

func (s *TrackService) SoftDelete(ctx context.Context, id int64) error {
	log := logger.FromContext(ctx)

	_ = s.cache.Delete(ctx, fmt.Sprintf("track:%d", id))

	err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		return err
	}

	log.Info("admin action: soft deleted track", zap.Int64("track_id", id))

	return nil
}

