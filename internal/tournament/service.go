//go:generate mockgen --source=service.go -destination=service_mock.go -package=tournament
package tournament

import (
	"fmt"
	"math/rand"
)

type Service interface {
	CreateTournament(teams [][]uint, format TournamentFormat, isSeeded bool) (*Tournament, error)
	CreateNextRound(tourn *Tournament) (*Tournament, error)
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) CreateTournament(teams [][]uint, format TournamentFormat, isSeeded bool) (*Tournament, error) {
	// If seeded, the order of the teams is asssumed to be the seeding order, with the first team being the top seed.
	if !isSeeded {
		rand.Shuffle(len(teams), func(i, j int) { teams[i], teams[j] = teams[j], teams[i] })
	}

	var (
		tourn *Tournament
		err   error
	)

	switch format {
	case FormatSingleElimination:
		tourn, err = s.createSingleEliminationTournament(teams)
	case FormatDoubleElimination:
		tourn, err = s.createDoubleEliminationTournament(teams)
	case FormatRoundRobin:
		tourn, err = s.createRoundRobinTournament(teams)
	case FormatSwiss:
		tourn, err = s.createSwissTournament(teams)
	default:
		return nil, fmt.Errorf("unknown tournament format: %s", format)
	}

	if err != nil {
		return nil, err
	}

	return tourn, nil
}

func (s *ServiceImpl) CreateNextRound(tourn *Tournament) (*Tournament, error) {
	panic("not implemented")
}
