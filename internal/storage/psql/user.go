package psql

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

func (s *Storage) User(
	ctx context.Context,
	uid int64,
) (*models.User, error) {
	const op = "storage.psql.User"

	stmt, args, err := s.builder.
		Select(
			"users.*",
			"user_statuses.*",
			"user_roles.*",
		).
		From("users").
		Join("user_roles ON user_roles.user_role_id = users.role_id").
		Join("user_statuses ON user_statuses.user_status_id = users.status_id").
		Where(squirrel.Eq{
			"users.user_id": uid,
		}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build delete query: %w", op, err)
	}

	var user struct {
		models.User
		models.UserStatus
		models.UserRole
	}
	if err := s.db.GetContext(ctx, &user, stmt, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to get user: %w", op, HandleDatabaseError(err))
	}

	return &models.User{
		UserID:   user.UserID,
		Username: user.Username,
		ChatID:   user.ChatID,
		RoleID:   user.RoleID,
		StatusID: user.StatusID,
		Role: models.UserRole{
			UserRoleID:   user.UserRoleID,
			UserRoleName: user.UserRoleName,
		},
		Status: models.UserStatus{
			UserStatusID:   user.UserStatusID,
			UserStatusName: user.UserStatusName,
		},
	}, nil
}
