package handler

import (
	"music-service/internal/features/history/repository/postgres"
)

type HistoryHandler struct {
	repo *postgres.HistoryRepository
}

func NewHistoryHandler(repo *postgres.HistoryRepository) *HistoryHandler {
	return &HistoryHandler{repo: repo}
}
