package playlists_transport_http

import (
	"context"
	"net/http"

	"music-service/internal/core/transport/http/server"
	"music-service/internal/features/playlists/model"
)

type PlaylistService interface {
	Create(ctx context.Context, p *model.Playlist, subscriptionType string) (*model.Playlist, error)
	GetByID(ctx context.Context, id, userID int64) (*model.Playlist, error)
	ListByUserID(ctx context.Context, userID int64) ([]model.Playlist, error)
	Update(ctx context.Context, id, userID int64, title string, description *string) (*model.Playlist, error)
	Delete(ctx context.Context, id, userID int64) error
	AddTrack(ctx context.Context, playlistID, trackID, userID int64) error
	RemoveTrack(ctx context.Context, playlistID, trackID, userID int64) error
}

type PlaylistHandler struct {
	playlistService PlaylistService
}

func NewPlaylistHandler(playlistService PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{playlistService: playlistService}
}

type CreatePlaylistRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=100"`
	Description *string `json:"description"`
}

func (h *PlaylistHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/playlists",
			Handler: h.ListPlaylists,
		},
		{
			Method:  http.MethodPost,
			Path:    "/playlists",
			Handler: h.CreatePlaylist,
		},
		{
			Method:  http.MethodGet,
			Path:    "/playlists/{id}",
			Handler: h.GetPlaylist,
		},
		{
			Method:  http.MethodPut,
			Path:    "/playlists/{id}",
			Handler: h.UpdatePlaylist,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/playlists/{id}",
			Handler: h.DeletePlaylist,
		},
		{
			Method:  http.MethodPost,
			Path:    "/playlists/{playlist_id}/tracks/{track_id}",
			Handler: h.AddTrack,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/playlists/{playlist_id}/tracks/{track_id}",
			Handler: h.RemoveTrack,
		},
	}
}
