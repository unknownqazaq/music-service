package playlists_transport_http

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	playlist_postgres "music-service/internal/features/playlists/repository/postgres"
	"music-service/internal/features/playlists/service"
)

// DeletePlaylist godoc
// @Summary      Delete a playlist
// @Description  Delete a user's playlist by ID. Users can only delete their own playlists.
// @Tags         playlists
// @Param        id   path      int64  true  "Playlist ID"
// @Success      204  "No Content (Playlist successfully deleted)"
// @Failure      400      {object}  response.ErrorResponse "Invalid playlist ID"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse "Forbidden (Not owner of the playlist)"
// @Failure      404      {object}  response.ErrorResponse "Playlist not found"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /playlists/{id} [delete]
func (h *PlaylistHandler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
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
		responseHandler.ErrorResponse(err, "invalid playlist ID")
		return
	}

	err = h.playlistService.Delete(ctx, id, claims.UserID)
	if err != nil {
		if errors.Is(err, playlist_postgres.ErrPlaylistNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, "playlist not found")
			return
		}
		if errors.Is(err, service.ErrForbiddenPlaylist) {
			responseHandler.ErrorResponse(core_errors.ErrForbidden, err.Error())
			return
		}
		responseHandler.ErrorResponse(err, "failed to delete playlist")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
