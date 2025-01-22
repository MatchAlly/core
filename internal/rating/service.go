package rating

import (
	"context"
	"fmt"

	"github.com/Sebsh1/openskill.go"
	"github.com/pkg/errors"
)

type Service interface {
	CreateRating(ctx context.Context, memberID, gameID int) (int, error)
	UpdateRatingsByRanks(ctx context.Context, teamsByMemberIDs [][]int, ranks []int) error
}

type service struct {
	repo  Repository
	rater openskill.Rater
}

func NewService(repo Repository) Service {
	return &service{
		repo:  repo,
		rater: openskill.DefaultPlackettLuceModel(),
	}
}

func (s *service) CreateRating(ctx context.Context, memberID, gameID int) (int, error) {
	rating := &Rating{
		MemberID: memberID,
		GameID:   gameID,
		Mu:       startMu,
		Sigma:    startSigma,
	}

	id, err := s.repo.CreateRating(ctx, rating)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create rating")
	}

	return id, nil
}

func (s *service) UpdateRatingsByRanks(ctx context.Context, teamsByMemberIDs [][]int, ranks []int) error {
	ids, shape := flatten(teamsByMemberIDs)
	ratings, err := s.repo.GetRatingsByMemberIds(ctx, ids)
	if err != nil {
		return errors.Wrap(err, "failed to get ratings")
	}

	openSkillRatings := make([]openskill.Rating, len(ratings))
	for i, rating := range ratings {
		openSkillRatings[i] = openskill.Rating{
			Mu:    rating.Mu,
			Sigma: rating.Sigma,
		}
	}

	teams, err := unflatten(openSkillRatings, shape)
	if err != nil {
		return errors.Wrap(err, "failed to unflatten ratings")
	}

	updatedRatings, err := s.rater.Rate(teams, ranks, nil, nil)
	if err != nil {
		return errors.Wrap(err, "failed to calculate new ratings")
	}

	newRatings, _ := flatten(updatedRatings)
	for i := range ratings {
		ratings[i].Mu = newRatings[i].Mu
		ratings[i].Sigma = newRatings[i].Sigma
	}

	if err := s.repo.UpdateRatings(ctx, ratings); err != nil {
		return errors.Wrap(err, "failed to update ratings")
	}

	return nil
}

func flatten[T any](matrix [][]T) ([]T, []int) {
	shape := make([]int, len(matrix))
	totalLen := 0
	flattened := make([]T, 0, totalLen)

	for i, row := range matrix {
		shape[i] = len(row)
		totalLen += len(row)
		flattened = append(flattened, row...)
	}

	return flattened, shape
}

func unflatten[T any](flattened []T, shape []int) ([][]T, error) {
	total := 0
	for _, v := range shape {
		total += v
	}
	if len(flattened) != total {
		return nil, fmt.Errorf("total length of shape does not match length of flattened list")
	}

	var matrix [][]T
	start := 0
	for _, rowLen := range shape {
		if rowLen < 0 {
			return nil, fmt.Errorf("invalid row length: %d", rowLen)
		}
		end := start + rowLen
		matrix = append(matrix, flattened[start:end])
		start = end
	}

	return matrix, nil
}
