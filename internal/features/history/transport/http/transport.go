package history_transport_http

import (
	"context"
	"net/http"

	"music-service/internal/core/transport/http/server"
	history_model "music-service/internal/features/history/model"
)

type HistoryRepository interface {
	GetByUserID(ctx context.Context, userID int64) ([]history_model.HistoryEntry, error)
}

type HistoryHandler struct {
	repo HistoryRepository
}

func NewHistoryHandler(repo HistoryRepository) *HistoryHandler {
	return &HistoryHandler{repo: repo}
}

func (h *HistoryHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/listening-history",
			Handler: h.GetHistory,
		},
	}
}
