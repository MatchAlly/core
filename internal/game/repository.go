package game

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetGame(ctx context.Context, id int) (*Game, error)
	GetGames(ctx context.Context, ids []int) ([]Game, error)
	CreateGame(ctx context.Context, game *Game) (int, error)
	UpdateGame(ctx context.Context, game *Game) error
	DeleteGame(ctx context.Context, id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) GetGame(ctx context.Context, id int) (*Game, error) {
	var game *Game

	err := r.db.GetContext(ctx, game, "SELECT * FROM games WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *repository) GetGames(ctx context.Context, ids []int) ([]Game, error) {
	var games []Game

	query, args, err := sqlx.In("SELECT * FROM games WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &games, query, args...)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (r *repository) CreateGame(ctx context.Context, game *Game) (int, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO games (club_id, name) VALUES ($1, $2)",
		game.ClubID, game.Name,
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

func (r *repository) UpdateGame(ctx context.Context, game *Game) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE games SET club_id = $1, name = $2 WHERE id = $3",
		game.ClubID, game.Name, game.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteGame(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM games WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
