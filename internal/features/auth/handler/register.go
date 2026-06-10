package handler

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	"music-service/internal/features/auth/service"
	users_model "music-service/internal/features/users/model"
	users_postgres "music-service/internal/features/users/repository/postgres"
)

var _ = users_model.User{}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user in the system with email, password, and username
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Registration payload"
// @Success      201      {object}  users_model.User
// @Failure      400      {object}  response.ErrorResponse "Invalid payload"
// @Failure      409      {object}  response.ErrorResponse "Conflict (Email or Username already taken)"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req RegisterRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request")
		return
	}

	u, err := h.authService.Register(ctx, service.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	})

	if err != nil {
		if errors.Is(err, users_postgres.ErrEmailAlreadyTaken) || errors.Is(err, users_postgres.ErrUsernameTaken) {
			responseHandler.ErrorResponse(core_errors.ErrConflict, err.Error())
			return
		}
		responseHandler.ErrorResponse(err, "failed to register user")
		return
	}

	responseHandler.JSONResponse(http.StatusCreated, u)
}

