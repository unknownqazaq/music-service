package service

import (
	"context"
	"music-service/internal/features/playlists/model"
)

func (s *PlaylistService) GetByID(ctx context.Context, id, userID int64) (*model.Playlist, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p.UserID != userID {
		return nil, ErrForbiddenPlaylist
	}

	return p, nil
}
