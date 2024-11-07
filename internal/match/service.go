package match

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	CreateMatch(ctx context.Context, clubID, gameID int, teamIDs []int, sets []string) (int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateMatch(ctx context.Context, clubID, gameID int, teamIDs []int, sets []string) (int, error) {
	// TODO: Implement¨
	return 0, errors.New("not implemented")
}
