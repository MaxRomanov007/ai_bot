package psql

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

func (s *Storage) Blocked(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.User, error) {
	const op = "storage.psql.Blocked"

	builder := s.builder.
		Select("users.*").
		From("users").
		Join("user_statuses ON users.status_id = user_statuses.user_status_id").
		Join("user_roles ON users.role_id = user_roles.user_role_id").
		Where(squirrel.And{
			squirrel.Eq{
				"user_statuses.user_status_name": models.UserStatusBlocked,
			},
			squirrel.Eq{
				"user_roles.user_role_name": models.UserRoleUser,
			},
		})

	if limit > 0 {
		builder = builder.Limit(uint64(limit))
	}
	if offset > 0 {
		builder = builder.Offset(uint64(offset))
	}

	stmt, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var users []models.User
	if err := s.db.SelectContext(ctx, &users, stmt, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to get users from pq: %w", op, HandleDatabaseError(err))
	}

	return users, nil
}

func (s *Storage) BlockedCount(
	ctx context.Context,
) (int, error) {
	const op = "storage.psql.BlockedCount"

	stmt, args, err := s.builder.
		Select("COUNT(*)").
		From("users").
		Join("user_statuses ON users.status_id = user_statuses.user_status_id").
		Join("user_roles ON users.role_id = user_roles.user_role_id").
		Where(squirrel.And{
			squirrel.Eq{
				"user_statuses.user_status_name": models.UserStatusBlocked,
			},
			squirrel.Eq{
				"user_roles.user_role_name": models.UserRoleUser,
			},
		}).ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var count int
	if err := s.db.GetContext(ctx, &count, stmt, args...); err != nil {
		return 0, fmt.Errorf("%s: failed to get users count from pq: %w", op, HandleDatabaseError(err))
	}

	return count, nil
}
