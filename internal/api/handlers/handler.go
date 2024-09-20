package handlers

import (
	"core/internal/authentication"
	"core/internal/club"
	"core/internal/match"
	"core/internal/member"
	"core/internal/rating"
	"core/internal/user"

	"go.uber.org/zap"
)

type Handler struct {
	l             *zap.SugaredLogger
	authService   authentication.Service
	userService   user.Service
	clubService   club.Service
	memberService member.Service
	matchService  match.Service
	ratingService rating.Service
}

func NewHandler(
	l *zap.SugaredLogger,
	authService authentication.Service,
	userService user.Service,
	clubService club.Service,
	memberService member.Service,
	matchService match.Service,
	ratingService rating.Service,
) *Handler {
	return &Handler{
		l:             l,
		authService:   authService,
		userService:   userService,
		clubService:   clubService,
		memberService: memberService,
		matchService:  matchService,
		ratingService: ratingService,
	}
}
