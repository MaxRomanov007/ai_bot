package aiBot

import (
	"ai-bot/internal/domain/components/adminPanel"
	models2 "ai-bot/internal/domain/models"
	"ai-bot/internal/lib/api/logger/sl"
	"ai-bot/internal/lib/response"
	"ai-bot/internal/services/ai"
	"ai-bot/internal/services/user"
	"context"
	"errors"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
)

func (s *Server) Admin() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, u *models.Update) {
		const op = "server.aiBot.Admin"

		log := s.Log.With(
			slog.String("operation", op),
			slog.Int64("update_id", u.ID),
		)

		isOwner := u.Message.From.ID == s.Cfg.Telegram.OwnerTelegramID
		if !isOwner && !s.isAdmin(ctx, u.Message.From.ID, log) {
			log.Warn("access denied")
			response.AccessDenied(ctx, b, u)
			return
		}

		ap := adminPanel.New(b, isOwner, handleError(ctx, log, b, u), s.UserService, s.UserService)
		ap.Show(ctx, b, u.Message.Chat.ID)
	}
}

func (s *Server) isAdmin(ctx context.Context, uid int64, log *slog.Logger) bool {
	u, err := s.UserService.User(ctx, uid)
	if err != nil {
		var userErr *user.Error
		if errors.As(err, &userErr) {
			log.Warn("user error", sl.Err(userErr))
			return false
		}

		log.Error("failed to get user role", sl.Err(err))
		return false
	}

	if u.Status.UserStatusName != models2.UserStatusAuthorized {
		log.Warn("user have not authorize status")
		return false
	}

	if u.Role.UserRoleName != models2.UserRoleAdmin {
		log.Warn("user is not admin")
		return false
	}

	return true
}

func handleError(ctx context.Context, log *slog.Logger, b *bot.Bot, u *models.Update) adminPanel.HandleErrorFunc {
	return func(err error) {
		var userErr *user.Error
		if errors.As(err, &userErr) {
			log.Warn("user error", sl.Err(err))
			response.UserError(ctx, b, u, userErr)
			return
		}

		var AIErr *ai.Error
		if errors.As(err, &AIErr) {
			log.Warn("ai error", sl.Err(err))
			response.AIError(ctx, b, u, AIErr)
			return
		}

		log.Error("unknown error", sl.Err(err))
		response.Internal(ctx, b, u)
	}
}
