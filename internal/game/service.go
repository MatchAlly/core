package game

import "context"

type Service interface {
	GetGame(ctx context.Context, id int) (*Game, error)
	GetGames(ctx context.Context, ids []int) ([]Game, error)
	CreateGame(ctx context.Context, game *Game) (int, error)
	UpdateGame(ctx context.Context, game *Game) error
	DeleteGame(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) GetGame(ctx context.Context, id int) (*Game, error) {
	return s.repo.GetGame(ctx, id)
}

func (s *service) GetGames(ctx context.Context, ids []int) ([]Game, error) {
	return s.repo.GetGames(ctx, ids)
}

func (s *service) CreateGame(ctx context.Context, game *Game) (int, error) {
	return s.repo.CreateGame(ctx, game)
}

func (s *service) UpdateGame(ctx context.Context, game *Game) error {
	return s.repo.UpdateGame(ctx, game)
}

func (s *service) DeleteGame(ctx context.Context, id int) error {
	return s.repo.DeleteGame(ctx, id)
}
