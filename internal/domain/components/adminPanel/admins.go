package adminPanel

import (
	"ai-bot/internal/domain/components/ui/paginator"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func adminsHandler(o owner, errFunc HandleErrorFunc) paginator.Handler {
	return func(ctx context.Context, b *bot.Bot, u *models.Update, npf paginator.NewPaginatorFunc, uid int64) {
		kb := inline.New(b, inline.NoDeleteAfterClick()).
			Row().
			Button("set user", []byte("set user"), adminsOnSelect(o, npf, errFunc, uid)).
			Row().
			Button("back", []byte("back"), adminsOnSelect(o, npf, errFunc, uid)).
			Button("exit", []byte("exit"), adminsOnSelect(o, npf, errFunc, uid))

		_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:          u.CallbackQuery.Message.Message.Chat.ID,
			MessageID:       u.CallbackQuery.Message.Message.ID,
			InlineMessageID: u.CallbackQuery.InlineMessageID,
			Text:            "choose the variant",
			ParseMode:       models.ParseModeMarkdown,
			ReplyMarkup:     kb,
		})
		if err != nil {
			errFunc(err)
			return
		}
	}
}

func adminsOnSelect(o owner, npf paginator.NewPaginatorFunc, errFunc HandleErrorFunc, uid int64) inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
		switch string(data) {
		case "set user":
			err := o.SetUserRole(ctx, uid)
			if err != nil {
				errFunc(err)
				return
			}
			pag := npf(b, -1).BuildSendParams(ctx, mes.Message.Chat.ID)
			GoToPage(ctx, b, mes.Message.Chat.ID, mes.Message.ID, pag.Text, pag.ReplyMarkup, errFunc)
		case "back":
			pag := npf(b, 0).BuildSendParams(ctx, mes.Message.Chat.ID)
			GoToPage(ctx, b, mes.Message.Chat.ID, mes.Message.ID, pag.Text, pag.ReplyMarkup, errFunc)
		case "exit":
			_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    mes.Message.Chat.ID,
				MessageID: mes.Message.ID,
			})
			if err != nil {
				errFunc(err)
				return
			}
		}
	}
}
