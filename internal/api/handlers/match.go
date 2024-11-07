package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type postMatchRequest struct {
	ClubID   int      `json:"clubId" validate:"required"`
	GameID   int      `json:"gameId" validate:"required"`
	TeamsIDs []int    `json:"teamsIds" validate:"required"`
	Sets     []string `json:"sets" validate:"required"`
}

type postMatchResponse struct {
	MatchID int `json:"matchId"`
}

func (h *Handler) PostMatch(ctx context.Context, req *postMatchRequest) (*postMatchResponse, error) {
	matchID, err := h.matchService.CreateMatch(ctx, req.ClubID, req.GameID, req.TeamsIDs, req.Sets)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to create match, try again later")
	}

	return &postMatchResponse{MatchID: matchID}, nil
}
