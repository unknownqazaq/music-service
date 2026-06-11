package users_transport_http

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	users_model "music-service/internal/features/users/model"
	users_postgres "music-service/internal/features/users/repository/postgres"
)

var _ = users_model.User{}

type UpdateProfileRequest struct {
	Email    *string `json:"email" validate:"omitempty,email"`
	Username *string `json:"username" validate:"omitempty,min=3,max=100"`
}

// UpdateProfile godoc
// @Summary      Update current user profile
// @Description  Update the email and/or username of the currently authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      UpdateProfileRequest  true  "Profile update payload"
// @Success      200      {object}  users_model.User
// @Failure      400      {object}  response.ErrorResponse "Invalid payload"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      409      {object}  response.ErrorResponse "Conflict (Email or Username already taken)"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /users/me [patch]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	claims, ok := middleware.UserClaimsFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request")
		return
	}

	updatedUser, err := h.userService.UpdateProfile(ctx, claims.UserID, req.Email, req.Username)
	if err != nil {
		if errors.Is(err, users_postgres.ErrUserNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, "user not found")
			return
		}
		if errors.Is(err, users_postgres.ErrEmailAlreadyTaken) || errors.Is(err, users_postgres.ErrUsernameTaken) {
			responseHandler.ErrorResponse(core_errors.ErrConflict, err.Error())
			return
		}
		responseHandler.ErrorResponse(err, "failed to update profile")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, updatedUser)
}
