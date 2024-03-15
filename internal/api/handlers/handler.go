package handlers

import (
	"core/internal/authentication"
	"core/internal/club"
	"core/internal/invite"
	"core/internal/match"
	"core/internal/rating"
	"core/internal/statistic"
	"core/internal/user"

	"go.uber.org/zap"
)

type Handler struct {
	l                *zap.SugaredLogger
	authService      authentication.Service
	userService      user.Service
	clubService      club.Service
	matchService     match.Service
	ratingService    rating.Service
	statisticService statistic.Service
	inviteService    invite.Service
}

func NewHandler(
	l *zap.SugaredLogger,
	authService authentication.Service,
	userService user.Service,
	clubService club.Service,
	matchService match.Service,
	ratingService rating.Service,
	statisticService statistic.Service,
	inviteService invite.Service,
) *Handler {
	return &Handler{
		l:                l,
		authService:      authService,
		userService:      userService,
		clubService:      clubService,
		matchService:     matchService,
		ratingService:    ratingService,
		statisticService: statisticService,
	}
}
