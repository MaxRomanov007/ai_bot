package user

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
)

func (s *Service) User(
	ctx context.Context,
	uid int64,
) (*models.User, error) {
	const op = "services.user.User"

	user, err := s.provider.User(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get user: %w", op, HandleStorageError(err))
	}

	return user, nil
}
