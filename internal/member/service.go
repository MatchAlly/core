package member

import (
	"context"
)

type Service interface {
	GetMembersInClub(ctx context.Context, clubId int) ([]Member, error)
	GetUserMemberships(ctx context.Context, userId int) ([]Member, error)
	UpdateRole(ctx context.Context, memberId int, role Role) error
	DeleteMember(ctx context.Context, memberId int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetMembersInClub(ctx context.Context, clubId int) ([]Member, error) {
	return s.repo.GetMembersInClub(ctx, clubId)
}

func (s *service) GetUserMemberships(ctx context.Context, userId int) ([]Member, error) {
	return s.repo.GetUserMemberships(ctx, userId)
}

func (s *service) UpdateRole(ctx context.Context, memberId int, role Role) error {
	return s.repo.UpdateRole(ctx, memberId, role)
}

func (s *service) DeleteMember(ctx context.Context, memberId int) error {
	return s.repo.DeleteMember(ctx, memberId)
}
