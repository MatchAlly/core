package club

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	ErrDuplicateEntry = errors.New("already exists")
	ErrNotFound       = errors.New("not found")
)

type Repository interface {
	GetClub(ctx context.Context, id uint) (*Club, error)
	GetClubs(ctx context.Context, ids []uint) ([]Club, error)
	CreateClub(ctx context.Context, Club *Club) (clubId uint, err error)
	DeleteClub(ctx context.Context, id uint) error
	UpdateClub(ctx context.Context, id uint, name string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetClub(ctx context.Context, id uint) (*Club, error) {
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

func (r *repository) GetClubs(ctx context.Context, ids []uint) ([]Club, error) {
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

func (r *repository) CreateClub(ctx context.Context, c *Club) (uint, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO clubs (name) VALUES ($1) RETURNING id", c.Name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (r *repository) DeleteClub(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM clubs WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateClub(ctx context.Context, id uint, name string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE clubs SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return err
	}

	return nil
}
