package club

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetClub(ctx context.Context, id uint) (*Club, error)
	GetClubs(ctx context.Context, ids []uint) ([]Club, error)
	CreateClub(ctx context.Context, name string, adminUserId uint) (clubId uint, err error)
	DeleteClub(ctx context.Context, id uint) error
	UpdateClub(ctx context.Context, id uint, name string) error
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

func (s *service) CreateClub(ctx context.Context, name string, adminUserId uint) (uint, error) {
	c := &Club{
		Name: name,
	}

	clubId, err := s.repo.CreateClub(ctx, c)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create club")
	}

	return clubId, nil
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
