package user

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
)

func (s *Service) Blocked(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.User, error) {
	const op = "services.user.Blocked"

	users, err := s.provider.Blocked(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get blocked users: %w", op, HandleStorageError(err))
	}

	return users, nil
}

func (s *Service) BlockedCount(
	ctx context.Context,
) (int, error) {
	const op = "services.user.BlockedCount"

	count, err := s.provider.BlockedCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get blocked users count: %w", op, HandleStorageError(err))
	}

	return count, nil
}
