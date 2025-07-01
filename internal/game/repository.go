package game

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetGame(ctx context.Context, id uuid.UUID) (*Game, error)
	GetGames(ctx context.Context, ids []uuid.UUID) ([]Game, error)
	CreateGame(ctx context.Context, game *Game) (uuid.UUID, error)
	UpdateGame(ctx context.Context, game *Game) error
	DeleteGame(ctx context.Context, id uuid.UUID) error
	GetGameModes(ctx context.Context, gameID uuid.UUID) ([]Gamemode, error)
	AddGameMode(ctx context.Context, gameID uuid.UUID, mode Mode) error
	RemoveGameMode(ctx context.Context, gameID uuid.UUID, mode Mode) error
	IsGameNameUnique(ctx context.Context, clubID uuid.UUID, name string, excludeGameID uuid.UUID) (bool, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) GetGame(ctx context.Context, id uuid.UUID) (*Game, error) {
	var game *Game

	err := r.db.GetContext(ctx, game, "SELECT * FROM games WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *repository) GetGames(ctx context.Context, ids []uuid.UUID) ([]Game, error) {
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

func (r *repository) CreateGame(ctx context.Context, game *Game) (uuid.UUID, error) {
	var id uuid.UUID

	err := r.db.QueryRowContext(ctx,
		"INSERT INTO games (club_id, name) VALUES ($1, $2) RETURNING id",
		game.ClubID, game.Name).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
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

func (r *repository) DeleteGame(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM games WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetGameModes(ctx context.Context, gameID uuid.UUID) ([]Gamemode, error) {
	var modes []Gamemode
	err := r.db.SelectContext(ctx, &modes, "SELECT * FROM game_modes WHERE game_id = $1", gameID)
	if err != nil {
		return nil, err
	}
	return modes, nil
}

func (r *repository) AddGameMode(ctx context.Context, gameID uuid.UUID, mode Mode) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO game_modes (game_id, mode) VALUES ($1, $2) ON CONFLICT (game_id, mode) DO NOTHING",
		gameID, mode)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) RemoveGameMode(ctx context.Context, gameID uuid.UUID, mode Mode) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM game_modes WHERE game_id = $1 AND mode = $2",
		gameID, mode)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) IsGameNameUnique(ctx context.Context, clubID uuid.UUID, name string, excludeGameID uuid.UUID) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM games WHERE club_id = $1 AND name = $2 AND id != $3",
		clubID, name, excludeGameID)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
