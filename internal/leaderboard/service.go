package leaderboard

import (
	"context"
	"core/internal/club"
	"core/internal/rating"
	"core/internal/statistic"
	"core/internal/user"

	"github.com/pkg/errors"
)

type Service interface {
	GetLeaderboard(ctx context.Context, clubId uint, topX int, leaderboardType LeaderboardType) (*Leaderboard, error)
}

type ServiceImpl struct {
	clubService      club.Service
	userService      user.Service
	ratingService    rating.Service
	statisticService statistic.Service
}

func NewService(clubService club.Service, userService user.Service, ratingService rating.Service, statisticService statistic.Service) Service {
	return &ServiceImpl{
		clubService:      clubService,
		userService:      userService,
		ratingService:    ratingService,
		statisticService: statisticService,
	}
}

func (s *ServiceImpl) GetLeaderboard(ctx context.Context, clubId uint, topX int, leaderboardType LeaderboardType) (*Leaderboard, error) {
	var userIds []uint
	var values []float64

	userIdsInClub, err := s.clubService.GetUserIdsInClub(ctx, clubId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get users in Club %d", clubId)
	}

	switch leaderboardType {
	case TypeWins:
		ids, wins, err := s.statisticService.GetTopXAmongUserIdsByMeasure(ctx, topX, userIdsInClub, statistic.MeasureWins)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIds by wins", topX)
		}

		userIds = ids
		values = s.convertIntToFloat64(wins)
	case TypeStreak:
		ids, winstreaks, err := s.statisticService.GetTopXAmongUserIdsByMeasure(ctx, topX, userIdsInClub, statistic.MeasureStreak)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIds by streak", topX)
		}

		userIds = ids
		values = s.convertIntToFloat64(winstreaks)
	case TypeRating:
		ids, ratings, err := s.ratingService.GetTopXAmongUserIdsByRating(ctx, topX, userIdsInClub)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIds by rating", topX)
		}

		userIds = ids
		values = s.convertIntToFloat64(ratings)
	default:
		return nil, errors.Errorf("unknown leaderboard type: %s", leaderboardType)
	}

	users, err := s.userService.GetUsers(ctx, userIds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}

	entries := make([]Entry, len(users))
	for i, u := range users {
		entries[i] = Entry{
			Value:  values[i],
			UserId: u.Id,
			Name:   u.Name,
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get leaderboard")
	}

	lboard := &Leaderboard{
		Type:    leaderboardType,
		Entries: entries,
	}

	return lboard, nil
}

func (s *ServiceImpl) convertIntToFloat64(values []int) []float64 {
	floats := make([]float64, len(values))
	for i, v := range values {
		floats[i] = float64(v)
	}
	return floats
}
