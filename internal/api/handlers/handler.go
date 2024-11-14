package handlers

import (
	"core/internal/authentication"
	"core/internal/authorization"
	"core/internal/club"
	"core/internal/match"
	"core/internal/member"
	"core/internal/rating"
	"core/internal/user"

	"go.uber.org/zap"
)

type Handler struct {
	l             *zap.SugaredLogger
	authNService  authentication.Service
	authZService  authorization.Service
	userService   user.Service
	clubService   club.Service
	memberService member.Service
	matchService  match.Service
	ratingService rating.Service
}

func NewHandler(
	l *zap.SugaredLogger,
	authService authentication.Service,
	authZService authorization.Service,
	userService user.Service,
	clubService club.Service,
	memberService member.Service,
	matchService match.Service,
	ratingService rating.Service,
) *Handler {
	return &Handler{
		l:             l,
		authNService:  authService,
		authZService:  authZService,
		userService:   userService,
		clubService:   clubService,
		memberService: memberService,
		matchService:  matchService,
		ratingService: ratingService,
	}
}
