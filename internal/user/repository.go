package user

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const UniqueViolationCode = pq.ErrorCode("23505")

var (
	ErrDuplicateEntry = errors.New("already exists")
	ErrNotFound       = errors.New("not found")
)

type Repository interface {
	GetUser(ctx context.Context, id int) (*User, error)
	GetUsers(ctx context.Context, ids []int) ([]User, error)
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
	var user *User

	err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *repository) GetUsers(ctx context.Context, ids []int) ([]User, error) {
	var users []User

	query, args, err := sqlx.In("SELECT * FROM users WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	if err = r.db.SelectContext(ctx, &users, query, args...); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user *User

	err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *repository) CreateUser(ctx context.Context, user *User) (int, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO users (email, name, hash) VALUES ($1, $2, $3)", user.Email, user.Name, user.Hash)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == UniqueViolationCode {
			return 0, ErrDuplicateEntry
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
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
