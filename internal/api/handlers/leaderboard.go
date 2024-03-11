package handlers

import (
	"core/internal/api/helpers"
	"core/internal/leaderboard"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetLeaderboard(c helpers.AuthContext) error {
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
		h.l.Error("failed to get leaderboard",
			"error", err)
		return echo.ErrInternalServerError
	}

	resp := response{
		Leaderboard: *leaderboard,
	}

	return c.JSON(http.StatusOK, resp)
}
