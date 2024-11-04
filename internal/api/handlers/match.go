package handlers

import (
	"core/internal/api/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type postMatchRequest struct {
	ClubID   int      `json:"clubId" validate:"required"`
	GameID   int      `json:"gameId" validate:"required"`
	TeamsIDs []int    `json:"teamsIds" validate:"required"`
	Sets     []string `json:"sets" validate:"required"`
}

func (h *Handler) PostMatch(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[postMatchRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	_, err = h.matchService.CreateMatch(ctx, req.ClubID, req.GameID, req.TeamsIDs, req.Sets)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}
