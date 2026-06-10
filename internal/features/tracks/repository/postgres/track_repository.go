package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"music-service/internal/features/tracks/model"
)

var ErrTrackNotFound = errors.New("track not found")

type TrackRepository struct {
	db *sqlx.DB
}

func NewTrackRepository(db *sqlx.DB) *TrackRepository {
	return &TrackRepository{db: db}
}

func (r *TrackRepository) Create(ctx context.Context, t *model.Track) (*model.Track, error) {
	query := `
		INSERT INTO tracks (title, artist_id, album_id, genre_id, duration_seconds, file_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, artist_id, album_id, genre_id, duration_seconds, file_url, is_active, created_at
	`
	var created model.Track
	err := r.db.QueryRowxContext(ctx, query, t.Title, t.ArtistID, t.AlbumID, t.GenreID, t.DurationSeconds, t.FileURL, t.IsActive).
		StructScan(&created)
	if err != nil {
		return nil, fmt.Errorf("create track: %w", err)
	}
	return &created, nil
}

func (r *TrackRepository) Update(ctx context.Context, t *model.Track) (*model.Track, error) {
	query := `
		UPDATE tracks 
		SET title = $1, artist_id = $2, album_id = $3, genre_id = $4, duration_seconds = $5, file_url = $6, is_active = $7
		WHERE id = $8
		RETURNING id, title, artist_id, album_id, genre_id, duration_seconds, file_url, is_active, created_at
	`
	var updated model.Track
	err := r.db.QueryRowxContext(ctx, query, t.Title, t.ArtistID, t.AlbumID, t.GenreID, t.DurationSeconds, t.FileURL, t.IsActive, t.ID).
		StructScan(&updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTrackNotFound
		}
		return nil, fmt.Errorf("update track: %w", err)
	}
	return &updated, nil
}

func (r *TrackRepository) SoftDelete(ctx context.Context, id int64) error {
	query := `UPDATE tracks SET is_active = FALSE WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("soft delete track: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrTrackNotFound
	}
	return nil
}

func (r *TrackRepository) GetByID(ctx context.Context, id int64) (*model.Track, error) {
	query := `
		SELECT t.id, t.title, t.artist_id, t.album_id, t.genre_id, t.duration_seconds, t.file_url, t.is_active, t.created_at,
		       art.name AS artist_name, alb.title AS album_title, g.name AS genre_name
		FROM tracks t
		JOIN artists art ON t.artist_id = art.id
		LEFT JOIN albums alb ON t.album_id = alb.id
		LEFT JOIN genres g ON t.genre_id = g.id
		WHERE t.id = $1
	`
	var t model.Track
	err := r.db.GetContext(ctx, &t, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTrackNotFound
		}
		return nil, fmt.Errorf("get track by id: %w", err)
	}
	return &t, nil
}

func (r *TrackRepository) List(ctx context.Context, limit, offset int) ([]model.Track, error) {
	query := `
		SELECT t.id, t.title, t.artist_id, t.album_id, t.genre_id, t.duration_seconds, t.file_url, t.is_active, t.created_at,
		       art.name AS artist_name, alb.title AS album_title, g.name AS genre_name
		FROM tracks t
		JOIN artists art ON t.artist_id = art.id
		LEFT JOIN albums alb ON t.album_id = alb.id
		LEFT JOIN genres g ON t.genre_id = g.id
		WHERE t.is_active = TRUE
		ORDER BY t.id ASC
		LIMIT $1 OFFSET $2
	`
	var tracks []model.Track
	err := r.db.SelectContext(ctx, &tracks, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list active tracks: %w", err)
	}
	return tracks, nil
}

func (r *TrackRepository) Search(ctx context.Context, q string) ([]model.Track, error) {
	query := `
		SELECT t.id, t.title, t.artist_id, t.album_id, t.genre_id, t.duration_seconds, t.file_url, t.is_active, t.created_at,
		       art.name AS artist_name, alb.title AS album_title, g.name AS genre_name
		FROM tracks t
		JOIN artists art ON t.artist_id = art.id
		LEFT JOIN albums alb ON t.album_id = alb.id
		LEFT JOIN genres g ON t.genre_id = g.id
		WHERE t.is_active = TRUE AND (
			t.title ILIKE $1 OR 
			art.name ILIKE $1 OR 
			alb.title ILIKE $1 OR 
			g.name ILIKE $1
		)
		ORDER BY t.id ASC
	`
	var tracks []model.Track
	searchPattern := "%" + q + "%"
	err := r.db.SelectContext(ctx, &tracks, query, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("search tracks: %w", err)
	}
	return tracks, nil
}

func (r *TrackRepository) SavePlay(ctx context.Context, userID, trackID int64) error {
	query := `INSERT INTO listening_history (user_id, track_id, listened_at) VALUES ($1, $2, NOW())`
	_, err := r.db.ExecContext(ctx, query, userID, trackID)
	if err != nil {
		return fmt.Errorf("save track play history: %w", err)
	}
	return nil
}
