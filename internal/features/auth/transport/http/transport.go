package auth_transport_http

import (
	"context"
	"net/http"

	"music-service/internal/core/transport/http/server"
	"music-service/internal/features/auth/service"
	users_model "music-service/internal/features/users/model"
)

type AuthService interface {
	Register(ctx context.Context, input service.RegisterInput) (*users_model.User, error)
	Login(ctx context.Context, input service.LoginInput) (*service.LoginResult, error)
	Refresh(ctx context.Context, refreshToken string) (*service.LoginResult, error)
	Logout(ctx context.Context, refreshToken string) error
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Username string `json:"username" validate:"required,min=3"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (h *AuthHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.Refresh,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.Logout,
		},
	}
}
