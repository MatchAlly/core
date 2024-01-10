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
	authGroup := e.Group("/auth")
	authGroup.POST("/signup", h.Signup)
	authGroup.POST("/login", h.Login)
	authGroup.POST("/refresh", h.Refresh)

	// Users
	userGroup := e.Group("/user", authGuard)
	userGroup.DELETE("", authHandler(h.DeleteUser))
	userGroup.GET("/invites", authHandler(h.GetUserInvites))
	//userGroup.POST("/invites/:inviteId", authHandler(h.RespondToInvite))

	// Clubs
	clubGroup := e.Group("/club", authGuard)
	clubGroup.POST("", authHandler(h.CreateClub))
	clubGroup.PUT("", authHandler(h.UpdateClub))
	clubGroup.DELETE("", authHandler(h.DeleteClub))
	clubGroup.GET("/users", authHandler(h.GetUsersInClub))
	clubGroup.POST("/invite", authHandler(h.InviteUsersToClub))
	clubGroup.POST("/users/virtual", authHandler(h.AddVirtualUserToClub))
	//clubGroup.POST("/users/:userId/virtual/:virtualUserId", authHandler(h.TransferVirtualUserToUser))
	clubGroup.DELETE("/users/:userId", authHandler(h.RemoveUserFromClub))
	clubGroup.PUT("/users/:userId", authHandler(h.UpdateUserRole))
	clubGroup.GET("/top/:topX/measures/:leaderboardType", authHandler(h.GetLeaderboard))
	clubGroup.POST("/matches", authHandler(h.PostMatch))
}
