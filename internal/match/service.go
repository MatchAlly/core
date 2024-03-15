package match

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type Service interface {
	CreateMatch(ctx context.Context, teamA, teamB []uint, scoresA, scoresB []int, result Result) error
	DetermineResult(ctx context.Context, teamA, teamB []uint, scoresA, scoresB []int) (result Result, winners []uint, losers []uint)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateMatch(ctx context.Context, teamA, teamB []uint, scoresA, scoresB []int, result Result) error {
	sets := make([]string, len(scoresA))
	for i, scoreA := range scoresA {
		sets[i] = fmt.Sprintf("%d-%d", scoreA, scoresB[i])
	}

	match := &Match{
		TeamA:  teamA,
		TeamB:  teamB,
		Sets:   sets,
		Result: result,
	}

	if err := s.repo.CreateMatch(ctx, match); err != nil {
		return errors.Wrap(err, "failed to create match")
	}

	return nil
}

func (s *service) DetermineResult(ctx context.Context, teamA, teamB []uint, scoresA, scoresB []int) (Result, []uint, []uint) {
	teamASetWins := 0
	teamBSetWins := 0

	for i, scoreA := range scoresA {
		if scoreA > scoresB[i] {
			teamASetWins++
		} else if scoreA < scoresB[i] {
			teamBSetWins++
		}
	}

	if teamASetWins > teamBSetWins {
		return TeamAWins, teamA, teamB
	} else if teamBSetWins > teamASetWins {
		return TeamBWins, teamB, teamA
	} else {
		return Draw, nil, nil
	}

}
