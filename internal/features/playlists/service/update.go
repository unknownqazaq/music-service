package service

import (
	"context"
	"music-service/internal/features/playlists/model"
)

func (s *PlaylistService) Update(ctx context.Context, id, userID int64, title string, description *string) (*model.Playlist, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p.UserID != userID {
		return nil, ErrForbiddenPlaylist
	}

	p.Title = title
	p.Description = description

	return s.repo.Update(ctx, p)
}
