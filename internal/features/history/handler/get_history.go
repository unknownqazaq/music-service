package handler

import (
	"net/http"

	core_errors "music-service/internal/core/errors"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
	"music-service/internal/core/response"
	history_model "music-service/internal/features/history/model"
)

var _ = history_model.HistoryEntry{}

// GetHistory godoc
// @Summary      Get listening history
// @Description  Get a log of all songs listened to by the current user
// @Tags         history
// @Produce      json
// @Success      200      {array}   history_model.HistoryEntry
// @Failure      401      {object}  response.ErrorResponse "Unauthorized"
// @Failure      500      {object}  response.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /listening-history [get]
func (h *HistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	claims, ok := middleware.UserClaimsFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "unauthorized")
		return
	}

	history, err := h.repo.GetByUserID(ctx, claims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get listening history")
		return
	}

	responseHandler.JSONResponse(http.StatusOK, history)
}
