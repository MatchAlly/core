package member

import (
	"context"
)

type Service interface {
	GetMembersInClub(ctx context.Context, clubId uint) ([]Member, error)
	GetUserMemberships(ctx context.Context, userId uint) ([]Member, error)
	UpdateRole(ctx context.Context, memberId uint, role Role) error
	DeleteMembership(ctx context.Context, memberId uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetMembersInClub(ctx context.Context, clubId uint) ([]Member, error) {
	return s.repo.GetMembersInClub(ctx, clubId)
}

func (s *service) GetUserMemberships(ctx context.Context, userId uint) ([]Member, error) {
	return s.repo.GetUserMemberships(ctx, userId)
}

func (s *service) UpdateRole(ctx context.Context, memberId uint, role Role) error {
	return s.repo.UpdateRole(ctx, memberId, role)
}

func (s *service) DeleteMembership(ctx context.Context, memberId uint) error {
	return s.repo.DeleteMembership(ctx, memberId)
}
