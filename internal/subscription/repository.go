package subscription

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, userID int) error
	GetByUserID(ctx context.Context, userID int) (*Subscription, error)
	Update(ctx context.Context, userID int, tier Tier) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO subscriptions (user_id) VALUES ($1)", userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByUserID(ctx context.Context, userID int) (*Subscription, error) {
	var s Subscription
	err := r.db.GetContext(ctx, &s, "SELECT * FROM subscriptions WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *repository) Update(ctx context.Context, userID int, tier Tier) error {
	now := time.Now().Format(time.RFC3339)
	_, err := r.db.ExecContext(ctx, "UPDATE subscriptions SET tier = $1, updated_at = $2 WHERE user_id = $3", tier, now, userID)
	if err != nil {
		return err
	}

	return nil
}
