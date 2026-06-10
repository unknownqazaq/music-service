package handler

import (
	users_service "music-service/internal/features/users/service"
)

type SubscriptionHandler struct {
	userService *users_service.UserService
}

func NewSubscriptionHandler(userService *users_service.UserService) *SubscriptionHandler {
	return &SubscriptionHandler{userService: userService}
}

type UpdateSubscriptionRequest struct {
	SubscriptionType string `json:"subscription_type" validate:"required,oneof=FREE PREMIUM"`
}

