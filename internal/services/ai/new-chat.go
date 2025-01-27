package ai

import (
	"context"
	"fmt"
)

func (s *Service) NewChat(
	ctx context.Context,
	uid int64,
) error {
	const op = "services.ai.NewChat"

	if err := s.owner.DeleteUserMessages(ctx, uid); err != nil {
		return fmt.Errorf("%s: failed to delete user messages: %w", op, HandleStorageError(err))
	}

	return nil
}
