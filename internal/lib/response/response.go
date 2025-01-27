package response

import (
	"ai-bot/internal/services/ai"
	"ai-bot/internal/services/user"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SendText(ctx context.Context, b *bot.Bot, u *models.Update, text string) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   text,
	})
}

func SendRepliedText(ctx context.Context, b *bot.Bot, u *models.Update, text string) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   text,
		ReplyParameters: &models.ReplyParameters{
			MessageID: u.Message.ID,
		},
	})
}

func AIError(ctx context.Context, b *bot.Bot, u *models.Update, err *ai.Error) {
	var text string
	switch err.Code {
	case user.ErrAlreadyExistsCode:
		text = "message already exists"
	case user.ErrNotFoundCode:
		text = "message not found"
	default:
		text = "internal error"
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   text,
	})
}

func UserError(ctx context.Context, b *bot.Bot, u *models.Update, err *user.Error) {
	var text string
	switch err.Code {
	case user.ErrAlreadyExistsCode:
		text = "user already exists"
	case user.ErrNotFoundCode:
		text = "user not found"
	default:
		text = "internal error"
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   text,
	})
}

func Internal(ctx context.Context, b *bot.Bot, u *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   "internal error",
	})
}

func AccessDenied(ctx context.Context, b *bot.Bot, u *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   "access denied",
	})
}

func Unauthorized(ctx context.Context, b *bot.Bot, u *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   "you haven't be authorized",
	})
}

func Blocked(ctx context.Context, b *bot.Bot, u *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   "you've been blocked",
	})
}
