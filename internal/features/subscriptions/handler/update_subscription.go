package handler

import (
	"errors"
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/response"
	core_http_request "music-service/internal/core/transport/http/request"
	users_postgres "music-service/internal/features/users/repository/postgres"
)

// UpdateSubscription godoc
// @Summary      Update user subscription
// @Description  Update the subscription type of a user (ADMIN only). Types: FREE, PREMIUM.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id       path      int64                     true  "User ID"
// @Param        request  body      UpdateSubscriptionRequest  true  "Updated subscription payload"
// @Success      200      {object}  map[string]string "Success response"
// @Failure      400      {object}  response.ErrorResponse "Invalid payload or ID"
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse "Forbidden (ADMIN required)"
// @Failure      404      {object}  response.ErrorResponse "User not found"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /admin/users/{id}/subscription [patch]
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid user ID")
		return
	}

	var req UpdateSubscriptionRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request")
		return
	}

	err = h.userService.UpdateSubscription(ctx, id, req.SubscriptionType)
	if err != nil {
		if errors.Is(err, users_postgres.ErrUserNotFound) {
			responseHandler.ErrorResponse(core_errors.ErrNotFound, "user not found")
			return
		}
		responseHandler.ErrorResponse(err, "failed to update user subscription")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, map[string]string{"status": "success"})
}

