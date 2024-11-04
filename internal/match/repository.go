package match

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	ErrDuplicateEntry = errors.New("already exists")
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
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO matches (club_id, game_id, team_ids, sets) VALUES ($1, $2, $3, $4) RETURNING id",
		match.ClubID, match.GameID, match.TeamIDs, match.Sets,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
