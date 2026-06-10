package service

import (
	"context"

	"go.uber.org/zap"
	"music-service/internal/core/logger"
)

func (s *FavoritesService) Add(ctx context.Context, userID, trackID int64, subscriptionType string) error {
	log := logger.FromContext(ctx)

	if subscriptionType == "FREE" {
		count, err := s.repo.CountByUserID(ctx, userID)
		if err != nil {
			return err
		}
		if count >= s.favoriteLimit {
			log.Warn("FREE subscription favorites limit exceeded",
				zap.Int64("user_id", userID),
				zap.Int64("track_id", trackID),
				zap.Int("current_favorites", count),
				zap.Int("limit", s.favoriteLimit),
			)
			return ErrFavoritesLimitExceeded
		}
	}

	err := s.repo.Add(ctx, userID, trackID)
	if err != nil {
		return err
	}

	log.Info("track added to favorites",
		zap.Int64("user_id", userID),
		zap.Int64("track_id", trackID),
		zap.String("subscription_type", subscriptionType),
	)

	return nil
}

