package psql

import (
	"ai-bot/internal/config"
	"ai-bot/internal/lib/postgres"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db      *sqlx.DB
	builder squirrel.StatementBuilderType
}

func New(cfg *config.DatabaseConfig) (*Storage, error) {
	const op = "storage.psql.New"

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	db, err := sqlx.Open("postgres", postgres.ConnString(cfg))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open postgres connection: %w", op, err)
	}

	return &Storage{
		db:      db,
		builder: builder,
	}, nil
}
