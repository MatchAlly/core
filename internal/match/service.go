package match

import (
	"context"
	"core/internal/game"
	"fmt"
)

type Service interface {
	CreateMatch(ctx context.Context, clubID, gameID int, teams []Team, sets []string, mode game.Mode) (int, error)
	GetMatches(ctx context.Context, clubID int, gameID *int) ([]Match, error)
	GetOrCreateTeams(ctx context.Context, clubID int, members [][]int) ([]Team, error)
}

type service struct {
	repo Repository
	game game.Service
}

func NewService(repo Repository, game game.Service) Service {
	return &service{
		repo: repo,
		game: game,
	}
}

func (s *service) validateTeamMembers(ctx context.Context, clubID int, memberIDs []int) error {
	for _, memberID := range memberIDs {
		isMember, err := s.repo.IsClubMember(ctx, clubID, memberID)
		if err != nil {
			return fmt.Errorf("failed to check club membership: %w", err)
		}
		if !isMember {
			return fmt.Errorf("user %d is not a member of club %d", memberID, clubID)
		}
	}
	return nil
}

func (s *service) validateGameMode(ctx context.Context, gameID int, mode game.Mode, numTeams int) error {
	modes, err := s.game.GetGameModes(ctx, gameID)
	if err != nil {
		return fmt.Errorf("failed to get game modes: %w", err)
	}

	// Check if the mode is supported by the game
	modeSupported := false
	for _, m := range modes {
		if m.Mode == mode {
			modeSupported = true
			break
		}
	}
	if !modeSupported {
		return fmt.Errorf("game mode %s is not supported for this game", mode)
	}

	// Validate number of teams based on mode
	switch mode {
	case game.ModeFreeForAll:
		if numTeams < 2 {
			return fmt.Errorf("free-for-all mode requires at least 2 teams")
		}
	case game.ModeTeam:
		if numTeams != 2 {
			return fmt.Errorf("team mode requires exactly 2 teams")
		}
	case game.ModeCoop:
		if numTeams != 1 {
			return fmt.Errorf("coop mode requires exactly 1 team")
		}
	}

	return nil
}

func (s *service) CreateMatch(ctx context.Context, clubID, gameID int, teams []Team, sets []string, mode game.Mode) (int, error) {
	// Validate game exists in club
	g, err := s.game.GetGame(ctx, gameID)
	if err != nil {
		return 0, fmt.Errorf("failed to get game: %w", err)
	}
	if g.ClubID != clubID {
		return 0, fmt.Errorf("game does not belong to the specified club")
	}

	// Validate team members
	for _, team := range teams {
		memberIDs := make([]int, len(team.Members))
		for i, member := range team.Members {
			memberIDs[i] = member.UserID
		}
		if err := s.validateTeamMembers(ctx, clubID, memberIDs); err != nil {
			return 0, err
		}
	}

	// Validate game mode and number of teams
	if err := s.validateGameMode(ctx, gameID, mode, len(teams)); err != nil {
		return 0, err
	}

	m := &Match{
		ClubID:   clubID,
		GameID:   gameID,
		Teams:    teams,
		Sets:     sets,
		Gamemode: mode,
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
	teams := make([]Team, 0, len(memberIDTeams))

	for _, memberIDTeam := range memberIDTeams {
		exists, teamID, err := s.repo.TeamOfMembersExists(ctx, clubID, memberIDTeam)
		if err != nil {
			return nil, fmt.Errorf("failed to check if team exists: %w", err)
		}

		var team *Team
		if exists {
			team, err = s.repo.GetTeam(ctx, teamID)
			if err != nil {
				return nil, fmt.Errorf("failed to get team: %w", err)
			}
		} else {
			teamID, err := s.repo.CreateTeam(ctx, clubID, memberIDTeam)
			if err != nil {
				return nil, fmt.Errorf("failed to create team: %w", err)
			}

			team, err = s.repo.GetTeam(ctx, teamID)
			if err != nil {
				return nil, fmt.Errorf("failed to get team: %w", err)
			}
		}

		teams = append(teams, *team)
	}

	return teams, nil
}
