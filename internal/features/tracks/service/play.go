package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"music-service/internal/core/logger"
	"music-service/internal/features/tracks/model"

	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

func (s *TrackService) Play(ctx context.Context, userID, trackID int64, subscriptionType string) (*model.Track, error) {
	log := logger.FromContext(ctx)

	track, err := s.GetByID(ctx, trackID)
	if err != nil {
		return nil, err
	}

	if !track.IsActive {
		log.Warn("attempt to play inactive track",
			zap.Int64("user_id", userID),
			zap.Int64("track_id", trackID),
		)
		return nil, errors.New("track is inactive")
	}

	if subscriptionType == "FREE" {
		dateStr := time.Now().Format("2006-01-02")
		redisKey := fmt.Sprintf("user:%d:daily_play_count:%s", userID, dateStr)

		valStr, err := s.rdb.Get(ctx, redisKey).Result()
		if err != nil && err != redis.Nil {
			log.Error("failed to read daily play count from redis",
				zap.Int64("user_id", userID),
				zap.String("redis_key", redisKey),
				zap.Error(err),
			)
			return nil, fmt.Errorf("read daily limit: %w", err)
		}

		var count int
		if err == nil {
			fmt.Sscanf(valStr, "%d", &count)
		}

		if count >= s.dailyLimit {
			log.Warn("FREE subscription daily play limit exceeded",
				zap.Int64("user_id", userID),
				zap.Int64("track_id", trackID),
				zap.Int("plays_today", count),
				zap.Int("daily_limit", s.dailyLimit),
			)
			return nil, ErrDailyLimitExceeded
		}

		pipe := s.rdb.TxPipeline()
		pipe.Incr(ctx, redisKey)
		pipe.Expire(ctx, redisKey, 24*time.Hour)
		_, err = pipe.Exec(ctx)
		if err != nil {
			log.Error("failed to increment play count in redis",
				zap.Int64("user_id", userID),
				zap.String("redis_key", redisKey),
				zap.Error(err),
			)
			return nil, fmt.Errorf("increment play count: %w", err)
		}
	}

	err = s.repo.SavePlay(ctx, userID, trackID)
	if err != nil {
		return nil, err
	}

	log.Info("track played successfully",
		zap.Int64("user_id", userID),
		zap.Int64("track_id", trackID),
		zap.String("title", track.Title),
		zap.String("subscription_type", subscriptionType),
	)

	return track, nil
}
