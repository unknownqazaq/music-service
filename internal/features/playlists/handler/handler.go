package handler

import (
	"music-service/internal/features/playlists/service"
)

type PlaylistHandler struct {
	playlistService *service.PlaylistService
}

func NewPlaylistHandler(playlistService *service.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{playlistService: playlistService}
}

type CreatePlaylistRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=100"`
	Description *string `json:"description"`
}

