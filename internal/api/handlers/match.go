package handlers

import (
	"context"
	"core/internal/game"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type postClubMatchRequest struct {
	Body struct {
		ClubID uuid.UUID                  `json:"clubId"`
		GameID uuid.UUID                  `json:"gameId"`
		Mode   string                     `json:"mode" enum:"FREE_FOR_ALL,TEAM,COOP"`
		Teams  []postClubMatchRequestTeam `json:"teams" minItems:"1"`
		Sets   []string                   `json:"sets,omitempty"`
	}
}

type postClubMatchRequestTeam struct {
	Members []uuid.UUID `json:"members" minItems:"1"`
}

type postClubMatchResponse struct {
	Body struct {
		MatchID uuid.UUID `json:"matchId"`
	}
}

func (h *Handler) PostClubMatch(ctx context.Context, req *postClubMatchRequest) (*postClubMatchResponse, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authorization.IsMember(ctx, userID, req.Body.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to get matches in this club")
	}

	tempTeams := make([][]uuid.UUID, len(req.Body.Teams))
	for i, t := range req.Body.Teams {
		tempTeams[i] = t.Members
	}

	teams, err := h.match.GetOrCreateTeams(ctx, req.Body.ClubID, tempTeams)
	if err != nil {
		h.l.Error("failed to get or create teams", "error", err)
		return nil, huma.Error500InternalServerError("failed to get or create teams, try again later")
	}

	var mode game.Mode
	switch req.Body.Mode {
	case "FREE_FOR_ALL":
		mode = game.ModeFreeForAll
	case "TEAM":
		mode = game.ModeTeam
	case "COOP":
		mode = game.ModeCoop
	default:
		return nil, huma.Error400BadRequest("invalid game mode")
	}

	matchID, err := h.match.CreateMatch(ctx, req.Body.ClubID, req.Body.GameID, teams, req.Body.Sets, mode)
	if err != nil {
		h.l.Error("failed to create match", "error", err)
		return nil, huma.Error500InternalServerError("failed to create match, try again later")
	}

	// TODO update statistics and rankings

	resp := &postClubMatchResponse{}
	resp.Body.MatchID = matchID

	return resp, nil
}

type getClubMatchesRequest struct {
	ClubID uuid.UUID  `path:"clubId"`
	GameID *uuid.UUID `query:"gameId" required:"false"`
}

type getClubMatchesResponse struct {
	Body struct {
		Matches []getClubMatchesResponseMatch `json:"matches"`
	}
}

type getClubMatchesResponseMatch struct {
	ID     uuid.UUID                    `json:"id"`
	GameID uuid.UUID                    `json:"game_id"`
	Sets   []string                     `json:"sets,omitempty"`
	Teams  []getClubMatchesResponseTeam `json:"teams"`
	Date   time.Time                    `json:"date"`
}

type getClubMatchesResponseTeam struct {
	ID      uuid.UUID                          `json:"id"`
	Members []getClubMatchesResponseTeamMember `json:"members"`
}

type getClubMatchesResponseTeamMember struct {
	ID uuid.UUID `json:"id"`
}

func (h *Handler) GetClubMatches(ctx context.Context, req *getClubMatchesRequest) (*getClubMatchesResponse, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authorization.IsMember(ctx, userID, req.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to get matches in this club")
	}

	var gameID *uuid.UUID
	if req.GameID != nil {
		gameID = req.GameID
	}
	matches, err := h.match.GetMatches(ctx, req.ClubID, gameID)
	if err != nil {
		h.l.Error("failed to get matches", "error", err)
		return nil, huma.Error500InternalServerError("failed to get matches, try again later")
	}

	mappedMatches := make([]getClubMatchesResponseMatch, len(matches))
	for i, m := range matches {
		teams := make([]getClubMatchesResponseTeam, len(m.Teams))
		for j, t := range m.Teams {
			members := make([]getClubMatchesResponseTeamMember, len(t.Members))
			for k, mem := range t.Members {
				members[k] = getClubMatchesResponseTeamMember{
					ID: mem.ID,
				}
			}
			teams[j] = getClubMatchesResponseTeam{
				ID:      t.ID,
				Members: members,
			}
		}

		mappedMatches[i] = getClubMatchesResponseMatch{
			ID:     m.ID,
			GameID: m.GameID,
			Sets:   m.Sets,
			Teams:  teams,
			Date:   m.CreatedAt,
		}
	}

	resp := &getClubMatchesResponse{}
	resp.Body.Matches = mappedMatches

	return resp, nil
}
