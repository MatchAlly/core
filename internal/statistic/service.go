package statistic

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	UpdateStatistics(ctx context.Context, memberID, gameID uuid.UUID, won, drawn bool) error
	GetStatistics(ctx context.Context, memberID, gameID uuid.UUID) (*Statistic, error)
	GetStatisticsByGame(ctx context.Context, gameID uuid.UUID) ([]Statistic, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) UpdateStatistics(ctx context.Context, memberID, gameID uuid.UUID, won, drawn bool) error {
	stats, err := s.repo.GetStatistics(ctx, memberID, gameID)
	if err != nil {
		// If statistics don't exist, create new ones
		stats = &Statistic{
			MemberId: memberID,
			GameId:   gameID,
			Wins:     0,
			Losses:   0,
			Draws:    0,
			Streak:   0,
		}
	}

	// Update statistics based on match result
	if won {
		stats.Wins++
		if stats.Streak > 0 {
			stats.Streak++
		} else {
			stats.Streak = 1
		}
	} else if drawn {
		stats.Draws++
		stats.Streak = 0
	} else {
		stats.Losses++
		if stats.Streak < 0 {
			stats.Streak--
		} else {
			stats.Streak = -1
		}
	}

	if stats.ID == uuid.Nil {
		// Create new statistics
		_, err = s.repo.CreateStatistics(ctx, stats)
	} else {
		// Update existing statistics
		err = s.repo.UpdateStatistics(ctx, stats)
	}

	if err != nil {
		return fmt.Errorf("failed to update statistics: %w", err)
	}

	return nil
}

func (s *service) GetStatistics(ctx context.Context, memberID, gameID uuid.UUID) (*Statistic, error) {
	stats, err := s.repo.GetStatistics(ctx, memberID, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}
	return stats, nil
}

func (s *service) GetStatisticsByGame(ctx context.Context, gameID uuid.UUID) ([]Statistic, error) {
	stats, err := s.repo.GetStatisticsByGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics by game: %w", err)
	}
	return stats, nil
}
