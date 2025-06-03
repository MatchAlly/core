package subscription

import "context"

type Service interface {
	Create(ctx context.Context, userID int) error
	GetByUserID(ctx context.Context, userID int) (*Subscription, error)
	Update(ctx context.Context, userID int, tier Tier) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, id int) error {
	return s.repo.Create(ctx, id)
}

func (s *service) GetByUserID(ctx context.Context, userID int) (*Subscription, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *service) Update(ctx context.Context, userID int, tier Tier) error {
	return s.repo.Update(ctx, userID, tier)
}
