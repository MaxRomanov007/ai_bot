package user

import (
	"context"
	"fmt"
)

func (s *Service) Authorize(
	ctx context.Context,
	uid int64,
) (err error) {
	const op = "services.user.Authorize"

	if err := s.owner.AuthorizeUser(ctx, uid); err != nil {
		return fmt.Errorf("%s: failed to authorize user: %w", op, HandleStorageError(err))
	}

	return nil
}

func (s *Service) SetAdminRole(
	ctx context.Context,
	uid int64,
) (err error) {
	const op = "services.user.SetAdminRole"

	if err := s.owner.SetAdminRole(ctx, uid); err != nil {
		return fmt.Errorf("%s: failed to set user role to user: %w", op, HandleStorageError(err))
	}

	return nil
}

func (s *Service) SetUserRole(
	ctx context.Context,
	uid int64,
) (err error) {
	const op = "services.user.SetUserRole"

	if err := s.owner.SetUserRole(ctx, uid); err != nil {
		return fmt.Errorf("%s: failed to set admin role to user: %w", op, HandleStorageError(err))
	}

	return nil
}

func (s *Service) Block(
	ctx context.Context,
	uid int64,
) (err error) {
	const op = "services.user.Block"

	if err := s.owner.BlockUser(ctx, uid); err != nil {
		return fmt.Errorf("%s: failed to block user: %w", op, HandleStorageError(err))
	}

	return nil
}

func (s *Service) Unblock(
	ctx context.Context,
	uid int64,
) (err error) {
	const op = "services.user.Unblock"

	if err := s.owner.UnblockUser(ctx, uid); err != nil {
		return fmt.Errorf("%s: failed to unblock user: %w", op, HandleStorageError(err))
	}

	return nil
}

func (s *Service) Delete(
	ctx context.Context,
	uid int64,
) (err error) {
	const op = "services.user.Delete"

	if err := s.owner.DeleteUser(ctx, uid); err != nil {
		return fmt.Errorf("%s: failed to delete user: %w", op, HandleStorageError(err))
	}

	return nil
}
