package user

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
)

func (s *Service) Users(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.User, error) {
	const op = "services.user.Users"

	users, err := s.provider.Users(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get users: %w", op, HandleStorageError(err))
	}

	return users, nil
}

func (s *Service) UsersCount(
	ctx context.Context,
) (int, error) {
	const op = "services.user.UsersCount"

	count, err := s.provider.UsersCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get users count: %w", op, HandleStorageError(err))
	}

	return count, nil
}
