package invite

import "context"

type Service interface {
	CreateInvite(ctx context.Context, userId, clubId uint) error
	CreateInvites(ctx context.Context, userIds []uint, clubId uint) error
	GetInvitesByUserId(ctx context.Context, userId uint) ([]Invite, error)
	GetInvitesByClubId(ctx context.Context, clubId uint) ([]Invite, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateInvite(ctx context.Context, userId, clubId uint) error {
	return s.repo.CreateInvite(ctx, userId, clubId)
}

func (s *service) CreateInvites(ctx context.Context, userIds []uint, clubId uint) error {
	return s.repo.CreateInvites(ctx, userIds, clubId)
}

func (s *service) GetInvitesByUserId(ctx context.Context, userId uint) ([]Invite, error) {
	return s.repo.GetInvitesByUserId(ctx, userId)
}

func (s *service) GetInvitesByClubId(ctx context.Context, clubId uint) ([]Invite, error) {
	return s.repo.GetInvitesByClubId(ctx, clubId)
}
