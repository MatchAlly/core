package api

import (
	"core/internal/api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func addRoutes(api huma.API, h *handlers.Handler) {
	// Authentication
	huma.Post(api, "/auth/signup", h.Signup) // Needs no auth middleware
	huma.Post(api, "/auth/login", h.Login)   // Needs no auth middleware
	huma.Post(api, "/auth/refresh", h.Refresh)
	huma.Post(api, "/auth/password", h.ChangePassword)

	// Users
	huma.Delete(api, "/users", h.DeleteUser)
	huma.Put(api, "/users", h.UpdateUser)

	// Clubs
	huma.Get(api, "/clubs", h.GetMemberships)
	huma.Post(api, "/clubs", h.CreateClub)
	huma.Put(api, "/clubs/:clubId", h.UpdateClub)
	huma.Delete(api, "/clubs/:clubId", h.DeleteClub)
	huma.Get(api, "/clubs/:clubId/members", h.GetMembersInClub)
	huma.Delete(api, "/clubs/:clubId/members/:memberId", h.RemoveMemberFromClub)
	huma.Put(api, "/clubs/:clubId/members/:memberId", h.UpdateMemberRole)
	huma.Post(api, "/clubs/:clubId/matches", h.PostMatch)
}
