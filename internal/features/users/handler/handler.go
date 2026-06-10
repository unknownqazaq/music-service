package handler

import (
	users_model "music-service/internal/features/users/model"
	"music-service/internal/features/users/service"
)

var _ = users_model.User{}

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
