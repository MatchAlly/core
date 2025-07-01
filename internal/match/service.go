package match

import (
	"context"
	"core/internal/game"
	"core/internal/rating"
	"core/internal/statistic"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	CreateMatch(ctx context.Context, clubID, gameID uuid.UUID, teams []Team, sets []string, mode game.Mode) (uuid.UUID, error)
	GetMatches(ctx context.Context, clubID uuid.UUID, gameID *uuid.UUID) ([]Match, error)
	GetOrCreateTeams(ctx context.Context, clubID uuid.UUID, members [][]uuid.UUID) ([]Team, error)
}

type service struct {
	repo      Repository
	game      game.Service
	rating    rating.Service
	statistic statistic.Service
}

func NewService(repo Repository, game game.Service, rating rating.Service, statistic statistic.Service) Service {
	return &service{
		repo:      repo,
		game:      game,
		rating:    rating,
		statistic: statistic,
	}
}

func (s *service) validateTeamMembers(ctx context.Context, clubID uuid.UUID, memberIDs []uuid.UUID) error {
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

func (s *service) validateGameMode(ctx context.Context, gameID uuid.UUID, mode game.Mode, numTeams int) error {
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

func (s *service) CreateMatch(ctx context.Context, clubID, gameID uuid.UUID, teams []Team, sets []string, mode game.Mode) (uuid.UUID, error) {
	// Validate game exists in club
	g, err := s.game.GetGame(ctx, gameID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get game: %w", err)
	}
	if g.ClubID != clubID {
		return uuid.Nil, fmt.Errorf("game does not belong to the specified club")
	}

	// Validate team members
	for _, team := range teams {
		memberIDs := make([]uuid.UUID, len(team.Members))
		for i, member := range team.Members {
			memberIDs[i] = member.UserID
		}
		if err := s.validateTeamMembers(ctx, clubID, memberIDs); err != nil {
			return uuid.Nil, err
		}
	}

	// Validate game mode and number of teams
	if err := s.validateGameMode(ctx, gameID, mode, len(teams)); err != nil {
		return uuid.Nil, err
	}

	m := &Match{
		ClubID:   clubID,
		GameID:   gameID,
		Teams:    teams,
		Sets:     sets,
		Gamemode: mode,
		Ranked:   true, // Set ranked to true by default
	}

	matchID, err := s.repo.CreateMatch(ctx, m)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create match: %w", err)
	}

	// Update statistics for each player
	for i, team := range teams {
		for _, member := range team.Members {
			// For now, we'll consider the first team as the winner
			won := i == 0
			drawn := false
			if err := s.statistic.UpdateStatistics(ctx, member.ID, gameID, won, drawn); err != nil {
				// Log the error but don't fail the match creation
				fmt.Printf("failed to update statistics for member %d: %v\n", member.ID, err)
			}
		}
	}

	// Update ratings if the match is ranked
	if m.Ranked {
		// Convert teams to member IDs for rating update
		teamsByMemberIDs := make([][]uuid.UUID, len(teams))
		for i, team := range teams {
			memberIDs := make([]uuid.UUID, len(team.Members))
			for j, member := range team.Members {
				memberIDs[j] = member.ID
			}
			teamsByMemberIDs[i] = memberIDs
		}

		// For now, we'll consider the first team as the winner
		// In a real implementation, you would determine the ranks based on the match results
		ranks := make([]int, len(teams))
		for i := range teams {
			if i == 0 {
				ranks[i] = 1 // First place
			} else {
				ranks[i] = 2 // Second place
			}
		}

		if err := s.rating.UpdateRatingsByRanks(ctx, teamsByMemberIDs, ranks); err != nil {
			// Log the error but don't fail the match creation
			fmt.Printf("failed to update ratings: %v\n", err)
		}
	}

	return matchID, nil
}

func (s *service) GetMatches(ctx context.Context, clubID uuid.UUID, gameID *uuid.UUID) ([]Match, error) {
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

func (s *service) GetOrCreateTeams(ctx context.Context, clubID uuid.UUID, memberIDTeams [][]uuid.UUID) ([]Team, error) {
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
