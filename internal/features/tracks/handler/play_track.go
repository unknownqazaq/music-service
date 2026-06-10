package handler

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	tracks_model "music-service/internal/features/tracks/model"
	track_postgres "music-service/internal/features/tracks/repository/postgres"
	"music-service/internal/features/tracks/service"
)

var _ = tracks_model.Track{}

// PlayTrack godoc
// @Summary      Play a track
// @Description  Listen to a track by ID. For FREE accounts, this is limited to 10 plays per day.
// @Tags         tracks
// @Produce      json
// @Param        id   path      int64  true  "Track ID"
// @Success      200      {object}  tracks_model.Track
// @Failure      400      {object}  response.ErrorResponse "Invalid track ID"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse "Forbidden (Daily listening limit exceeded)"
// @Failure      404      {object}  response.ErrorResponse "Track not found"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /tracks/{id}/play [post]
func (h *TrackHandler) PlayTrack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	claims, ok := middleware.UserClaimsFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "unauthorized")
		return
	}

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid track ID")
		return
	}

	track, err := h.trackService.Play(ctx, claims.UserID, id, claims.SubscriptionType)
	if err != nil {
		if errors.Is(err, service.ErrDailyLimitExceeded) {
			responseHandler.ErrorResponse(core_errors.ErrForbidden, "daily listening limit exceeded")
			return
		}
		if errors.Is(err, track_postgres.ErrTrackNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, "track not found")
			return
		}
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponse(http.StatusOK, track)
}

