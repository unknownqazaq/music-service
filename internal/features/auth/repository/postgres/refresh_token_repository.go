package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrTokenNotFound = errors.New("refresh token not found or expired")

type RefreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Save сохраняет refresh token для пользователя.
func (r *RefreshTokenRepository) Save(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (token) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}

// GetUserIDByToken возвращает user_id по валидному refresh token.
func (r *RefreshTokenRepository) GetUserIDByToken(ctx context.Context, token string) (int64, error) {
	query := `
		SELECT user_id FROM refresh_tokens
		WHERE token = $1 AND expires_at > NOW()
	`
	var userID int64
	err := r.db.QueryRowContext(ctx, query, token).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrTokenNotFound
		}
		return 0, fmt.Errorf("get user by refresh token: %w", err)
	}
	return userID, nil
}

// Delete удаляет refresh token (logout).
func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

// DeleteAllByUserID удаляет все токены пользователя.
func (r *RefreshTokenRepository) DeleteAllByUserID(ctx context.Context, userID int64) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete all refresh tokens: %w", err)
	}
	return nil
}
