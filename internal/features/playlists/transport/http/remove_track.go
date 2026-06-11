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

// RemoveTrack godoc
// @Summary      Remove a track from a playlist
// @Description  Remove a track by ID from a user's playlist by ID. Users can only modify their own playlists.
// @Tags         playlists
// @Param        playlist_id  path      int64  true  "Playlist ID"
// @Param        track_id     path      int64  true  "Track ID"
// @Success      204  "No Content (Track successfully removed from playlist)"
// @Failure      400      {object}  response.ErrorResponse "Invalid playlist ID or track ID"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse "Forbidden (Not owner of the playlist)"
// @Failure      404      {object}  response.ErrorResponse "Playlist not found"
// @Security     BearerAuth
// @Router       /playlists/{playlist_id}/tracks/{track_id} [delete]
func (h *PlaylistHandler) RemoveTrack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	claims, ok := middleware.UserClaimsFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "unauthorized")
		return
	}

	playlistID, err := core_http_request.GetIntPathValue(r, "playlist_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid playlist ID")
		return
	}

	trackID, err := core_http_request.GetIntPathValue(r, "track_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid track ID")
		return
	}

	err = h.playlistService.RemoveTrack(ctx, playlistID, trackID, claims.UserID)
	if err != nil {
		if errors.Is(err, playlist_postgres.ErrPlaylistNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, "playlist not found")
			return
		}
		if errors.Is(err, service.ErrForbiddenPlaylist) {
			responseHandler.ErrorResponse(core_errors.ErrForbidden, err.Error())
			return
		}
		responseHandler.ErrorResponse(err, "failed to remove track from playlist")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
