//go:generate mockgen --source=service.go -destination=service_mock.go -package=rating
package rating

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetTopXAmongUserIdsByRating(ctx context.Context, topX int, userIds []uint) (topXUserIds []uint, ratings []int, err error)
	CreateRating(ctx context.Context, userId uint) error
	UpdateRatings(ctx context.Context, draw bool, winningUserIds, losingUserIds []uint) error
	TransferRatings(ctx context.Context, fromUserId, toUserId uint) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) GetTopXAmongUserIdsByRating(ctx context.Context, topX int, userIds []uint) ([]uint, []int, error) {
	userIds, ratings, err := s.repo.GetTopXAmongUserIdsByRating(ctx, topX, userIds)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get top %d user ids by rating", topX)
	}

	return userIds, ratings, nil
}

func (s *ServiceImpl) CreateRating(ctx context.Context, userId uint) error {
	rating := &Rating{
		UserId:     userId,
		Value:      startRating,
		Deviation:  maxDeviation,
		Volatility: startVolatility,
	}

	if err := s.repo.CreateRating(ctx, rating); err != nil {
		return errors.Wrap(err, "failed to create rating")
	}

	return nil
}

func (s *ServiceImpl) UpdateRatings(ctx context.Context, draw bool, winningUserIds, losingUserIds []uint) error {
	var updatedRatings []Rating

	winnerRatings, err := s.repo.GetRatingsByUserIds(ctx, winningUserIds)
	if err != nil {
		return errors.Wrap(err, "failed to get winning users")
	}

	loserRatings, err := s.repo.GetRatingsByUserIds(ctx, losingUserIds)
	if err != nil {
		return errors.Wrap(err, "failed to get losing users")
	}

	winnerAverageRating, winnerAverageDeviation := s.getAverageRatingAndDeviation(winnerRatings)
	loserAverageRating, loserAverageDeviation := s.getAverageRatingAndDeviation(loserRatings)

	var winnerResult, loserResult MatchResult

	if draw {
		winnerResult = MatchResult{
			OpponentRating:    loserAverageRating,
			OpponentDeviation: loserAverageDeviation,
			Result:            resultMultiplierDraw,
		}

		loserResult = MatchResult{
			OpponentRating:    winnerAverageRating,
			OpponentDeviation: winnerAverageDeviation,
			Result:            resultMultiplierDraw,
		}
	} else {
		winnerResult = MatchResult{
			OpponentRating:    loserAverageRating,
			OpponentDeviation: loserAverageDeviation,
			Result:            resultMultiplierWin,
		}

		loserResult = MatchResult{
			OpponentRating:    winnerAverageRating,
			OpponentDeviation: winnerAverageDeviation,
			Result:            resultMultiplierLoss,
		}
	}

	for _, winnerRating := range winnerRatings {
		updatedRating := ApplyActiveRatingPeriod(winnerRating, []MatchResult{winnerResult, loserResult})
		updatedRatings = append(updatedRatings, updatedRating)
	}

	for _, loserRating := range loserRatings {
		updatedRating := ApplyActiveRatingPeriod(loserRating, []MatchResult{loserResult, winnerResult})
		updatedRatings = append(updatedRatings, updatedRating)
	}

	if err := s.repo.UpdateRatings(ctx, updatedRatings); err != nil {
		return errors.Wrap(err, "failed to update ratings")
	}

	return nil
}

func (s *ServiceImpl) TransferRatings(ctx context.Context, fromUserId, toUserId uint) error {
	fromUserRating, err := s.repo.GetRatingByUserId(ctx, fromUserId)
	if err != nil {
		return errors.Wrap(err, "failed to get from user rating")
	}

	toUserRating, err := s.repo.GetRatingByUserId(ctx, toUserId)
	if err != nil {
		return errors.Wrap(err, "failed to get to user rating")
	}

	fromUserRating.UserId = toUserId
	toUserRating.UserId = fromUserId

	transferedRatings := []Rating{*fromUserRating, *toUserRating}
	if err := s.repo.UpdateRatings(ctx, transferedRatings); err != nil {
		return errors.Wrap(err, "failed to update ratings")
	}

	return nil
}

func (s *ServiceImpl) getAverageRatingAndDeviation(ratings []Rating) (float64, float64) {
	var totalRating, totalDeviation float64

	for _, rating := range ratings {
		totalRating += rating.Value
		totalDeviation += rating.Deviation
	}

	return totalRating / float64(len(ratings)), totalDeviation / float64(len(ratings))
}
