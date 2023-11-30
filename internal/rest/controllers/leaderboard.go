package controllers

import (
	"core/internal/leaderboard"
	"core/internal/rest/handlers"
	"core/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) GetLeaderboard(c handlers.AuthenticatedContext) error {
	type request struct {
		ClubId          uint                        `query:"clubId" validate:"required,gt=0"`
		TopX            int                         `query:"topX" validate:"required,gt=0,lte=50"`
		LeaderboardType leaderboard.LeaderboardType `query:"type" validate:"required,oneof=wins streak rating"`
	}

	type response struct {
		Leaderboard leaderboard.Leaderboard `json:"leaderboard"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	leaderboard, err := h.leaderboardService.GetLeaderboard(ctx, req.ClubId, req.TopX, req.LeaderboardType)
	if err != nil {
		h.logger.Error("failed to get leaderboard",
			"error", err)
		return echo.ErrInternalServerError
	}

	resp := response{
		Leaderboard: *leaderboard,
	}

	return c.JSON(http.StatusOK, resp)
}
