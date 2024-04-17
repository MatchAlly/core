package statistic

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetStatisticsMemberId(ctx context.Context, memberId uint) ([]Statistic, error)
	CreateDefaultStatistic(ctx context.Context, memberId, gameId uint) error
	UpdateGameStatisticsByMemberIds(ctx context.Context, memberIds []uint, gameId uint, result MatchResult) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetStatisticsMemberId(ctx context.Context, memberId uint) ([]Statistic, error) {
	stats, err := s.repo.GetStatisticsByMemberId(ctx, memberId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statistics for user %d", memberId)
	}

	return stats, nil
}

func (s *service) CreateDefaultStatistic(ctx context.Context, userId, gameId uint) error {
	if err := s.repo.CreateDefaultStatistic(ctx, userId, gameId); err != nil {
		return errors.Wrapf(err, "failed to create default statistics for user %d", userId)
	}

	return nil
}

func (s *service) UpdateGameStatisticsByMemberIds(ctx context.Context, memberIds []uint, gameId uint, result MatchResult) error {
	oldStatistics, err := s.repo.GetGameStatisticsForMemberIds(ctx, memberIds, gameId)
	if err != nil {
		return errors.Wrapf(err, "failed to get statistics for members %v", memberIds)
	}

	updatedStatistics := make([]Statistic, len(memberIds))
	for i := range memberIds {
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

		updatedStatistics[i] = stats
	}

	if err := s.repo.UpdateStatistics(ctx, updatedStatistics); err != nil {
		return errors.Wrap(err, "failed to update statistics")
	}

	return nil
}
