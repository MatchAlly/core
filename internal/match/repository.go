package match

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository interface {
	CreateMatch(ctx context.Context, match *Match) (int, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateMatch(ctx context.Context, match *Match) (int, error) {
	// TODO: Implement
	return 0, errors.New("not implemented")
}
