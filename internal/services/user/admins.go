package user

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
)

func (s *Service) Admins(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.User, error) {
	const op = "services.user.Admins"

	users, err := s.provider.Admins(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get admins: %w", op, HandleStorageError(err))
	}

	return users, nil
}

func (s *Service) AdminsCount(
	ctx context.Context,
) (int, error) {
	const op = "services.user.AdminsCount"

	count, err := s.provider.AdminsCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get admins count: %w", op, HandleStorageError(err))
	}

	return count, nil
}
