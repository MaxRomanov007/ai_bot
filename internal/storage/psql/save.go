package psql

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
)

func (s *Storage) SaveUser(
	ctx context.Context,
	user *models.User,
) error {
	const op = "storage.psql.SaveUser"

	query := `
		INSERT INTO users(user_id, username, chat_id, role_id, status_id) 
		VALUES 
		(
			$1, $2, $3,
			(
				SELECT user_role_id
				FROM user_roles
				WHERE user_role_name = $4
			),
			(
				SELECT user_status_id
				FROM user_statuses
				WHERE user_status_name = $5
			)
		)`

	if _, err := s.db.ExecContext(ctx, query,
		user.UserID, user.Username, user.ChatID,
		models.UserRoleUser,
		models.UserStatusUnauthorized,
	); err != nil {
		return fmt.Errorf("%s: failed to insert user: %w", op, HandleDatabaseError(err))
	}

	return nil
}
