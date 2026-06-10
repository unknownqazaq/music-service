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

// Login godoc
// @Summary      Authenticate user and retrieve token
// @Description  Login with email and password to receive a JWT access token and refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Credentials payload"
// @Success      200      {object}  TokenResponse
// @Failure      400      {object}  response.ErrorResponse "Invalid payload"
// @Failure      401      {object}  response.ErrorResponse "Invalid credentials"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req LoginRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request")
		return
	}

	result, err := h.authService.Login(ctx, service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			responseHandler.ErrorResponse(core_errors.ErrUnauthorized, err.Error())
			return
		}
		responseHandler.ErrorResponse(err, "failed to authenticate")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}


