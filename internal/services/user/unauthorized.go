package user

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
)

func (s *Service) Unauthorized(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.User, error) {
	const op = "services.user.Unauthorized"

	users, err := s.provider.Unauthorized(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get unauthorized users: %w", op, HandleStorageError(err))
	}

	return users, nil
}

func (s *Service) UnauthorizedCount(
	ctx context.Context,
) (int, error) {
	const op = "services.user.UnauthorizedCount"

	count, err := s.provider.UnauthorizedCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get unauthorized users count: %w", op, HandleStorageError(err))
	}

	return count, nil
}
