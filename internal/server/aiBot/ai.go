package aiBot

import (
	"ai-bot/internal/lib/api/logger/sl"
	"ai-bot/internal/lib/response"
	"ai-bot/internal/services/ai"
	"context"
	"errors"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
)

func (s *Server) SendMessage() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, u *models.Update) {
		const op = "server.aiBot.SendMessage"

		log := s.Log.With(
			slog.String("operation", op),
			slog.Int64("update_id", u.ID),
		)

		resp, err := s.AIService.SendMessage(ctx, u.Message.From.ID, u.Message.Text)
		if err != nil {
			var aiErr *ai.Error
			if errors.As(err, &aiErr) {
				log.Warn("ai error", sl.Err(aiErr))
				response.AIError(ctx, b, u, aiErr)
				return
			}

			log.Error("failed to send message to ai", sl.Err(err))
			response.Internal(ctx, b, u)
			return
		}

		response.SendText(ctx, b, u, resp)
	}
}

func (s *Server) RestartChat() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, u *models.Update) {
		const op = "server.aiBot.RestartChat"

		log := s.Log.With(
			slog.String("operation", op),
			slog.Int64("update_id", u.ID),
		)

		if err := s.AIService.NewChat(ctx, u.Message.From.ID); err != nil && !errors.Is(err, ai.ErrNotFound) {
			var aiErr *ai.Error
			if errors.As(err, &aiErr) {
				log.Warn("ai error", sl.Err(err))
				response.AIError(ctx, b, u, aiErr)
				return
			}

			log.Error("fail to start chat", sl.Err(err))
			response.Internal(ctx, b, u)
			return
		}

		response.SendText(ctx, b, u, "new chat started")
	}
}
