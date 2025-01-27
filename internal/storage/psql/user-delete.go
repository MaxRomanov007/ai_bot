package psql

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

func (s *Storage) DeleteUser(
	ctx context.Context,
	uid int64,
) error {
	const op = "storage.psql.Delete"

	stmt, args, err := s.builder.
		Delete("users").
		Where(squirrel.Eq{
			"user_id": uid,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build delete query: %w", op, err)
	}

	result, err := s.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to delete user: %w", op, HandleDatabaseError(err))
	}

	ar, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get affected rows: %w", op, err)
	}
	if ar == 0 {
		return fmt.Errorf("%s: user to delete not found: %w", op, ErrNotFound)
	}

	return nil
}
