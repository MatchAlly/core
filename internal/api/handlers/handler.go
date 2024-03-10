package handlers

import (
	"core/internal/authentication"
	"core/internal/club"
	"core/internal/leaderboard"
	"core/internal/match"
	"core/internal/rating"
	"core/internal/statistic"
	"core/internal/user"

	"go.uber.org/zap"
)

type Handler struct {
	logger             *zap.SugaredLogger
	authService        authentication.Service
	userService        user.Service
	clubService        club.Service
	matchService       match.Service
	ratingService      rating.Service
	statisticService   statistic.Service
	leaderboardService leaderboard.Service
}

func NewHandler(
	logger *zap.SugaredLogger,
	authService authentication.Service,
	userService user.Service,
	clubService club.Service,
	matchService match.Service,
	ratingService rating.Service,
	statisticService statistic.Service,
	leaderboardService leaderboard.Service,
) *Handler {
	return &Handler{
		logger:             logger,
		authService:        authService,
		userService:        userService,
		clubService:        clubService,
		matchService:       matchService,
		ratingService:      ratingService,
		statisticService:   statisticService,
		leaderboardService: leaderboardService,
	}
}
