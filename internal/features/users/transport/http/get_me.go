package users_transport_http

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
	"music-service/internal/core/response"
	users_model "music-service/internal/features/users/model"
	users_postgres "music-service/internal/features/users/repository/postgres"
)

var _ = users_model.User{}

// GetMe godoc
// @Summary      Get current user profile
// @Description  Get detailed profile information of the currently authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  users_model.User
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      404  {object}  response.ErrorResponse "User not found"
// @Failure      500  {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /users/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	claims, ok := middleware.UserClaimsFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "unauthorized")
		return
	}

	u, err := h.userService.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, users_postgres.ErrUserNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, "user not found")
			return
		}
		responseHandler.ErrorResponse(err, "failed to retrieve profile")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, u)
}
