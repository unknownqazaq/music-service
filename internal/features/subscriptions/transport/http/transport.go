package subscriptions_transport_http

import (
	"context"
	"net/http"

	"music-service/internal/core/transport/http/server"
)

type UserService interface {
	UpdateSubscription(ctx context.Context, id int64, subscriptionType string) error
}

type SubscriptionHandler struct {
	userService UserService
}

func NewSubscriptionHandler(userService UserService) *SubscriptionHandler {
	return &SubscriptionHandler{userService: userService}
}

type UpdateSubscriptionRequest struct {
	SubscriptionType string `json:"subscription_type" validate:"required,oneof=FREE PREMIUM"`
}

func (h *SubscriptionHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPatch,
			Path:    "/admin/users/{id}/subscription",
			Handler: h.UpdateSubscription,
		},
	}
}
