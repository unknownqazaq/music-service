package model

import "time"

const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

const (
	SubscriptionFree    = "FREE"
	SubscriptionPremium = "PREMIUM"
)

type User struct {
	ID               int64     `db:"id" json:"id"`
	Email            string    `db:"email" json:"email"`
	PasswordHash     string    `db:"password_hash" json:"-"`
	Username         string    `db:"username" json:"username"`
	Role             string    `db:"role" json:"role"`
	SubscriptionType string    `db:"subscription_type" json:"subscription_type"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
