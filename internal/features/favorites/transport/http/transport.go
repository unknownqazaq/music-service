package favorites_transport_http

import (
	"context"
	"net/http"

	"music-service/internal/core/transport/http/server"
	tracks_model "music-service/internal/features/tracks/model"
)

type FavoritesService interface {
	Add(ctx context.Context, userID, trackID int64, subscriptionType string) error
	Remove(ctx context.Context, userID, trackID int64) error
	List(ctx context.Context, userID int64) ([]tracks_model.Track, error)
}

type FavoritesHandler struct {
	favoritesService FavoritesService
}

func NewFavoritesHandler(favoritesService FavoritesService) *FavoritesHandler {
	return &FavoritesHandler{favoritesService: favoritesService}
}

func (h *FavoritesHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/favorites/tracks",
			Handler: h.ListFavorites,
		},
		{
			Method:  http.MethodPost,
			Path:    "/favorites/tracks/{track_id}",
			Handler: h.AddFavorite,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/favorites/tracks/{track_id}",
			Handler: h.RemoveFavorite,
		},
	}
}
