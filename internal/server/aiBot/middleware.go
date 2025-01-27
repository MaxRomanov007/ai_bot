package aiBot

import (
	models2 "ai-bot/internal/domain/models"
	"ai-bot/internal/lib/api/logger/sl"
	"ai-bot/internal/lib/response"
	"ai-bot/internal/services/user"
	"context"
	"errors"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
)

func (s *Server) WithDefaultAuthorization(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, u *models.Update) {
		const op = "server.aiBot.WithDefaultAuthorization"

		if u.Message == nil {
			next(ctx, b, u)
			return
		}

		log := s.Log.With(
			slog.String("operation", op),
			slog.Int64("update_id", u.ID),
		)

		if u.Message.Text == "/admin" && u.Message.From.ID == s.Cfg.Telegram.OwnerTelegramID {
			next(ctx, b, u)
			return
		}

		us, err := s.UserService.User(ctx, u.Message.From.ID)
		if errors.Is(err, user.ErrNotFound) {
			if err := s.UserService.Save(ctx, &models2.User{
				UserID:   u.Message.From.ID,
				Username: u.Message.From.Username,
				ChatID:   u.Message.Chat.ID,
			}); err != nil {
				var userErr *user.Error
				if errors.As(err, &userErr) {
					log.Warn("user error while saving user", sl.Err(err))
					response.UserError(ctx, b, u, userErr)
					return
				}

				log.Error("failed to save new user", sl.Err(err))
				response.Internal(ctx, b, u)
				return
			}
			response.SendText(ctx, b, u, "authorization request has been sent")
			return
		}
		if err != nil {
			var userErr *user.Error
			if errors.As(err, &userErr) {
				log.Warn("user error while getting user", sl.Err(err))
				response.UserError(ctx, b, u, userErr)
				return
			}

			log.Error("failed to get user", sl.Err(err))
			response.Internal(ctx, b, u)
			return
		}

		if us.Status.UserStatusName == models2.UserStatusUnauthorized {
			log.Warn("user is unauthorized")
			response.Unauthorized(ctx, b, u)
			return
		}
		if us.Status.UserStatusName == models2.UserStatusBlocked {
			log.Warn("user blocked")
			response.Blocked(ctx, b, u)
			return
		}

		next(ctx, b, u)
	}
}
