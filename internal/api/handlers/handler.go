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
	"log/slog"
	"time"
)

type Config struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}
type Handler struct {
	l              *slog.Logger
	config         Config
	authentication authentication.Service
	authorization  authorization.Service
	user           user.Service
	club           club.Service
	member         member.Service
	match          match.Service
	rating         rating.Service
	game           game.Service
	subscription   subscription.Service
}

func NewHandler(
	l *slog.Logger,
	config Config,
	authentication authentication.Service,
	authorization authorization.Service,
	user user.Service,
	club club.Service,
	member member.Service,
	match match.Service,
	rating rating.Service,
	game game.Service,
	subscription subscription.Service,
) *Handler {
	return &Handler{
		l:              l,
		config:         config,
		authentication: authentication,
		authorization:  authorization,
		user:           user,
		club:           club,
		member:         member,
		match:          match,
		rating:         rating,
		game:           game,
		subscription:   subscription,
	}
}
