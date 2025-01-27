package ai

import (
	"ai-bot/internal/domain/models"
	"context"
	"fmt"
)

func (s *Service) SendMessage(
	ctx context.Context,
	uid int64,
	message string,
) (string, error) {
	const op = "services.ai.SendMessage"

	messages, err := s.provider.UserMessages(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("%s: failed to get user messages: %w", op, HandleStorageError(err))
	}

	mes := models.Message{
		UserID:  uid,
		Content: message,
		Role: models.MessageRole{
			MessageRoleName: models.MessageRoleUser,
		}}

	if err := s.owner.SaveMessage(ctx, mes); err != nil {
		return "", fmt.Errorf("%s: failed to save user message: %w", op, HandleStorageError(err))
	}

	messages = append(messages, mes)

	chatCompletion, err := s.client.Chat.Completions.New(ctx, s.buildAIMessage(messages))
	if err != nil {
		return "", fmt.Errorf("%s: failed to get ai answer: %w", op, err)
	}

	if err := s.owner.SaveMessage(ctx, models.Message{
		UserID:  uid,
		Content: chatCompletion.Choices[0].Message.Content,
		Role: models.MessageRole{
			MessageRoleName: models.MessageRoleAssistant,
		},
	}); err != nil {
		return "", fmt.Errorf("%s: failed to save assistant message: %w", op, HandleStorageError(err))
	}

	return chatCompletion.Choices[0].Message.Content, nil
}
