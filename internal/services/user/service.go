package user

import (
	"ai-bot/internal/domain/models"
	"context"
)

type provider interface {
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

	User(
		ctx context.Context,
		uid int64,
	) (user *models.User, err error)
}

type owner interface {
	AuthorizeUser(
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

	BlockUser(
		ctx context.Context,
		uid int64,
	) (err error)

	UnblockUser(
		ctx context.Context,
		uid int64,
	) (err error)

	DeleteUser(
		ctx context.Context,
		uid int64,
	) (err error)

	SaveUser(
		ctx context.Context,
		user *models.User,
	) (err error)
}

type Service struct {
	owner    owner
	provider provider
}

func New(owner owner, provider provider) *Service {
	return &Service{
		owner:    owner,
		provider: provider,
	}
}
