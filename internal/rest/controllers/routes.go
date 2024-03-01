package controllers

import (
	"core/internal/authentication"
	"core/internal/club"
	"core/internal/leaderboard"
	"core/internal/match"
	"core/internal/rating"
	"core/internal/rest/handlers"
	"core/internal/rest/middleware"
	"core/internal/statistic"
	"core/internal/user"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handlers struct {
	logger             *zap.SugaredLogger
	authService        authentication.Service
	userService        user.Service
	clubService        club.Service
	matchService       match.Service
	ratingService      rating.Service
	statisticService   statistic.Service
	leaderboardService leaderboard.Service
}

func Register(
	e *echo.Group,
	logger *zap.SugaredLogger,
	authService authentication.Service,
	userService user.Service,
	clubService club.Service,
	matchService match.Service,
	ratingService rating.Service,
	statisticService statistic.Service,
	leaderboardService leaderboard.Service,
) {
	h := &Handlers{
		logger:             logger,
		authService:        authService,
		userService:        userService,
		clubService:        clubService,
		matchService:       matchService,
		ratingService:      ratingService,
		statisticService:   statisticService,
		leaderboardService: leaderboardService,
	}

	authHandler := handlers.AuthenticatedHandlerFactory(logger)

	authGuard := middleware.AuthGuard(authService)

	// Authentication
	auth := e.Group("/auth")
	auth.POST("/signup", h.Signup)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh, authGuard)

	// Users
	users := e.Group("/user", authGuard)
	users.DELETE("", authHandler(h.DeleteUser))
	users.PUT("", authHandler(h.UpdateUser))
	users.GET("/invites", authHandler(h.GetUserInvites))
	users.POST("/invites/:inviteId", authHandler(h.RespondToInvite))

	// Clubs
	clubs := e.Group("/club", authGuard)
	clubs.POST("", authHandler(h.CreateClub))
	clubs.PUT("", authHandler(h.UpdateClub))
	clubs.DELETE("", authHandler(h.DeleteClub))
	clubs.GET("/users", authHandler(h.GetUsersInClub))
	clubs.DELETE("/users/:userId", authHandler(h.RemoveUserFromClub))
	clubs.PUT("/users/:userId", authHandler(h.UpdateUserRole))
	clubs.POST("/invites", authHandler(h.InviteUsersToClub))
	clubs.GET("/leaderboards", authHandler(h.GetLeaderboard))
	clubs.POST("/matches", authHandler(h.PostMatch))
}
