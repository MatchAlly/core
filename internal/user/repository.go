package user

import (
	"context"
	"database/sql"
	"fmt"

	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const UniqueViolationCode = "23505"

var (
	ErrDuplicateEntry = fmt.Errorf("already exists")
	ErrNotFound       = fmt.Errorf("not found")
)

type Repository interface {
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) (uuid.UUID, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, hash string) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
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

func (r *repository) CreateUser(ctx context.Context, user *User) (uuid.UUID, error) {
	var id uuid.UUID

	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (email, name, hash) VALUES ($1, $2, $3) RETURNING id",
		user.Email, user.Name, user.Hash).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
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

func (r *repository) UpdatePassword(ctx context.Context, userID uuid.UUID, hash string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET hash = $1 WHERE id = $2", hash, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1", userID)
	if err != nil {
		return err
	}

	return nil
}
