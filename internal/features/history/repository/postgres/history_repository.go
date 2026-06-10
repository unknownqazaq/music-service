package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"music-service/internal/features/history/model"
)

type HistoryRepository struct {
	db *sqlx.DB
}

func NewHistoryRepository(db *sqlx.DB) *HistoryRepository {
	return &HistoryRepository{db: db}
}

func (r *HistoryRepository) GetByUserID(ctx context.Context, userID int64) ([]model.HistoryEntry, error) {
	query := `
		SELECT lh.id, lh.track_id, t.title AS track_title, art.name AS artist_name, lh.listened_at
		FROM listening_history lh
		JOIN tracks t ON lh.track_id = t.id
		JOIN artists art ON t.artist_id = art.id
		WHERE lh.user_id = $1
		ORDER BY lh.listened_at DESC
	`
	var history []model.HistoryEntry
	err := r.db.SelectContext(ctx, &history, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get user history: %w", err)
	}
	return history, nil
}
