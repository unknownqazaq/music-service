package service

import (
	"context"

	"go.uber.org/zap"
	"music-service/internal/core/logger"
)

func (s *UserService) UpdateSubscription(ctx context.Context, id int64, subscriptionType string) error {
	log := logger.FromContext(ctx)

	err := s.repo.UpdateSubscription(ctx, id, subscriptionType)
	if err != nil {
		return err
	}

	log.Info("admin action: updated user subscription",
		zap.Int64("user_id", id),
		zap.String("new_subscription_type", subscriptionType),
	)

	return nil
}

