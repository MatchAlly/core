package api

import (
	"core/internal/api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func addPublicRoutes(g *huma.Group, h *handlers.Handler) {
	// Authentication
	huma.Post(g, "/auth/signup", h.Signup)
	huma.Post(g, "/auth/login", h.Login)
}

func addAuthRoutes(g *huma.Group, h *handlers.Handler) {
	// Authentication
	huma.Post(g, "/auth/logout", h.Logout)
	huma.Post(g, "/auth/refresh", h.Refresh)
	huma.Post(g, "/auth/password", h.ChangePassword)

	// Users
	huma.Delete(g, "/users/:userId", h.DeleteUser)
	huma.Put(g, "/users/:userId", h.UpdateUser)
	huma.Get(g, "/users/:userId/clubs", h.GetMemberships)

	// Clubs
	huma.Post(g, "/clubs", h.CreateClub)
	huma.Put(g, "/clubs/:clubId", h.UpdateClub)
	huma.Delete(g, "/clubs/:clubId", h.DeleteClub)
	huma.Get(g, "/clubs/:clubId/members", h.GetMembersInClub)
	huma.Delete(g, "/clubs/:clubId/members/:memberId", h.RemoveMemberFromClub)
	huma.Put(g, "/clubs/:clubId/members/:memberId", h.UpdateMemberRole)
	huma.Post(g, "/clubs/:clubId/matches", h.PostClubMatch)
	huma.Get(g, "/clubs/:clubId/matches", h.GetClubMatches)
	huma.Get(g, "/clubs/:clubId/games", h.GetClubGames)
	huma.Post(g, "/clubs/:clubId/games", h.PostClubGame)
}
