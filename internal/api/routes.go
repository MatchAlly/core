package api

import (
	"core/internal/api/handlers"
	"core/internal/api/helpers"
	"core/internal/api/middleware"
	"core/internal/authentication"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Register(h *handlers.Handler, e *echo.Group, l *zap.SugaredLogger, authService authentication.Service) {
	authCtx := helpers.AuthenticatedContextFactory(l)
	authGuard := middleware.AuthGuard(authService)

	// Authentication
	auth := e.Group("/auth")
	auth.POST("/signup", h.Signup)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh, authGuard)
	auth.POST("/password", authCtx(h.ChangePassword))

	// Users
	users := e.Group("/users", authGuard)
	users.DELETE("", authCtx(h.DeleteUser))
	users.PUT("", authCtx(h.UpdateUser))

	// Clubs
	clubs := e.Group("/clubs", authGuard)
	clubs.GET("", authCtx(h.GetMemberships))
	clubs.POST("", authCtx(h.CreateClub))
	clubs.PUT(":clubId", authCtx(h.UpdateClub))
	clubs.DELETE(":clubId", authCtx(h.DeleteClub))
	clubs.GET(":clubId/members", authCtx(h.GetMembersInClub))
	clubs.DELETE(":clubId/members/:memberId", authCtx(h.RemoveMemberFromClub))
	clubs.PUT(":clubId/members/:memberId", authCtx(h.UpdateMemberRole))
	clubs.POST(":clubId/matches", authCtx(h.PostMatch))
}
