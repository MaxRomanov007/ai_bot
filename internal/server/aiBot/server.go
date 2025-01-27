package aiBot

import (
	"ai-bot/internal/config"
	"ai-bot/internal/domain/models"
	"context"
	"log/slog"
)

type UserService interface {
	Unauthorized(
		ctx context.Context,
		limit int,
		offset int,
	) (users []models.User, err error)

	UnauthorizedCount(
		ctx context.Context,
	) (count int, err error)

	Users(
		ctx context.Context,
		limit int,
		offset int,
	) (users []models.User, err error)

	UsersCount(
		ctx context.Context,
	) (count int, err error)

	Blocked(
		ctx context.Context,
		limit int,
		offset int,
	) (users []models.User, err error)

	BlockedCount(
		ctx context.Context,
	) (count int, err error)

	Admins(
		ctx context.Context,
		limit int,
		offset int,
	) (users []models.User, err error)

	AdminsCount(
		ctx context.Context,
	) (count int, err error)

	Authorize(
		ctx context.Context,
		uid int64,
	) (err error)

	SetAdminRole(
		ctx context.Context,
		uid int64,
	) (err error)

	SetUserRole(
		ctx context.Context,
		uid int64,
	) (err error)

	Block(
		ctx context.Context,
		uid int64,
	) (err error)

	Unblock(
		ctx context.Context,
		uid int64,
	) (err error)

	Delete(
		ctx context.Context,
		uid int64,
	) (err error)

	Save(
		ctx context.Context,
		user *models.User,
	) (err error)

	User(
		ctx context.Context,
		uid int64,
	) (user *models.User, err error)
}

type AIService interface {
	SendMessage(
		ctx context.Context,
		uid int64,
		message string,
	) (content string, err error)

	NewChat(
		ctx context.Context,
		uid int64,
	) (err error)
}

type Server struct {
	Log         *slog.Logger
	Cfg         config.Config
	UserService UserService
	AIService   AIService
}

func New(log *slog.Logger, cfg config.Config, us UserService, as AIService) *Server {
	return &Server{
		Log:         log,
		Cfg:         cfg,
		UserService: us,
		AIService:   as,
	}
}
