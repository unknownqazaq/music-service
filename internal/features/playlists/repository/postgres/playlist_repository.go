package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"music-service/internal/features/playlists/model"
)

var (
	ErrPlaylistNotFound       = errors.New("playlist not found")
	ErrTrackAlreadyInPlaylist = errors.New("track is already in playlist")
)

type PlaylistRepository struct {
	db *sqlx.DB
}

func NewPlaylistRepository(db *sqlx.DB) *PlaylistRepository {
	return &PlaylistRepository{db: db}
}

func (r *PlaylistRepository) Create(ctx context.Context, p *model.Playlist) (*model.Playlist, error) {
	query := `
		INSERT INTO playlists (user_id, title, description, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, user_id, title, description, created_at
	`
	var created model.Playlist
	err := r.db.QueryRowxContext(ctx, query, p.UserID, p.Title, p.Description).StructScan(&created)
	if err != nil {
		return nil, fmt.Errorf("create playlist: %w", err)
	}
	return &created, nil
}

func (r *PlaylistRepository) GetByID(ctx context.Context, id int64) (*model.Playlist, error) {
	query := `SELECT id, user_id, title, description, created_at FROM playlists WHERE id = $1`
	var p model.Playlist
	err := r.db.GetContext(ctx, &p, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlaylistNotFound
		}
		return nil, fmt.Errorf("get playlist by id: %w", err)
	}
	return &p, nil
}

func (r *PlaylistRepository) ListByUserID(ctx context.Context, userID int64) ([]model.Playlist, error) {
	query := `SELECT id, user_id, title, description, created_at FROM playlists WHERE user_id = $1 ORDER BY id DESC`
	var playlists []model.Playlist
	err := r.db.SelectContext(ctx, &playlists, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list user playlists: %w", err)
	}
	return playlists, nil
}

func (r *PlaylistRepository) Update(ctx context.Context, p *model.Playlist) (*model.Playlist, error) {
	query := `
		UPDATE playlists 
		SET title = $1, description = $2
		WHERE id = $3
		RETURNING id, user_id, title, description, created_at
	`
	var updated model.Playlist
	err := r.db.QueryRowxContext(ctx, query, p.Title, p.Description, p.ID).StructScan(&updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlaylistNotFound
		}
		return nil, fmt.Errorf("update playlist: %w", err)
	}
	return &updated, nil
}

func (r *PlaylistRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM playlists WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete playlist: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrPlaylistNotFound
	}
	return nil
}

func (r *PlaylistRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM playlists WHERE user_id = $1`
	var count int
	err := r.db.GetContext(ctx, &count, query, userID)
	if err != nil {
		return 0, fmt.Errorf("count user playlists: %w", err)
	}
	return count, nil
}

func (r *PlaylistRepository) AddTrack(ctx context.Context, playlistID, trackID int64) error {
	query := `INSERT INTO playlist_tracks (playlist_id, track_id, added_at) VALUES ($1, $2, NOW())`
	_, err := r.db.ExecContext(ctx, query, playlistID, trackID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "23505") {
			return ErrTrackAlreadyInPlaylist
		}
		return fmt.Errorf("add track to playlist: %w", err)
	}
	return nil
}

func (r *PlaylistRepository) RemoveTrack(ctx context.Context, playlistID, trackID int64) error {
	query := `DELETE FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`
	res, err := r.db.ExecContext(ctx, query, playlistID, trackID)
	if err != nil {
		return fmt.Errorf("remove track from playlist: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("track not found in playlist")
	}
	return nil
}
