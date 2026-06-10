package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"music-service/internal/features/favorites/service"
)

type mockFavoritesRepository struct {
	favoritesCount int
	service.FavoritesRepository
}

func (m *mockFavoritesRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	return m.favoritesCount, nil
}

func (m *mockFavoritesRepository) Add(ctx context.Context, userID, trackID int64) error {
	return nil
}

func TestFavoritesService_Add_LimitFREE(t *testing.T) {
	repo := &mockFavoritesRepository{favoritesCount: 20}
	svc := service.NewFavoritesService(repo, 20)

	// Add for FREE subscription (count is already 20)
	err := svc.Add(context.Background(), 1, 101, "FREE")
	assert.ErrorIs(t, err, service.ErrFavoritesLimitExceeded)

	// Add for PREMIUM subscription (ignores limit)
	errPremium := svc.Add(context.Background(), 1, 101, "PREMIUM")
	assert.NoError(t, errPremium)
}
