package botApp

import (
	"ai-bot/internal/config"
	"ai-bot/internal/server/aiBot"
	"ai-bot/internal/server/middleware/logger"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"log/slog"
)

type BotApp struct {
	Log   *slog.Logger
	tgBot *bot.Bot
	Cfg   *config.TelegramConfig
}

func New(log *slog.Logger, cfg *config.TelegramConfig, api *aiBot.Server) (*BotApp, error) {
	const op = "app.BotApp.New"

	opts := []bot.Option{
		bot.WithMiddlewares(
			api.WithDefaultAuthorization,
			logger.Logger(log),
		),
		bot.WithDefaultHandler(api.SendMessage()),
	}

	log = log.With(slog.String("operation", op))

	b, err := bot.New(cfg.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create bot: %w", op, err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, api.RestartChat())
	b.RegisterHandler(bot.HandlerTypeMessageText, "/admin", bot.MatchTypeExact, api.Admin())

	return &BotApp{
		Log:   log,
		tgBot: b,
		Cfg:   cfg,
	}, nil
}

func (a *BotApp) Run(ctx context.Context) {
	a.tgBot.Start(ctx)
}
