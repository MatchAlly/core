package subscription

import "context"

type Service interface {
	CreateSubscription(ctx context.Context, subscription Subscription) error
	GetSubscriptionByUserID(ctx context.Context, userID int) (*Subscription, error)
	UpdateSubscription(ctx context.Context, subscription Subscription) error
	DeleteSubscription(ctx context.Context, userID int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateSubscription(ctx context.Context, subscription Subscription) error {
	return s.repo.CreateSubscription(ctx, subscription)
}

func (s *service) GetSubscriptionByUserID(ctx context.Context, userID int) (*Subscription, error) {
	return s.repo.GetSubscriptionByUserID(ctx, userID)
}

func (s *service) UpdateSubscription(ctx context.Context, subscription Subscription) error {
	return s.repo.UpdateSubscription(ctx, subscription)
}

func (s *service) DeleteSubscription(ctx context.Context, userID int) error {
	return s.repo.DeleteSubscription(ctx, userID)
}
