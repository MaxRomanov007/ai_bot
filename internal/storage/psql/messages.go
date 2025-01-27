package psql

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

func (s *Storage) UserMessages(
	ctx context.Context,
	uid int64,
) ([]models.Message, error) {
	const op = "storage.psql.UserMessages"

	stmt, args, err := s.builder.
		Select(
			"messages.*",
			"message_roles.*",
		).
		From("messages").
		Join("message_roles ON messages.role_id = message_roles.message_role_id").
		Where(squirrel.Eq{
			"messages.user_id": uid,
		}).OrderBy("messages.message_id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var dbMessages []struct {
		models.Message
		models.MessageRole
	}
	if err := s.db.SelectContext(ctx, &dbMessages, stmt, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to get messages from database: %w", op, HandleDatabaseError(err))
	}

	if dbMessages == nil {
		return []models.Message{}, nil
	}

	result := make([]models.Message, len(dbMessages))
	for i := 0; i < len(dbMessages); i++ {
		result[i] = models.Message{
			MessageID: dbMessages[i].MessageID,
			UserID:    dbMessages[i].UserID,
			RoleID:    dbMessages[i].RoleID,
			Content:   dbMessages[i].Content,
			Role: models.MessageRole{
				MessageRoleID:   dbMessages[i].MessageRoleID,
				MessageRoleName: dbMessages[i].MessageRoleName,
			},
		}
	}

	return result, nil
}

func (s *Storage) DeleteUserMessages(
	ctx context.Context,
	uid int64,
) error {
	const op = "storage.psql.DeleteUserMessages"

	stmt, args, err := s.builder.
		Delete("messages").
		Where(squirrel.Eq{
			"user_id": uid,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	result, err := s.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to delete user messages: %w", op, HandleDatabaseError(err))
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

func (s *Storage) SaveMessage(
	ctx context.Context,
	message models.Message,
) error {
	const op = "storage.psql.SaveMessage"

	query := `
		INSERT INTO messages(user_id, content, role_id)
		VALUES
		(
			$1, $2,
			(
				SELECT message_role_id
				FROM message_roles
				WHERE message_role_name = $3
			)
		)`

	if _, err := s.db.ExecContext(ctx, query,
		message.UserID, message.Content,
		message.Role.MessageRoleName,
	); err != nil {
		return fmt.Errorf("%s: failed to insert message: %w", op, HandleDatabaseError(err))
	}

	return nil
}
