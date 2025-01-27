package user

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
)

func (s *Service) Save(
	ctx context.Context,
	user *models.User,
) error {
	const op = "services.user.Save"

	if err := s.owner.SaveUser(ctx, user); err != nil {
		return fmt.Errorf("%s: failed to save user: %w", op, HandleStorageError(err))
	}

	return nil
}
