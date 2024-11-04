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
	match := &Match{
		ClubID:  clubID,
		GameID:  gameID,
		TeamIDs: teamIDs,
		Sets:    sets,
	}

	id, err := s.repo.CreateMatch(ctx, match)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create match")
	}

	return id, nil
}
