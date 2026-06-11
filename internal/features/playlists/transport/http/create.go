package playlists_transport_http

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	"music-service/internal/features/playlists/model"
	"music-service/internal/features/playlists/service"
)

// CreatePlaylist godoc
// @Summary      Create a playlist
// @Description  Create a new playlist for the current user. FREE users can create up to 3 playlists.
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        request  body      CreatePlaylistRequest  true  "Playlist creation payload"
// @Success      201      {object}  model.Playlist
// @Failure      400      {object}  response.ErrorResponse "Invalid payload or missing title"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse "Forbidden (Playlist limit exceeded)"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /playlists [post]
func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	claims, ok := middleware.UserClaimsFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "unauthorized")
		return
	}

	var req CreatePlaylistRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request")
		return
	}

	p := &model.Playlist{
		UserID:      claims.UserID,
		Title:       req.Title,
		Description: req.Description,
	}

	created, err := h.playlistService.Create(ctx, p, claims.SubscriptionType)
	if err != nil {
		if errors.Is(err, service.ErrPlaylistLimitExceeded) {
			responseHandler.ErrorResponse(core_errors.ErrForbidden, err.Error())
			return
		}
		responseHandler.ErrorResponse(err, "failed to create playlist")
		return
	}

	responseHandler.JSONResponse(http.StatusCreated, created)
}
