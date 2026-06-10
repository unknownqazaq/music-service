package service

import (
	"context"

	"go.uber.org/zap"
	"music-service/internal/core/logger"
	"music-service/internal/features/playlists/model"
)

func (s *PlaylistService) Create(ctx context.Context, p *model.Playlist, subscriptionType string) (*model.Playlist, error) {
	log := logger.FromContext(ctx)

	if subscriptionType == "FREE" {
		count, err := s.repo.CountByUserID(ctx, p.UserID)
		if err != nil {
			return nil, err
		}
		if count >= s.playlistLimitFREE {
			log.Warn("FREE subscription playlist limit exceeded",
				zap.Int64("user_id", p.UserID),
				zap.Int("current_playlists", count),
				zap.Int("limit", s.playlistLimitFREE),
			)
			return nil, ErrPlaylistLimitExceeded
		}
	}

	created, err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	log.Info("playlist created successfully",
		zap.Int64("user_id", created.UserID),
		zap.Int64("playlist_id", created.ID),
		zap.String("title", created.Title),
		zap.String("subscription_type", subscriptionType),
	)

	return created, nil
}

