package auth_transport_http

import (
	"net/http"

	"music-service/internal/core/logger"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
)

// Logout godoc
// @Summary      Logout user
// @Description  Invalidates the given refresh_token so it can no longer be used
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LogoutRequest  true  "Refresh token to invalidate"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  response.ErrorResponse "Invalid payload"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req LogoutRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request")
		return
	}

	if err := h.authService.Logout(ctx, req.RefreshToken); err != nil {
		responseHandler.ErrorResponse(err, "failed to logout")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, map[string]string{"message": "logged out successfully"})
}
