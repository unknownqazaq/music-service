package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"music-service/internal/features/users/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrUsernameTaken     = errors.New("username already taken")
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (email, password_hash, username, role, subscription_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, email, username, role, subscription_type, created_at, updated_at
	`
	var created model.User
	err := r.db.QueryRowxContext(ctx, query, u.Email, u.PasswordHash, u.Username, u.Role, u.SubscriptionType).
		StructScan(&created)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "users_email_key") {
			return nil, ErrEmailAlreadyTaken
		}
		if strings.Contains(errStr, "users_username_key") {
			return nil, ErrUsernameTaken
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &created, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT id, email, password_hash, username, role, subscription_type, created_at, updated_at FROM users WHERE id = $1`
	var u model.User
	err := r.db.GetContext(ctx, &u, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, password_hash, username, role, subscription_type, created_at, updated_at FROM users WHERE email = $1`
	var u model.User
	err := r.db.GetContext(ctx, &u, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) UpdateSubscription(ctx context.Context, id int64, subscriptionType string) error {
	query := `UPDATE users SET subscription_type = $1, updated_at = NOW() WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, subscriptionType, id)
	if err != nil {
		return fmt.Errorf("update user subscription: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id int64, email, username *string) (*model.User, error) {
	query := `
		UPDATE users 
		SET 
			email = COALESCE($1, email), 
			username = COALESCE($2, username), 
			updated_at = NOW() 
		WHERE id = $3
		RETURNING id, email, username, role, subscription_type, created_at, updated_at
	`
	var updated model.User
	err := r.db.QueryRowxContext(ctx, query, email, username, id).StructScan(&updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		errStr := err.Error()
		if strings.Contains(errStr, "users_email_key") {
			return nil, ErrEmailAlreadyTaken
		}
		if strings.Contains(errStr, "users_username_key") {
			return nil, ErrUsernameTaken
		}
		return nil, fmt.Errorf("update user profile: %w", err)
	}
	return &updated, nil
}


