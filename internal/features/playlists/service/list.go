package service

import (
	"context"
	"music-service/internal/features/playlists/model"
)

func (s *PlaylistService) ListByUserID(ctx context.Context, userID int64) ([]model.Playlist, error) {
	return s.repo.ListByUserID(ctx, userID)
}
