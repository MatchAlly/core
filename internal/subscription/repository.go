package subscription

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateSubscription(ctx context.Context, subscription Subscription) error
	GetSubscriptionByUserID(ctx context.Context, userID int) (*Subscription, error)
	UpdateSubscription(ctx context.Context, subscription Subscription) error
	DeleteSubscription(ctx context.Context, userID int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateSubscription(ctx context.Context, subscription Subscription) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO subscriptions (user_id, managed_organization_ids, total_managed_users, tier) VALUES (:user_id, :managed_organization_ids, :total_managed_users, :tier)", subscription)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetSubscriptionByUserID(ctx context.Context, userID int) (*Subscription, error) {
	var s Subscription
	err := r.db.GetContext(ctx, &s, "SELECT * FROM subscriptions WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *repository) UpdateSubscription(ctx context.Context, subscription Subscription) error {
	_, err := r.db.NamedExecContext(ctx, "UPDATE subscriptions SET managed_organization_ids = :managed_organization_ids, total_managed_users = :total_managed_users, tier = :tier WHERE user_id = :user_id", subscription)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteSubscription(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM subscriptions WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	return nil
}
