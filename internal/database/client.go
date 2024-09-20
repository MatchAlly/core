package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func NewClient(ctx context.Context, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
