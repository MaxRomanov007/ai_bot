package psql

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
)

func (s *Storage) AuthorizeUser(
	ctx context.Context,
	uid int64,
) error {
	const op = "storage.psql.AuthorizeUser"

	query := `
		UPDATE users
        SET status_id = (
            SELECT user_status_id
            FROM user_statuses
            WHERE user_status_name = $1
        )
        WHERE user_id = $2`

	result, err := s.db.ExecContext(ctx, query, models.UserStatusAuthorized, uid)
	if err != nil {
		return fmt.Errorf("%s: failed to update user: %w", op, HandleDatabaseError(err))
	}

	ar, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get affected rows: %w", op, err)
	}
	if ar == 0 {
		return fmt.Errorf("%s: user to update not found: %w", op, ErrNotFound)
	}

	return nil
}

func (s *Storage) SetAdminRole(
	ctx context.Context,
	uid int64,
) error {
	const op = "storage.psql.SetAdminRole"

	query := `
		UPDATE users
        SET role_id = (
            SELECT user_role_id
            FROM user_roles
            WHERE user_role_name = $1
        )
        WHERE user_id = $2`

	result, err := s.db.ExecContext(ctx, query, models.UserRoleAdmin, uid)
	if err != nil {
		return fmt.Errorf("%s: failed to update user: %w", op, HandleDatabaseError(err))
	}

	ar, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get affected rows: %w", op, err)
	}
	if ar == 0 {
		return fmt.Errorf("%s: user to update not found: %w", op, ErrNotFound)
	}

	return nil
}

func (s *Storage) SetUserRole(
	ctx context.Context,
	uid int64,
) error {
	const op = "storage.psql.SetUserRole"

	query := `
		UPDATE users
        SET role_id = (
            SELECT user_role_id
            FROM user_roles
            WHERE user_role_name = $1
        )
        WHERE user_id = $2`

	result, err := s.db.ExecContext(ctx, query, models.UserRoleUser, uid)
	if err != nil {
		return fmt.Errorf("%s: failed to update user: %w", op, HandleDatabaseError(err))
	}

	ar, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get affected rows: %w", op, err)
	}
	if ar == 0 {
		return fmt.Errorf("%s: user to update not found: %w", op, ErrNotFound)
	}

	return nil
}

func (s *Storage) BlockUser(
	ctx context.Context,
	uid int64,
) error {
	const op = "storage.psql.BlockUser"

	query := `
		UPDATE users
        SET status_id = (
            SELECT user_status_id
            FROM user_statuses
            WHERE user_status_name = $1
        )
        WHERE user_id = $2`

	result, err := s.db.ExecContext(ctx, query, models.UserStatusBlocked, uid)
	if err != nil {
		return fmt.Errorf("%s: failed to update user: %w", op, HandleDatabaseError(err))
	}

	ar, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get affected rows: %w", op, err)
	}
	if ar == 0 {
		return fmt.Errorf("%s: user to update not found: %w", op, ErrNotFound)
	}

	return nil
}

func (s *Storage) UnblockUser(
	ctx context.Context,
	uid int64,
) error {
	const op = "storage.psql.UnblockUser"

	query := `
		UPDATE users
        SET status_id = (
            SELECT user_status_id
            FROM user_statuses
            WHERE user_status_name = $1
        )
        WHERE user_id = $2`

	result, err := s.db.ExecContext(ctx, query, models.UserStatusAuthorized, uid)
	if err != nil {
		return fmt.Errorf("%s: failed to update user: %w", op, HandleDatabaseError(err))
	}

	ar, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get affected rows: %w", op, err)
	}
	if ar == 0 {
		return fmt.Errorf("%s: user to update not found: %w", op, ErrNotFound)
	}

	return nil
}
