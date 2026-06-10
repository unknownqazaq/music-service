package service

import (
	"context"
	"errors"

	"music-service/internal/features/playlists/model"
)

var (
	ErrPlaylistLimitExceeded = errors.New("free subscription limits to at most 3 playlists")
	ErrForbiddenPlaylist     = errors.New("access denied: you do not own this playlist")
)

type PlaylistRepository interface {
	Create(ctx context.Context, p *model.Playlist) (*model.Playlist, error)
	GetByID(ctx context.Context, id int64) (*model.Playlist, error)
	ListByUserID(ctx context.Context, userID int64) ([]model.Playlist, error)
	Update(ctx context.Context, p *model.Playlist) (*model.Playlist, error)
	Delete(ctx context.Context, id int64) error
	CountByUserID(ctx context.Context, userID int64) (int, error)
	AddTrack(ctx context.Context, playlistID, trackID int64) error
	RemoveTrack(ctx context.Context, playlistID, trackID int64) error
}

type PlaylistService struct {
	repo              PlaylistRepository
	playlistLimitFREE int
}

func NewPlaylistService(repo PlaylistRepository, playlistLimitFREE int) *PlaylistService {
	return &PlaylistService{
		repo:              repo,
		playlistLimitFREE: playlistLimitFREE,
	}
}
