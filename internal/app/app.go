package app

import (
	"ai-bot/internal/app/botApp"
	"ai-bot/internal/config"
	"ai-bot/internal/server/aiBot"
	"ai-bot/internal/services/ai"
	"ai-bot/internal/services/user"
	"ai-bot/internal/storage/psql"
	"fmt"
	"log/slog"
)

type App struct {
	AIBot *botApp.BotApp
}

func MustLoad(
	log *slog.Logger,
	cfg config.Config,
) *App {
	app, err := New(log, cfg)
	if err != nil {
		panic(err)
	}

	return app
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	const op = "app.New"

	storage, err := psql.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create storage: %w", op, err)
	}

	aiService := ai.New(cfg.AI, storage, storage)
	userService := user.New(storage, storage)

	server := aiBot.New(log, cfg, userService, aiService)

	tgBot, err := botApp.New(log, cfg.Telegram, server)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create bot app: %w", op, err)
	}

	return &App{
		AIBot: tgBot,
	}, nil
}
