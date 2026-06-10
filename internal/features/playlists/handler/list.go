package handler

import (
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
	"music-service/internal/core/response"
	playlists_model "music-service/internal/features/playlists/model"
)

var _ = playlists_model.Playlist{}

// ListPlaylists godoc
// @Summary      List playlists
// @Description  Get all playlists created by the current user
// @Tags         playlists
// @Produce      json
// @Success      200      {array}   playlists_model.Playlist
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /playlists [get]
func (h *PlaylistHandler) ListPlaylists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	claims, ok := middleware.UserClaimsFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "unauthorized")
		return
	}

	playlists, err := h.playlistService.ListByUserID(ctx, claims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to list playlists")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, playlists)
}
