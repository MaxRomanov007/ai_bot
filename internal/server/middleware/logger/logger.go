package logger

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"time"
)

func Logger(log *slog.Logger) bot.Middleware {
	log = log.With(
		slog.String("component", "middleware/logger"),
	)
	log.Info("logger middleware enabled")
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, u *models.Update) {
			entry := log.With(
				slog.Int64("update_id", u.ID),
			)
			switch {
			case u.Message != nil:
				var text string
				runes := []rune(u.Message.Text)
				if len(runes) >= 50 {
					text = string(runes[:50]) + "..."
				} else {
					text = string(runes)
				}

				entry = entry.With(
					slog.String("type", "message"),
					slog.Int64("from_id", u.Message.From.ID),
					slog.String("from_username", u.Message.From.Username),
					slog.String("text", text),
				)
			case u.CallbackQuery != nil:
				entry = entry.With(
					slog.String("type", "callback_query"),
					slog.Int64("from_id", u.CallbackQuery.From.ID),
					slog.String("from_username", u.CallbackQuery.From.Username),
					slog.String("data", u.CallbackQuery.Data),
				)
			default:
				return
			}

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next(ctx, b, u)
		}
	}
}
