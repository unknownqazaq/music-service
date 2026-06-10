package service_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_redis "music-service/internal/core/redis"
	"music-service/internal/features/tracks/model"
	"music-service/internal/features/tracks/service"
)

type mockTrackRepository struct {
	track *model.Track
	service.TrackRepository
}

func (m *mockTrackRepository) GetByID(ctx context.Context, id int64) (*model.Track, error) {
	return m.track, nil
}

func (m *mockTrackRepository) SavePlay(ctx context.Context, userID, trackID int64) error {
	return nil
}

func TestTrackService_Play_DailyLimitFREE(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rdb.Close()

	cache := core_redis.NewCache(rdb)

	track := &model.Track{
		ID:              1,
		Title:           "Test Song",
		ArtistID:        10,
		DurationSeconds: 180,
		FileURL:         "http://test.com",
		IsActive:        true,
	}

	repo := &mockTrackRepository{track: track}
	svc := service.NewTrackService(repo, cache, rdb, 10)

	// Simulate 10 plays for user 1 (FREE)
	for i := 0; i < 10; i++ {
		played, err := svc.Play(context.Background(), 1, 1, "FREE")
		assert.NoError(t, err)
		assert.NotNil(t, played)
	}

	// 11th play should exceed daily limit
	played, err := svc.Play(context.Background(), 1, 1, "FREE")
	assert.ErrorIs(t, err, service.ErrDailyLimitExceeded)
	assert.Nil(t, played)

	// Premium user should have no limits
	playedPremium, err := svc.Play(context.Background(), 1, 1, "PREMIUM")
	assert.NoError(t, err)
	assert.NotNil(t, playedPremium)
}
