package service

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	core_redis "music-service/internal/core/redis"
	"music-service/internal/features/tracks/model"
)

var (
	ErrDailyLimitExceeded = errors.New("daily listening limit exceeded")
)

type TrackRepository interface {
	Create(ctx context.Context, t *model.Track) (*model.Track, error)
	Update(ctx context.Context, t *model.Track) (*model.Track, error)
	SoftDelete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.Track, error)
	List(ctx context.Context, limit, offset int) ([]model.Track, error)
	Search(ctx context.Context, q string) ([]model.Track, error)
	SavePlay(ctx context.Context, userID, trackID int64) error
}

type TrackService struct {
	repo       TrackRepository
	cache      *core_redis.Cache
	rdb        *redis.Client
	dailyLimit int
}

func NewTrackService(repo TrackRepository, cache *core_redis.Cache, rdb *redis.Client, dailyLimit int) *TrackService {
	return &TrackService{
		repo:       repo,
		cache:      cache,
		rdb:        rdb,
		dailyLimit: dailyLimit,
	}
}
