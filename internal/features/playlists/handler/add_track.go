package handler

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

// AddTrack godoc
// @Summary      Add a track to a playlist
// @Description  Add a track by ID to a user's playlist by ID. Users can only modify their own playlists.
// @Tags         playlists
// @Param        playlist_id  path      int64  true  "Playlist ID"
// @Param        track_id     path      int64  true  "Track ID"
// @Success      201  "Created (Track added to playlist)"
// @Failure      400      {object}  response.ErrorResponse "Invalid playlist ID or track ID"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse "Forbidden (Not owner of the playlist)"
// @Failure      404      {object}  response.ErrorResponse "Playlist not found"
// @Failure      409      {object}  response.ErrorResponse "Conflict (Track already in playlist)"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /playlists/{playlist_id}/tracks/{track_id} [post]
func (h *PlaylistHandler) AddTrack(w http.ResponseWriter, r *http.Request) {
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

	err = h.playlistService.AddTrack(ctx, playlistID, trackID, claims.UserID)
	if err != nil {
		if errors.Is(err, playlist_postgres.ErrPlaylistNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, "playlist not found")
			return
		}
		if errors.Is(err, service.ErrForbiddenPlaylist) {
			responseHandler.ErrorResponse(core_errors.ErrForbidden, err.Error())
			return
		}
		if errors.Is(err, playlist_postgres.ErrTrackAlreadyInPlaylist) {
			responseHandler.ErrorResponse(core_errors.ErrConflict, err.Error())
			return
		}
		responseHandler.ErrorResponse(err, "failed to add track to playlist")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

