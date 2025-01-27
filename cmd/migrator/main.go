package main

import (
	"ai-bot/internal/config"
	psql "ai-bot/internal/lib/postgres"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"log"
	"os"
)

func main() {
	cfg := config.MustLoad()
	migPath := os.Getenv("MIGRATIONS_PATH")
	if migPath == "" {
		log.Fatal("migrations path not set")
	}

	m, err := migrate.New(
		"file://"+migPath,
		psql.ConnString(cfg.Database),
	)
	if err != nil {
		log.Fatal("failed to create migrator: " + err.Error())
	}
	m.Up()
}
