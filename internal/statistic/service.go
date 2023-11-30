package statistic

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetStatisticByUserId(ctx context.Context, userId uint) (*Statistic, error)
	GetTopXAmongUserIdsByMeasure(ctx context.Context, topX int, userIds []uint, measure Measure) (topXUserIds []uint, values []int, err error)
	CreateStatistic(ctx context.Context, userId uint) error
	UpdateStatisticsByUserIds(ctx context.Context, userIds []uint, result MatchResult) error
	TransferStatistics(ctx context.Context, fromUserId, toUserId uint) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) GetStatisticByUserId(ctx context.Context, userId uint) (*Statistic, error) {
	stats, err := s.repo.GetStatisticByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statistics for user %d", userId)
	}

	return stats, nil
}

func (s *ServiceImpl) GetTopXAmongUserIdsByMeasure(ctx context.Context, topX int, userIds []uint, measure Measure) ([]uint, []int, error) {
	var topXUserIds []uint
	var values []int

	switch measure {
	case MeasureWins:
		ids, wins, err := s.repo.GetTopXAmongUserIdsByWins(ctx, topX, userIds)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get top %d users by wins", topX)
		}

		topXUserIds = ids
		values = wins
	case MeasureStreak:
		ids, streaks, err := s.repo.GetTopXAmongUserIdsByStreak(ctx, topX, userIds)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get top %d users by win streaks", topX)
		}

		topXUserIds = ids
		values = streaks
	default:
		return nil, nil, errors.Errorf("unsupported measure: %s", measure)
	}

	return topXUserIds, values, nil
}

func (s *ServiceImpl) CreateStatistic(ctx context.Context, userId uint) error {
	if err := s.repo.CreateStatistic(ctx, userId); err != nil {
		return errors.Wrapf(err, "failed to create statistics for user %d", userId)
	}

	return nil
}

func (s *ServiceImpl) UpdateStatisticsByUserIds(ctx context.Context, userIds []uint, result MatchResult) error {
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

func (s *ServiceImpl) TransferStatistics(ctx context.Context, fromUserId, toUserId uint) error {
	fromStats, err := s.repo.GetStatisticByUserId(ctx, fromUserId)
	if err != nil {
		return errors.Wrapf(err, "failed to get statistics for user %d", fromUserId)
	}

	toStats, err := s.repo.GetStatisticByUserId(ctx, toUserId)
	if err != nil {
		return errors.Wrapf(err, "failed to get statistics for user %d", toUserId)
	}

	fromStats.UserId = toUserId
	toStats.UserId = fromUserId

	if err := s.repo.UpdateStatistics(ctx, []Statistic{*fromStats, *toStats}); err != nil {
		return errors.Wrap(err, "failed to update statistics")
	}

	return nil
}
