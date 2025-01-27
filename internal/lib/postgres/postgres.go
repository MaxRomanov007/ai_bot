package postgres

import (
	"ai-bot/internal/config"
	"fmt"
)

func ConnString(cfg *config.DatabaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Server, cfg.Port, cfg.Database,
	)
}
