package handler

import (
	"music-service/internal/features/favorites/service"
)

type FavoritesHandler struct {
	favoritesService *service.FavoritesService
}

func NewFavoritesHandler(favoritesService *service.FavoritesService) *FavoritesHandler {
	return &FavoritesHandler{favoritesService: favoritesService}
}
