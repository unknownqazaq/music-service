package handler

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	"music-service/internal/features/auth/service"
)

// Refresh godoc
// @Summary      Refresh access token
// @Description  Exchange a valid refresh_token for a new access_token and refresh_token pair
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest  true  "Refresh token payload"
// @Success      200      {object}  TokenResponse
// @Failure      400      {object}  response.ErrorResponse "Invalid payload"
// @Failure      401      {object}  response.ErrorResponse "Invalid or expired refresh token"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req RefreshRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request")
		return
	}

	result, err := h.authService.Refresh(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorized, err.Error())
			return
		}
		responseHandler.ErrorResponse(err, "failed to refresh token")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}
