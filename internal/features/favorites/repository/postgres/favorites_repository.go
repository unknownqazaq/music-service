package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	tracks_model "music-service/internal/features/tracks/model"
)

var (
	ErrFavoriteAlreadyExists = errors.New("track is already in favorites")
	ErrFavoriteNotFound      = errors.New("track is not in favorites")
)

type FavoritesRepository struct {
	db *sqlx.DB
}

func NewFavoritesRepository(db *sqlx.DB) *FavoritesRepository {
	return &FavoritesRepository{db: db}
}

func (r *FavoritesRepository) Add(ctx context.Context, userID, trackID int64) error {
	query := `INSERT INTO favorites (user_id, track_id, created_at) VALUES ($1, $2, NOW())`
	_, err := r.db.ExecContext(ctx, query, userID, trackID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "23505") {
			return ErrFavoriteAlreadyExists
		}
		return fmt.Errorf("add favorite: %w", err)
	}
	return nil
}

func (r *FavoritesRepository) Remove(ctx context.Context, userID, trackID int64) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND track_id = $2`
	res, err := r.db.ExecContext(ctx, query, userID, trackID)
	if err != nil {
		return fmt.Errorf("remove favorite: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrFavoriteNotFound
	}
	return nil
}

func (r *FavoritesRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM favorites WHERE user_id = $1`
	var count int
	err := r.db.GetContext(ctx, &count, query, userID)
	if err != nil {
		return 0, fmt.Errorf("count user favorites: %w", err)
	}
	return count, nil
}

func (r *FavoritesRepository) List(ctx context.Context, userID int64) ([]tracks_model.Track, error) {
	query := `
		SELECT t.id, t.title, t.artist_id, t.album_id, t.genre_id, t.duration_seconds, t.file_url, t.is_active, t.created_at,
		       art.name AS artist_name, alb.title AS album_title, g.name AS genre_name
		FROM tracks t
		JOIN favorites f ON t.id = f.track_id
		JOIN artists art ON t.artist_id = art.id
		LEFT JOIN albums alb ON t.album_id = alb.id
		LEFT JOIN genres g ON t.genre_id = g.id
		WHERE f.user_id = $1 AND t.is_active = TRUE
		ORDER BY f.created_at DESC
	`
	var tracks []tracks_model.Track
	err := r.db.SelectContext(ctx, &tracks, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list user favorites: %w", err)
	}
	return tracks, nil
}
