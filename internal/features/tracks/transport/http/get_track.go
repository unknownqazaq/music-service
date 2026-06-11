package tracks_transport_http

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	tracks_model "music-service/internal/features/tracks/model"
	track_postgres "music-service/internal/features/tracks/repository/postgres"
)

var _ = tracks_model.Track{}

// GetTrack godoc
// @Summary      Get a track by ID
// @Description  Get detailed information of a track by ID (cached via Redis)
// @Tags         tracks
// @Produce      json
// @Param        id   path      int64  true  "Track ID"
// @Success      200      {object}  tracks_model.Track
// @Failure      400      {object}  response.ErrorResponse "Invalid track ID"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      404      {object}  response.ErrorResponse "Track not found"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /tracks/{id} [get]
func (h *TrackHandler) GetTrack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid track ID")
		return
	}

	track, err := h.trackService.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, track_postgres.ErrTrackNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, "track not found")
			return
		}
		responseHandler.ErrorResponse(err, "failed to get track")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, track)
}
