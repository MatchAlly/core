package statistic

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetStatisticsByUserId(ctx context.Context, userId uint) ([]Statistic, error)
	CreateStatistic(ctx context.Context, userId, gameId uint) error
	UpdateStatisticsByUserIds(ctx context.Context, userIds []uint, gameId uint, result MatchResult) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetStatisticsByUserId(ctx context.Context, userId uint) ([]Statistic, error) {
	stats, err := s.repo.GetStatisticsByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statistics for user %d", userId)
	}

	return stats, nil
}

func (s *service) CreateStatistic(ctx context.Context, userId, gameId uint) error {
	if err := s.repo.CreateStatistic(ctx, userId, gameId); err != nil {
		return errors.Wrapf(err, "failed to create statistics for user %d", userId)
	}

	return nil
}

func (s *service) UpdateStatisticsByUserIds(ctx context.Context, userIds []uint, gameId uint, result MatchResult) error {
	oldStatistics, err := s.repo.GetStatisticsByUserIds(ctx, userIds)
	if err != nil {
		return errors.Wrapf(err, "failed to get statistics for users %v", userIds)
	}

	updatedStatistics := make([]Statistic, len(userIds))
	for i := range userIds {
		stats := oldStatistics[i]

		switch result {
		case ResultWin:
			stats.Wins++
			if stats.Streak >= 0 {
				stats.Streak++
			} else {
				stats.Streak = 1
			}
		case ResultLoss:
			stats.Losses++
			if stats.Streak <= 0 {
				stats.Streak--
			} else {
				stats.Streak = -1
			}
		case ResultDraw:
			stats.Draws++
			stats.Streak = 0
		}

		updatedStatistics[i] = *stats
	}

	if err := s.repo.UpdateStatistics(ctx, updatedStatistics); err != nil {
		return errors.Wrap(err, "failed to update statistics")
	}

	return nil
}
