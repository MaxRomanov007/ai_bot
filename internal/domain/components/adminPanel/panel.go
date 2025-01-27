package adminPanel

import (
	"ai-bot/internal/domain/components/ui/paginator"
	models2 "ai-bot/internal/domain/models"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

type provider interface {
	Unauthorized(ctx context.Context, limit int, offset int) (users []models2.User, err error)
	UnauthorizedCount(ctx context.Context) (count int, err error)
	Users(ctx context.Context, limit int, offset int) (users []models2.User, err error)
	UsersCount(ctx context.Context) (count int, err error)
	Blocked(ctx context.Context, limit int, offset int) (users []models2.User, err error)
	BlockedCount(ctx context.Context) (count int, err error)
	Admins(ctx context.Context, limit int, offset int) (users []models2.User, err error)
	AdminsCount(ctx context.Context) (count int, err error)
}

type owner interface {
	Authorize(ctx context.Context, uid int64) (err error)
	SetAdminRole(ctx context.Context, uid int64) (err error)
	SetUserRole(ctx context.Context, uid int64) (err error)
	Block(ctx context.Context, uid int64) (err error)
	Unblock(ctx context.Context, uid int64) (err error)
	Delete(ctx context.Context, uid int64) (err error)
}

type HandleErrorFunc func(err error)

type AdminPanel struct {
	isOwner   bool
	kb        *inline.Keyboard
	errorFunc HandleErrorFunc
	provider  provider
	owner     owner
}

func New(b *bot.Bot, isOwner bool, errHandle HandleErrorFunc, p provider, o owner) *AdminPanel {
	kb := inline.New(b, inline.NoDeleteAfterClick())
	if isOwner {
		kb = kb.
			Row().
			Button("admins", []byte("admins"), onSelect(isOwner, errHandle, p, o))
	}
	kb = kb.
		Row().
		Button("unauthorized", []byte("unauthorized"), onSelect(isOwner, errHandle, p, o)).
		Row().
		Button("users", []byte("users"), onSelect(isOwner, errHandle, p, o)).
		Row().
		Button("blocked", []byte("blocked"), onSelect(isOwner, errHandle, p, o)).
		Row().
		Button("exit", []byte("exit"), onSelect(isOwner, errHandle, p, o))

	return &AdminPanel{
		isOwner:   isOwner,
		kb:        kb,
		errorFunc: errHandle,
		provider:  p,
		owner:     o,
	}
}

func (ap *AdminPanel) Show(ctx context.Context, b *bot.Bot, chatID any) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "choose the variant",
		ReplyMarkup: ap.kb,
	})
	if err != nil {
		ap.errorFunc(err)
		return
	}
}

func onSelect(isOwner bool, errFunc HandleErrorFunc, p provider, o owner) inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		pag := &paginator.Paginator{}

		options := []paginator.Option{
			paginator.PerPage(3),
			paginator.NoDeleteBeforeHandler(),
			paginator.OnError(paginator.OnErrorHandler(errFunc)),
		}
		switch string(data) {
		case "admins":
			options = append(options, paginator.WithText("select admin"))
			count, err := p.AdminsCount(ctx)
			if err != nil {
				errFunc(err)
				return
			}
			pag = paginator.NewPaginator(b, p.Admins, count, adminsHandler(o, errFunc), options...)
		case "unauthorized":
			options = append(options, paginator.WithText("select unauthorized user"))
			count, err := p.UnauthorizedCount(ctx)
			if err != nil {
				errFunc(err)
				return
			}
			pag = paginator.NewPaginator(b, p.Unauthorized, count, unauthorizedHandler(o, errFunc), options...)
		case "users":
			options = append(options, paginator.WithText("select user"))
			count, err := p.UsersCount(ctx)
			if err != nil {
				errFunc(err)
				return
			}
			pag = paginator.NewPaginator(b, p.Users, count, usersHandler(o, errFunc, isOwner), options...)
		case "blocked":
			options = append(options, paginator.WithText("select blocked user"))
			count, err := p.BlockedCount(ctx)
			if err != nil {
				errFunc(err)
				return
			}
			pag = paginator.NewPaginator(b, p.Blocked, count, blockedHandler(o, errFunc), options...)
		case "exit":
			_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    mes.Message.Chat.ID,
				MessageID: mes.Message.ID,
			})
			if err != nil {
				errFunc(err)
				return
			}
			return
		}

		pag.SetBack("back", &bot.SendMessageParams{
			ChatID:      mes.Message.Chat.ID,
			Text:        "choose the user",
			ReplyMarkup: New(b, isOwner, errFunc, p, o).kb,
		})

		send := pag.BuildSendParams(ctx, mes.Message.Chat.ID)

		_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      send.ChatID,
			MessageID:   mes.Message.ID,
			Text:        send.Text,
			ParseMode:   send.ParseMode,
			ReplyMarkup: send.ReplyMarkup,
		})
		if err != nil {
			errFunc(err)
			return
		}
	}
}

func (ap *AdminPanel) GoBack(ctx context.Context, b *bot.Bot, chatID any) {
	ap.Show(ctx, b, chatID)
}

func GoToPage(ctx context.Context, b *bot.Bot, chatID any, messageID int, text string, rm models.ReplyMarkup, errFunc HandleErrorFunc) {
	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: rm,
	})
	if err != nil {
		errFunc(err)
		return
	}
}
