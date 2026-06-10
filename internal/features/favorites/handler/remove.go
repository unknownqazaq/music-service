package handler

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	favorites_postgres "music-service/internal/features/favorites/repository/postgres"
)

// RemoveFavorite godoc
// @Summary      Remove track from favorites
// @Description  Remove a track by ID from the user's favorites list
// @Tags         favorites
// @Param        track_id  path      int64  true  "Track ID"
// @Success      204  "No Content (Track successfully removed from favorites)"
// @Failure      400      {object}  response.ErrorResponse "Invalid track ID"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      404      {object}  response.ErrorResponse "Track not found in favorites"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /favorites/tracks/{track_id} [delete]
func (h *FavoritesHandler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	claims, ok := middleware.UserClaimsFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "unauthorized")
		return
	}

	trackID, err := core_http_request.GetIntPathValue(r, "track_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid track ID")
		return
	}

	err = h.favoritesService.Remove(ctx, claims.UserID, trackID)
	if err != nil {
		if errors.Is(err, favorites_postgres.ErrFavoriteNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, err.Error())
			return
		}
		responseHandler.ErrorResponse(err, "failed to remove track from favorites")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

