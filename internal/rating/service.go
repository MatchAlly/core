//go:generate mockgen --source=service.go -destination=service_mock.go -package=rating
package rating

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetTopMembersByRating(ctx context.Context, topX int, bemberIds []uint) (topXMemberIds []uint, ratings []int, err error)
	CreateRating(ctx context.Context, memberId uint) error
	UpdateRatings(ctx context.Context, draw bool, winningMemberIds, losingMemberIds []uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetTopMembersByRating(ctx context.Context, topX int, memberIds []uint) ([]uint, []int, error) {
	topMemberIds, ratings, err := s.repo.GetTopMembersByRating(ctx, topX, memberIds)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get top %d member ids by rating", topX)
	}

	return topMemberIds, ratings, nil
}

func (s *service) CreateRating(ctx context.Context, memberId uint) error {
	rating := &Rating{
		MemberId:   memberId,
		Value:      startRating,
		Deviation:  maxDeviation,
		Volatility: startVolatility,
	}

	if err := s.repo.CreateRating(ctx, rating); err != nil {
		return errors.Wrap(err, "failed to create rating")
	}

	return nil
}

func (s *service) UpdateRatings(ctx context.Context, draw bool, winningMemberIds, losingMemberIds []uint) error {
	var updatedRatings []Rating

	winnerRatings, err := s.repo.GetRatingsByMemberIds(ctx, winningMemberIds)
	if err != nil {
		return errors.Wrap(err, "failed to get winning members")
	}

	loserRatings, err := s.repo.GetRatingsByMemberIds(ctx, losingMemberIds)
	if err != nil {
		return errors.Wrap(err, "failed to get losing members")
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

func (s *service) getAverageRatingAndDeviation(ratings []Rating) (float64, float64) {
	var totalRating, totalDeviation float64

	for _, rating := range ratings {
		totalRating += rating.Value
		totalDeviation += rating.Deviation
	}

	return totalRating / float64(len(ratings)), totalDeviation / float64(len(ratings))
}
