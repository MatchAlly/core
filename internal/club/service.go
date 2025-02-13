package club

import (
	"context"
	"core/internal/game"
	"fmt"
)

type Service interface {
	GetClub(ctx context.Context, id int) (*Club, error)
	GetClubs(ctx context.Context, ids []int) ([]Club, error)
	CreateClub(ctx context.Context, name string, adminUserId int) (clubId int, err error)
	DeleteClub(ctx context.Context, id int) error
	UpdateClub(ctx context.Context, id int, name string) error
	GetGames(ctx context.Context, clubID int) ([]game.Game, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetClub(ctx context.Context, id int) (*Club, error) {
	club, err := s.repo.GetClub(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get club: %w", err)
	}

	return club, nil
}

func (s *service) GetClubs(ctx context.Context, ids []int) ([]Club, error) {
	clubs, err := s.repo.GetClubs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get clubs: %w", err)
	}

	return clubs, nil
}

func (s *service) CreateClub(ctx context.Context, name string, adminUserId int) (int, error) {
	c := &Club{
		Name: name,
	}

	clubId, err := s.repo.CreateClub(ctx, c)
	if err != nil {
		return 0, fmt.Errorf("failed to create club: %w", err)
	}

	return clubId, nil
}

func (s *service) DeleteClub(ctx context.Context, id int) error {
	if err := s.repo.DeleteClub(ctx, id); err != nil {
		return fmt.Errorf("failed to delete club: %w", err)
	}

	return nil
}

func (s *service) UpdateClub(ctx context.Context, id int, name string) error {
	if err := s.repo.UpdateClub(ctx, id, name); err != nil {
		return fmt.Errorf("failed to update club: %w", err)
	}

	return nil
}

func (s *service) GetGames(ctx context.Context, clubID int) ([]game.Game, error) {
	games, err := s.repo.GetGames(ctx, clubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get games: %w", err)
	}

	return games, nil
}
