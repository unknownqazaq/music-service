package tracks_transport_http

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	"music-service/internal/features/tracks/model"
	track_postgres "music-service/internal/features/tracks/repository/postgres"
)

// UpdateTrack godoc
// @Summary      Update a track
// @Description  Update details of a music track (ADMIN only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id       path      int64               true  "Track ID"
// @Param        request  body      CreateTrackRequest  true  "Updated track configuration payload"
// @Success      200      {object}  model.Track
// @Failure      400      {object}  response.ErrorResponse "Invalid payload or ID"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse "Forbidden (ADMIN required)"
// @Failure      404      {object}  response.ErrorResponse "Track not found"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /admin/tracks/{id} [put]
func (h *TrackHandler) UpdateTrack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid track ID")
		return
	}

	var req CreateTrackRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request")
		return
	}

	track := &model.Track{
		ID:              id,
		Title:           req.Title,
		ArtistID:        req.ArtistID,
		AlbumID:         req.AlbumID,
		GenreID:         req.GenreID,
		DurationSeconds: req.DurationSeconds,
		FileURL:         req.FileURL,
		IsActive:        req.IsActive,
	}

	updated, err := h.trackService.Update(ctx, track)
	if err != nil {
		if errors.Is(err, track_postgres.ErrTrackNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, "track not found")
			return
		}
		responseHandler.ErrorResponse(err, "failed to update track")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, updated)
}
