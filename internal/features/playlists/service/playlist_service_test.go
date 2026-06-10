package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"music-service/internal/features/playlists/model"
	"music-service/internal/features/playlists/service"
)

type mockPlaylistRepository struct {
	playlistsCount int
	service.PlaylistRepository
}

func (m *mockPlaylistRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	return m.playlistsCount, nil
}

func (m *mockPlaylistRepository) Create(ctx context.Context, p *model.Playlist) (*model.Playlist, error) {
	return p, nil
}

func TestPlaylistService_Create_LimitFREE(t *testing.T) {
	repo := &mockPlaylistRepository{playlistsCount: 3}
	svc := service.NewPlaylistService(repo, 3)

	p := &model.Playlist{UserID: 1, Title: "My Playlist"}

	// Create for FREE subscription (count is already 3)
	created, err := svc.Create(context.Background(), p, "FREE")
	assert.ErrorIs(t, err, service.ErrPlaylistLimitExceeded)
	assert.Nil(t, created)

	// Create for PREMIUM subscription (ignores limit)
	createdPremium, err := svc.Create(context.Background(), p, "PREMIUM")
	assert.NoError(t, err)
	assert.NotNil(t, createdPremium)
}
