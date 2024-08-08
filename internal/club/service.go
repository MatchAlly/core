package club

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetClub(ctx context.Context, id uint) (*Club, error)
	GetClubs(ctx context.Context, ids []uint) ([]Club, error)
	GetClubIDsWithUserID(ctx context.Context, userId uint) ([]uint, error)
	GetMembers(ctx context.Context, id uint) ([]Member, error)
	CreateClub(ctx context.Context, name string, adminUserId uint) (clubId uint, err error)
	DeleteMember(ctx context.Context, memberId uint) error
	DeleteClub(ctx context.Context, id uint) error
	UpdateClub(ctx context.Context, id uint, name string) error
	UpdateMemberRole(ctx context.Context, memberId uint, role Role) error
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
		return nil, errors.Wrap(err, "failed to get club")
	}

	return club, nil
}

func (s *service) GetClubs(ctx context.Context, ids []uint) ([]Club, error) {
	clubs, err := s.repo.GetClubs(ctx, ids)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get clubs")
	}

	return clubs, nil
}

func (s *service) GetClubIDsWithUserID(ctx context.Context, userId uint) ([]uint, error) {
	clubIds, err := s.repo.GetClubIDsWithUserID(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get clubIds with userId")
	}

	return clubIds, nil
}

func (s *service) GetMembers(ctx context.Context, id uint) ([]Member, error) {
	members, err := s.repo.GetMembers(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get userIds in club")
	}

	return members, nil
}

func (s *service) CreateClub(ctx context.Context, name string, adminUserId uint) (uint, error) {
	c := &Club{
		Name: name,
	}

	clubId, err := s.repo.CreateClub(ctx, c)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create club")
	}

	if err := s.repo.AddUserToClub(ctx, adminUserId, clubId, AdminRole); err != nil {
		return 0, errors.Wrap(err, "failed to add creating user to club")
	}

	return clubId, nil
}

func (s *service) DeleteMember(ctx context.Context, memberId uint) error {
	if err := s.repo.DeleteMember(ctx, memberId); err != nil {
		return errors.Wrap(err, "failed to delete member")
	}

	return nil
}

func (s *service) DeleteClub(ctx context.Context, id uint) error {
	if err := s.repo.DeleteClub(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete club")
	}

	return nil
}

func (s *service) UpdateClub(ctx context.Context, id uint, name string) error {
	if err := s.repo.UpdateClub(ctx, id, name); err != nil {
		return errors.Wrap(err, "failed to update club")
	}

	return nil
}

func (s *service) UpdateMemberRole(ctx context.Context, memberId uint, role Role) error {
	if err := s.repo.UpdateMemberRole(ctx, memberId, role); err != nil {
		return errors.Wrap(err, "failed to update member role")
	}

	return nil
}
