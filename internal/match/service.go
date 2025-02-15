package match

import (
	"context"
	"fmt"
)

type Service interface {
	CreateMatch(ctx context.Context, clubID, gameID int, teams []Team, sets []string) (int, error)
	GetMatches(ctx context.Context, clubID int, gameID *int) ([]Match, error)
	GetOrCreateTeams(ctx context.Context, clubID int, members [][]int) ([]Team, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateMatch(ctx context.Context, clubID, gameID int, teams []Team, sets []string) (int, error) {
	m := &Match{
		ClubID: clubID,
		GameID: gameID,
		Teams:  teams,
		Sets:   sets,
	}

	matchID, err := s.repo.CreateMatch(ctx, m)
	if err != nil {
		return 0, fmt.Errorf("failed to create match: %w", err)
	}

	return matchID, nil
}

func (s *service) GetMatches(ctx context.Context, clubID int, gameID *int) ([]Match, error) {
	var matches []Match
	var err error

	if gameID == nil {
		matches, err = s.repo.GetMatches(ctx, clubID)
		if err != nil {
			return nil, fmt.Errorf("failed to get matches: %w", err)
		}
	} else {
		matches, err = s.repo.GetMatchesByGame(ctx, clubID, *gameID)
		if err != nil {
			return nil, fmt.Errorf("failed to get matches by game: %w", err)
		}
	}

	return matches, nil
}

func (s *service) GetOrCreateTeams(ctx context.Context, clubID int, memberIDTeams [][]int) ([]Team, error) {
	teams := make([]Team, len(memberIDTeams))

	for _, memberIDTeam := range memberIDTeams {
		exists, teamID, err := s.repo.TeamOfMembersExists(ctx, clubID, memberIDTeam)
		if err != nil {
			return nil, fmt.Errorf("failed to check if team exists: %w", err)
		}

		if exists {
			team, err := s.repo.GetTeam(ctx, teamID)
			if err != nil {
				return nil, fmt.Errorf("failed to get team: %w", err)
			}

			teams = append(teams, *team)
		} else {
			teamID, err := s.repo.CreateTeam(ctx, clubID, memberIDTeam)
			if err != nil {
				return nil, fmt.Errorf("failed to create team: %w", err)
			}

			team, err := s.repo.GetTeam(ctx, teamID)
			if err != nil {
				return nil, fmt.Errorf("failed to get team: %w", err)
			}

			teams = append(teams, *team)
		}
	}

	return teams, nil
}
