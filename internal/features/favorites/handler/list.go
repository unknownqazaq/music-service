package handler

import (
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
	"music-service/internal/core/response"
	tracks_model "music-service/internal/features/tracks/model"
)

var _ = tracks_model.Track{}

// ListFavorites godoc
// @Summary      List favorites
// @Description  Get all tracks added to favorites by the current user
// @Tags         favorites
// @Produce      json
// @Success      200      {array}   tracks_model.Track
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /favorites/tracks [get]
func (h *FavoritesHandler) ListFavorites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	claims, ok := middleware.UserClaimsFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "unauthorized")
		return
	}

	tracks, err := h.favoritesService.List(ctx, claims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to retrieve favorites list")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, tracks)
}
