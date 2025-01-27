package ai

import (
	"ai-bot/internal/config"
	"ai-bot/internal/domain/models"
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type provider interface {
	UserMessages(
		ctx context.Context,
		uid int64,
	) (messages []models.Message, err error)
}

type owner interface {
	DeleteUserMessages(
		ctx context.Context,
		uid int64,
	) (err error)

	SaveMessage(
		ctx context.Context,
		message models.Message,
	) (err error)
}

type Service struct {
	client   *openai.Client
	cfg      *config.AIConfig
	provider provider
	owner    owner
}

func New(
	cfg *config.AIConfig,
	provider provider,
	owner owner,
) *Service {

	client := openai.NewClient(
		option.WithBaseURL(cfg.BaseURL),
		option.WithAPIKey(cfg.Key),
	)

	return &Service{
		client:   client,
		cfg:      cfg,
		provider: provider,
		owner:    owner,
	}
}

func (s *Service) buildAIMessage(messages []models.Message) openai.ChatCompletionNewParams {
	aiMessages := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	for i := 0; i < len(messages); i++ {
		switch messages[i].Role.MessageRoleName {
		case models.MessageRoleUser:
			aiMessages[i] = openai.UserMessage(messages[i].Content)
		case models.MessageRoleAssistant:
			aiMessages[i] = openai.AssistantMessage(messages[i].Content)
		}
	}
	if s.cfg.Prompt != "" {
		aiMessages = append([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(s.cfg.Prompt),
		}, aiMessages...)
	}

	return openai.ChatCompletionNewParams{
		Messages:            openai.F(aiMessages),
		Model:               openai.F(s.cfg.Model),
		MaxCompletionTokens: openai.F(s.cfg.MaxCompletionTokens),
	}
}
