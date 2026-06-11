package users_transport_http

import (
	"context"
	"net/http"

	"music-service/internal/core/transport/http/server"
	users_model "music-service/internal/features/users/model"
)

type UserService interface {
	GetByID(ctx context.Context, id int64) (*users_model.User, error)
	UpdateProfile(ctx context.Context, id int64, email, username *string) (*users_model.User, error)
	UpdateSubscription(ctx context.Context, id int64, subscriptionType string) error
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users/me",
			Handler: h.GetMe,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/me",
			Handler: h.UpdateProfile,
		},
	}
}
