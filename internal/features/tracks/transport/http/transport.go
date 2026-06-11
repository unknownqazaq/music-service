package tracks_transport_http

import (
	"context"
	"net/http"

	"music-service/internal/core/transport/http/server"
	"music-service/internal/features/tracks/model"
)

type TrackService interface {
	Create(ctx context.Context, t *model.Track) (*model.Track, error)
	Update(ctx context.Context, t *model.Track) (*model.Track, error)
	SoftDelete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.Track, error)
	List(ctx context.Context, limit, offset int) ([]model.Track, error)
	Search(ctx context.Context, q string) ([]model.Track, error)
	Play(ctx context.Context, userID, trackID int64, subscriptionType string) (*model.Track, error)
}

type TrackHandler struct {
	trackService TrackService
}

func NewTrackHandler(trackService TrackService) *TrackHandler {
	return &TrackHandler{trackService: trackService}
}

type CreateTrackRequest struct {
	Title           string  `json:"title" validate:"required,min=1,max=200"`
	ArtistID        int64   `json:"artist_id" validate:"required,gt=0"`
	AlbumID         *int64  `json:"album_id"`
	GenreID         *int64  `json:"genre_id"`
	DurationSeconds int     `json:"duration_seconds" validate:"required,gt=0"`
	FileURL         string  `json:"file_url" validate:"required,url"`
	IsActive        bool    `json:"is_active"`
}

func (h *TrackHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/tracks",
			Handler: h.ListTracks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tracks/{id}",
			Handler: h.GetTrack,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tracks/search",
			Handler: h.SearchTracks,
		},
		{
			Method:  http.MethodPost,
			Path:    "/tracks/{id}/play",
			Handler: h.PlayTrack,
		},
	}
}

func (h *TrackHandler) AdminRoutes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/admin/tracks",
			Handler: h.CreateTrack,
		},
		{
			Method:  http.MethodPut,
			Path:    "/admin/tracks/{id}",
			Handler: h.UpdateTrack,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/admin/tracks/{id}",
			Handler: h.DeleteTrack,
		},
	}
}
