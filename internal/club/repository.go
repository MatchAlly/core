package club

import (
	"context"
	"core/internal/game"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrDuplicateEntry = fmt.Errorf("duplicate entry")
	ErrNotFound       = fmt.Errorf("not found")
)

type Repository interface {
	GetClub(ctx context.Context, id int) (*Club, error)
	GetClubs(ctx context.Context, ids []int) ([]Club, error)
	CreateClub(ctx context.Context, Club *Club) (clubId int, err error)
	DeleteClub(ctx context.Context, id int) error
	UpdateClub(ctx context.Context, id int, name string) error
	GetGames(ctx context.Context, clubID int) ([]game.Game, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetClub(ctx context.Context, id int) (*Club, error) {
	var c *Club

	err := r.db.GetContext(ctx, c, "SELECT * FROM clubs WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return c, nil
}

func (r *repository) GetClubs(ctx context.Context, ids []int) ([]Club, error) {
	var clubs []Club

	query, args, err := sqlx.In("SELECT * FROM clubs WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &clubs, query, args...)
	if err != nil {
		return nil, err
	}

	return clubs, nil
}

func (r *repository) CreateClub(ctx context.Context, c *Club) (int, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO clubs (name) VALUES ($1) RETURNING id", c.Name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) DeleteClub(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM clubs WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateClub(ctx context.Context, id int, name string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE clubs SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetGames(ctx context.Context, clubID int) ([]game.Game, error) {
	var games []game.Game

	err := r.db.SelectContext(ctx, &games, "SELECT * FROM games WHERE club_id = $1", clubID)
	if err != nil {
		return nil, err
	}

	return games, nil
}
