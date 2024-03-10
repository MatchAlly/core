package api

import (
	"core/internal/api/handlers"
	"core/internal/api/helpers"
	"core/internal/api/middleware"
	"core/internal/authentication"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Register(h *handlers.Handler, e *echo.Group, logger *zap.SugaredLogger, authService authentication.Service) {
	authCtx := helpers.AuthenticatedContextFactory(logger)
	authGuard := middleware.AuthGuard(authService)

	// Authentication
	auth := e.Group("/auth")
	auth.POST("/signup", h.Signup)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh, authGuard)

	// Users
	users := e.Group("/user", authGuard)
	users.DELETE("", authCtx(h.DeleteUser))
	users.PUT("", authCtx(h.UpdateUser))
	users.GET("/invites", authCtx(h.GetUserInvites))
	users.POST("/invites/:inviteId", authCtx(h.RespondToInvite))

	// Clubs
	clubs := e.Group("/club", authGuard)
	clubs.POST("", authCtx(h.CreateClub))
	clubs.PUT("", authCtx(h.UpdateClub))
	clubs.DELETE("", authCtx(h.DeleteClub))
	clubs.GET("/users", authCtx(h.GetUsersInClub))
	clubs.DELETE("/users/:userId", authCtx(h.RemoveUserFromClub))
	clubs.PUT("/users/:userId", authCtx(h.UpdateUserRole))
	clubs.POST("/invites", authCtx(h.InviteUsersToClub))
	clubs.GET("/leaderboards", authCtx(h.GetLeaderboard))
	clubs.POST("/matches", authCtx(h.PostMatch))
}
