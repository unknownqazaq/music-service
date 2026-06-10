package handler

import (
	"music-service/internal/features/tracks/service"
)

type TrackHandler struct {
	trackService *service.TrackService
}

func NewTrackHandler(trackService *service.TrackService) *TrackHandler {
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

