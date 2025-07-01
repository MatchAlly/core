package subscription

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Subscription, error)
	Update(ctx context.Context, userID uuid.UUID, tier Tier) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, id uuid.UUID) error {
	return s.repo.Create(ctx, id)
}

func (s *service) GetByUserID(ctx context.Context, userID uuid.UUID) (*Subscription, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *service) Update(ctx context.Context, userID uuid.UUID, tier Tier) error {
	return s.repo.Update(ctx, userID, tier)
}
