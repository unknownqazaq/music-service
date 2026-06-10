package service

import (
	"context"
	"errors"

	tracks_model "music-service/internal/features/tracks/model"
)

var (
	ErrFavoritesLimitExceeded = errors.New("free subscription limits to at most 20 favorites")
)

type FavoritesRepository interface {
	Add(ctx context.Context, userID, trackID int64) error
	Remove(ctx context.Context, userID, trackID int64) error
	CountByUserID(ctx context.Context, userID int64) (int, error)
	List(ctx context.Context, userID int64) ([]tracks_model.Track, error)
}

type FavoritesService struct {
	repo          FavoritesRepository
	favoriteLimit int
}

func NewFavoritesService(repo FavoritesRepository, favoriteLimit int) *FavoritesService {
	return &FavoritesService{
		repo:          repo,
		favoriteLimit: favoriteLimit,
	}
}
