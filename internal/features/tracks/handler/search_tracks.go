package handler

import (
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/response"
	tracks_model "music-service/internal/features/tracks/model"
)

var _ = tracks_model.Track{}

// SearchTracks godoc
// @Summary      Search tracks
// @Description  Search active tracks by title or artist name
// @Tags         tracks
// @Produce      json
// @Param        query  query     string  true  "Search query"
// @Success      200      {array}   tracks_model.Track
// @Failure      400      {object}  response.ErrorResponse "Query parameter is required"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /tracks/search [get]
func (h *TrackHandler) SearchTracks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	query := r.URL.Query().Get("query")
	if query == "" {
		responseHandler.ErrorResponse(core_errors.ErrBadRequest, "query parameter is required")
		return
	}

	tracks, err := h.trackService.Search(ctx, query)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to search tracks")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, tracks)
}
