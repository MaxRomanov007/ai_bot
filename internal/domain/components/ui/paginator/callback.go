package paginator

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (p *Paginator) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		p.onError(err)
		return
	}
	if !ok {
		p.onError(fmt.Errorf("callback answer failed"))
	}
}

func (p *Paginator) callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	cmd := strings.TrimPrefix(update.CallbackQuery.Data, p.prefix)

	switch getCmd(cmd) {
	case cmdUser:
		p.callbackAnswer(ctx, b, update.CallbackQuery)

		uid, err := strconv.ParseInt(getUID(cmd), 16, 64)
		if err != nil {
			p.onError(err)
			return
		}

		if p.isDeleteBeforeHandler {
			p.close(ctx, b, update)
		} else {
			b.UnregisterHandler(p.callbackHandlerID)
		}
		p.handler(ctx, b, update, p.newPaginator, uid)
		p.callbackAnswer(ctx, b, update.CallbackQuery)

		return
	case cmdNop:
		p.callbackAnswer(ctx, b, update.CallbackQuery)
		return
	case cmdStart:
		if p.currentPage == 1 {
			p.callbackAnswer(ctx, b, update.CallbackQuery)
			return
		}
		p.currentPage = 1
	case cmdEnd:
		if p.currentPage == p.pagesCount {
			p.callbackAnswer(ctx, b, update.CallbackQuery)
			return
		}
		p.currentPage = p.pagesCount
	case cmdClose:
		p.close(ctx, b, update)
		p.callbackAnswer(ctx, b, update.CallbackQuery)
		return
	case cmdBack:
		b.UnregisterHandler(p.callbackHandlerID)
		_, errEdit := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:          update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:       update.CallbackQuery.Message.Message.ID,
			InlineMessageID: update.CallbackQuery.InlineMessageID,
			Text:            p.backMessage.Text,
			ParseMode:       models.ParseModeMarkdown,
			ReplyMarkup:     p.backMessage.ReplyMarkup,
		})
		if errEdit != nil {
			p.onError(errEdit)
		}
		p.callbackAnswer(ctx, b, update.CallbackQuery)
		return
	default:
		page, _ := strconv.Atoi(cmd)
		p.currentPage = page
	}

	data, err := p.getDataFunc(ctx, p.perPage, p.perPage*(p.currentPage-1))
	if err != nil {
		p.onError(err)
	}

	_, errEdit := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:          update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:       update.CallbackQuery.Message.Message.ID,
		InlineMessageID: update.CallbackQuery.InlineMessageID,
		Text:            p.Text,
		ParseMode:       models.ParseModeMarkdown,
		ReplyMarkup:     p.buildMarkup(data),
	})
	if errEdit != nil {
		p.onError(errEdit)
	}

	p.callbackAnswer(ctx, b, update.CallbackQuery)
}

func (p *Paginator) close(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.UnregisterHandler(p.callbackHandlerID)

	_, errDelete := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})
	if errDelete != nil {
		p.onError(errDelete)
	}
}

func getCmd(s string) string {
	idx := strings.Index(s, ".")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func getUID(s string) string {
	idx := strings.Index(s, ".")
	if idx == -1 {
		return s
	}
	return s[idx+1:]
}

func (p *Paginator) newPaginator(b *bot.Bot, delta int) *Paginator {
	p.totalCount += delta
	opts := []Option{
		WithText(p.Text),
		PerPage(p.perPage),
	}
	if !p.isDeleteBeforeHandler {
		opts = append(opts, NoDeleteBeforeHandler())
	}
	pag := NewPaginator(b, p.getDataFunc, p.totalCount, p.handler, opts...)
	pag.SetBack(p.backButton, p.backMessage)
	return pag
}
