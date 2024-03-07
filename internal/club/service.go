package club

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type Service interface {
	GetClub(ctx context.Context, id uint) (*Club, error)
	GetClubs(ctx context.Context, ids []uint) ([]Club, error)
	GetUserIdsInClub(ctx context.Context, id uint) ([]uint, error)
	GetInvitesByUserId(ctx context.Context, userId uint) ([]Invite, error)
	InviteToClub(ctx context.Context, userIds []uint, clubId uint) error
	CreateClub(ctx context.Context, name string, adminUserId uint) (clubId uint, err error)
	RemoveUserFromClub(ctx context.Context, userId uint, clubId uint) error
	DeleteClub(ctx context.Context, id uint) error
	UpdateClub(ctx context.Context, id uint, name string) error
	UpdateUserRole(ctx context.Context, userId uint, clubId uint, role Role) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetClub(ctx context.Context, id uint) (*Club, error) {
	club, err := s.repo.GetClub(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Club")
	}

	return club, nil
}

func (s *service) GetClubs(ctx context.Context, ids []uint) ([]Club, error) {
	clubs, err := s.repo.GetClubs(ctx, ids)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Clubs")
	}

	return clubs, nil
}

func (s *service) GetUserIdsInClub(ctx context.Context, id uint) ([]uint, error) {
	userIds, err := s.repo.GetUserIdsInClub(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get userIds in Club")
	}

	return userIds, nil
}

func (s *service) GetInvitesByUserId(ctx context.Context, userId uint) ([]Invite, error) {
	invites, err := s.repo.GetInvitesByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user invites")
	}

	return invites, nil
}

func (s *service) CreateClub(ctx context.Context, name string, adminUserId uint) (uint, error) {
	club := &Club{
		Name: name,
	}

	fmt.Println("Creating club with name: ", name)

	clubId, err := s.repo.CreateClub(ctx, club)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create Club")
	}

	if err := s.repo.AddUserToClub(ctx, adminUserId, clubId, AdminRole); err != nil {
		return 0, errors.Wrap(err, "failed to add creating user to Club")
	}

	return clubId, nil
}

func (s *service) RemoveUserFromClub(ctx context.Context, userId uint, clubId uint) error {
	if err := s.repo.RemoveUserFromClub(ctx, userId, clubId); err != nil {
		return errors.Wrap(err, "failed to remove user from Club")
	}

	return nil
}

func (s *service) DeleteClub(ctx context.Context, id uint) error {
	if err := s.repo.DeleteClub(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete Club")
	}

	return nil
}

func (s *service) UpdateClub(ctx context.Context, id uint, name string) error {
	if err := s.repo.UpdateClub(ctx, id, name); err != nil {
		return errors.Wrap(err, "failed to update Club")
	}

	return nil
}

func (s *service) UpdateUserRole(ctx context.Context, userId uint, clubId uint, role Role) error {
	if err := s.repo.UpdateUserRole(ctx, userId, clubId, role); err != nil {
		return errors.Wrap(err, "failed to update user role")
	}

	return nil
}

func (s *service) InviteToClub(ctx context.Context, userIds []uint, clubId uint) error {
	if err := s.repo.InviteToClub(ctx, userIds, clubId); err != nil {
		return errors.Wrap(err, "failed to invite users to club")
	}

	return nil
}
