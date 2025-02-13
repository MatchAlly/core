package handlers

import (
	"core/internal/authentication"
	"core/internal/authorization"
	"core/internal/club"
	"core/internal/game"
	"core/internal/match"
	"core/internal/member"
	"core/internal/rating"
	"core/internal/subscription"
	"core/internal/user"
	"time"

	"go.uber.org/zap"
)

type Config struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}
type Handler struct {
	l                   *zap.SugaredLogger
	config              Config
	authNService        authentication.Service
	authZService        authorization.Service
	userService         user.Service
	clubService         club.Service
	memberService       member.Service
	matchService        match.Service
	ratingService       rating.Service
	gameService         game.Service
	subscriptionService subscription.Service
}

func NewHandler(
	l *zap.SugaredLogger,
	config Config,
	authService authentication.Service,
	authZService authorization.Service,
	userService user.Service,
	clubService club.Service,
	memberService member.Service,
	matchService match.Service,
	ratingService rating.Service,
	gameService game.Service,
	subscriptionService subscription.Service,
) *Handler {
	return &Handler{
		l:                   l,
		config:              config,
		authNService:        authService,
		authZService:        authZService,
		userService:         userService,
		clubService:         clubService,
		memberService:       memberService,
		matchService:        matchService,
		ratingService:       ratingService,
		gameService:         gameService,
		subscriptionService: subscriptionService,
	}
}
