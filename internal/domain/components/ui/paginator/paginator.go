package paginator

import (
	"context"
	"log"

	models2 "ai-bot/internal/domain/models"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type OnErrorHandler func(err error)

const (
	cmdNop   = "nop"
	cmdStart = "start"
	cmdEnd   = "end"
	cmdClose = "close"
	cmdUser  = "user"
	cmdBack  = "back"
)

type NewPaginatorFunc func(b *bot.Bot, delta int) *Paginator
type Handler func(ctx context.Context, bot *bot.Bot, update *models.Update, newPaginatorFunc NewPaginatorFunc, uid int64)
type GetDataFunc func(ctx context.Context, limit int, offset int) (users []models2.User, err error)

type Paginator struct {
	Text                  string
	getDataFunc           GetDataFunc
	handler               Handler
	totalCount            int
	perPage               int
	currentPage           int
	pagesCount            int
	prefix                string
	closeButton           string
	backButton            string
	isDeleteBeforeHandler bool
	backMessage           *bot.SendMessageParams
	onError               OnErrorHandler

	callbackHandlerID string
}

func NewPaginator(b *bot.Bot, gdf GetDataFunc, totalCount int, handler Handler, opts ...Option) *Paginator {
	p := &Paginator{
		prefix:                bot.RandomString(16),
		getDataFunc:           gdf,
		currentPage:           1,
		perPage:               10,
		totalCount:            totalCount,
		isDeleteBeforeHandler: true,
		onError:               defaultOnError,
		handler:               handler,
	}

	for _, opt := range opts {
		opt(p)
	}

	if p.Text == "" {
		p.Text = "select:"
	}

	p.pagesCount = totalCount / p.perPage
	if totalCount%p.perPage != 0 {
		p.pagesCount++
	}

	p.callbackHandlerID = b.RegisterHandler(bot.HandlerTypeCallbackQueryData, p.prefix, bot.MatchTypePrefix, p.callback)

	return p
}

// Prefix returns the prefix of the widget
func (p *Paginator) Prefix() string {
	return p.prefix
}

func defaultOnError(err error) {
	log.Printf("[TG-UI-PAGINATOR] [ERROR] %s", err)
}

func (p *Paginator) Show(ctx context.Context, b *bot.Bot, chatID any, opts ...ShowOption) (*models.Message, error) {
	params := p.BuildSendParams(ctx, chatID, opts...)

	return b.SendMessage(ctx, params)
}

func (p *Paginator) BuildSendParams(ctx context.Context, chatID any, opts ...ShowOption) *bot.SendMessageParams {
	data, err := p.getDataFunc(ctx, p.perPage, p.perPage*(p.currentPage-1))
	if err != nil {
		p.onError(err)
	}

	params := &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        p.Text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: p.buildMarkup(data),
	}

	for _, o := range opts {
		o(params)
	}

	return params
}

func (p *Paginator) SetBack(text string, prev *bot.SendMessageParams) {
	p.backButton = text
	p.backMessage = prev
}
