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
	"music-service/internal/features/favorites/service"
)

// AddFavorite godoc
// @Summary      Add track to favorites
// @Description  Add a track by ID to the user's favorites list. FREE users can add up to 20 favorites.
// @Tags         favorites
// @Param        track_id  path      int64  true  "Track ID"
// @Success      201  "Created (Track added to favorites)"
// @Failure      400      {object}  response.ErrorResponse "Invalid track ID"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse "Forbidden (Favorites limit exceeded)"
// @Failure      409      {object}  response.ErrorResponse "Conflict (Track already in favorites)"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /favorites/tracks/{track_id} [post]
func (h *FavoritesHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
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

	err = h.favoritesService.Add(ctx, claims.UserID, trackID, claims.SubscriptionType)
	if err != nil {
		if errors.Is(err, service.ErrFavoritesLimitExceeded) {
			responseHandler.ErrorResponse(core_errors.ErrForbidden, err.Error())
			return
		}
		if errors.Is(err, favorites_postgres.ErrFavoriteAlreadyExists) {
			responseHandler.ErrorResponse(core_errors.ErrConflict, err.Error())
			return
		}
		responseHandler.ErrorResponse(err, "failed to add track to favorites")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

