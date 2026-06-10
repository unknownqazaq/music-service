package handler

import (
	"net/http"

	"music-service/internal/core/logger"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	"music-service/internal/features/tracks/model"
)

// CreateTrack godoc
// @Summary      Create a new track
// @Description  Create a new music track (ADMIN only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request  body      CreateTrackRequest  true  "Track configuration payload"
// @Success      201      {object}  model.Track
// @Failure      400      {object}  response.ErrorResponse "Invalid request body or missing required fields"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse "Forbidden (ADMIN required)"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /admin/tracks [post]
func (h *TrackHandler) CreateTrack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req CreateTrackRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request")
		return
	}

	track := &model.Track{
		Title:           req.Title,
		ArtistID:        req.ArtistID,
		AlbumID:         req.AlbumID,
		GenreID:         req.GenreID,
		DurationSeconds: req.DurationSeconds,
		FileURL:         req.FileURL,
		IsActive:        req.IsActive,
	}

	created, err := h.trackService.Create(ctx, track)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create track")
		return
	}

	responseHandler.JSONResponse(http.StatusCreated, created)
}

