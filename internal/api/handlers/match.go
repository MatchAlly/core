package handlers

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type postClubMatchRequest struct {
	Body struct {
		ClubID int                        `json:"clubId" minimum:"1"`
		GameID int                        `json:"gameId" minimum:"1"`
		Teams  []postClubMatchRequestTeam `json:"teams" minItems:"1"`
		Sets   []string                   `json:"sets,omitempty"`
	}
}

type postClubMatchRequestTeam struct {
	Members []int `json:"members" minItems:"1"`
}

type postClubMatchResponse struct {
	Body struct {
		MatchID int `json:"matchId"`
	}
}

func (h *Handler) PostClubMatch(ctx context.Context, req *postClubMatchRequest) (*postClubMatchResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authZService.IsMember(ctx, userID, req.Body.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to get matches in this club")
	}

	tempTeams := make([][]int, len(req.Body.Teams))
	for i, t := range req.Body.Teams {
		tempTeams[i] = t.Members
	}

	teams, err := h.matchService.GetOrCreateTeams(ctx, req.Body.ClubID, tempTeams)
	if err != nil {
		h.l.Error("failed to get or create teams", "error", err)
		return nil, huma.Error500InternalServerError("failed to get or create teams, try again later")
	}

	matchID, err := h.matchService.CreateMatch(ctx, req.Body.ClubID, req.Body.GameID, teams, req.Body.Sets)
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
	ClubID int `path:"clubId" minimum:"1"`
	GameID int `query:"gameId" required:"false" minimum:"1"`
}

type getClubMatchesResponse struct {
	Body struct {
		Matches []getClubMatchesResponseMatch `json:"matches"`
	}
}

type getClubMatchesResponseMatch struct {
	ID     int                          `json:"id"`
	GameID int                          `json:"game_id"`
	Sets   []string                     `json:"sets,omitempty"`
	Teams  []getClubMatchesResponseTeam `json:"teams"`
	Date   time.Time                    `json:"date"`
}

type getClubMatchesResponseTeam struct {
	ID      int                                `json:"id"`
	Members []getClubMatchesResponseTeamMember `json:"members"`
}

type getClubMatchesResponseTeamMember struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetClubMatches(ctx context.Context, req *getClubMatchesRequest) (*getClubMatchesResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authZService.IsMember(ctx, userID, req.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to get matches in this club")
	}

	var gameID *int
	if req.GameID != 0 {
		gameID = &req.GameID
	}
	matches, err := h.matchService.GetMatches(ctx, req.ClubID, gameID)
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
					ID:   mem.ID,
					Name: mem.DisplayName,
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
