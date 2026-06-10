package handler

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	playlists_model "music-service/internal/features/playlists/model"
	playlist_postgres "music-service/internal/features/playlists/repository/postgres"
	"music-service/internal/features/playlists/service"
)

var _ = playlists_model.Playlist{}

// UpdatePlaylist godoc
// @Summary      Update a playlist
// @Description  Update the title and description of a user's playlist. Users can only update their own playlists.
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        id       path      int64                  true  "Playlist ID"
// @Param        request  body      CreatePlaylistRequest  true  "Updated playlist payload"
// @Success      200      {object}  playlists_model.Playlist
// @Failure      400      {object}  response.ErrorResponse "Invalid payload or ID"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse "Forbidden (Not owner of the playlist)"
// @Failure      404      {object}  response.ErrorResponse "Playlist not found"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /playlists/{id} [put]
func (h *PlaylistHandler) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
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

	var req CreatePlaylistRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request")
		return
	}

	updated, err := h.playlistService.Update(ctx, id, claims.UserID, req.Title, req.Description)
	if err != nil {
		if errors.Is(err, playlist_postgres.ErrPlaylistNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, "playlist not found")
			return
		}
		if errors.Is(err, service.ErrForbiddenPlaylist) {
			responseHandler.ErrorResponse(core_errors.ErrForbidden, err.Error())
			return
		}
		responseHandler.ErrorResponse(err, "failed to update playlist")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, updated)
}

