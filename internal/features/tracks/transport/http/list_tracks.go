package tracks_transport_http

import (
	"net/http"

	"music-service/internal/core/logger"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	tracks_model "music-service/internal/features/tracks/model"
)

var _ = tracks_model.Track{}

// ListTracks godoc
// @Summary      List tracks
// @Description  Get a paginated list of active music tracks
// @Tags         tracks
// @Produce      json
// @Param        limit  query     int  false  "Limit (default 20)"
// @Param        page   query     int  false  "Page (default 1)"
// @Success      200      {array}   tracks_model.Track
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /tracks [get]
func (h *TrackHandler) ListTracks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	limitPtr, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid limit parameter")
		return
	}

	pagePtr, err := core_http_request.GetIntQueryParam(r, "page")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid page parameter")
		return
	}

	limit := 20
	if limitPtr != nil && *limitPtr > 0 {
		limit = *limitPtr
	}

	page := 1
	if pagePtr != nil && *pagePtr > 0 {
		page = *pagePtr
	}

	offset := (page - 1) * limit
	tracks, err := h.trackService.List(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to list tracks")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, tracks)
}
