package user

import (
	"context"
	"database/sql"
	"fmt"

	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const UniqueViolationCode = "23505"

var (
	ErrDuplicateEntry = fmt.Errorf("already exists")
	ErrNotFound       = fmt.Errorf("not found")
)

type Repository interface {
	GetUser(ctx context.Context, id int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) (int, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, userID int, hash string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetUser(ctx context.Context, id int) (*User, error) {
	var user User

	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) CreateUser(ctx context.Context, user *User) (int, error) {
	var id int

	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (email, name, hash) VALUES ($1, $2, $3) RETURNING id",
		user.Email, user.Name, user.Hash).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok && pqErr.Code == UniqueViolationCode {
			return 0, ErrDuplicateEntry
		}
		return 0, err
	}

	return id, nil
}

func (r *repository) DeleteUser(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET email = $1, name = $2 WHERE id = $3", user.Email, user.Name, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdatePassword(ctx context.Context, userID int, hash string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET hash = $1 WHERE id = $2", hash, userID)
	if err != nil {
		return err
	}

	return nil
}
