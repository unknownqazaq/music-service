package tracks_transport_http

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	track_postgres "music-service/internal/features/tracks/repository/postgres"
)

// DeleteTrack godoc
// @Summary      Delete a track
// @Description  Soft delete a music track by ID (ADMIN only)
// @Tags         admin
// @Param        id   path      int64  true  "Track ID"
// @Success      204  "No Content (Track successfully deleted)"
// @Failure      400      {object}  response.ErrorResponse "Invalid track ID"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse "Forbidden (ADMIN required)"
// @Failure      404      {object}  response.ErrorResponse "Track not found"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /admin/tracks/{id} [delete]
func (h *TrackHandler) DeleteTrack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid track ID")
		return
	}

	err = h.trackService.SoftDelete(ctx, id)
	if err != nil {
		if errors.Is(err, track_postgres.ErrTrackNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, "track not found")
			return
		}
		responseHandler.ErrorResponse(err, "failed to delete track")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
