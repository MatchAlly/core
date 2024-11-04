package rating

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	CreateRating(ctx context.Context, memberID, gameID int) (int, error)
	UpdateRatings(ctx context.Context, draw bool, winningMemberIds, losingMemberIds []int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateRating(ctx context.Context, memberID, gameID int) (int, error) {
	rating := &Rating{
		MemberID:   memberID,
		GameID:     gameID,
		Value:      startRating,
		Deviation:  maxDeviation,
		Volatility: startVolatility,
	}

	id, err := s.repo.CreateRating(ctx, rating)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create rating")
	}

	return id, nil
}

func (s *service) UpdateRatings(ctx context.Context, draw bool, winningMemberIds, losingMemberIds []int) error {
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
