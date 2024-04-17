package handlers

import (
	"core/internal/api/helpers"
	"core/internal/match"
	"core/internal/statistic"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type postMatchRequest struct {
	GameId  uint   `json:"gameId" validate:"required"`
	Rated   bool   `json:"rated" validate:"required"`
	TeamA   []uint `json:"teamA" validate:"required"`
	TeamB   []uint `json:"teamB" validate:"required"`
	ScoresA []int  `json:"scoresA" validate:"required"`
	ScoresB []int  `json:"scoresB" validate:"required"`
}

func (h *Handler) PostMatch(c helpers.AuthContext) error {

	ctx := c.Request().Context()

	req, err := helpers.Bind[postMatchRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if len(req.ScoresA) != len(req.ScoresB) {
		return echo.ErrBadRequest
	}

	result, winners, losers := h.matchService.DetermineResult(ctx, req.TeamA, req.TeamB, req.ScoresA, req.ScoresB)

	if err = h.matchService.CreateMatch(ctx, req.TeamA, req.TeamB, req.ScoresA, req.ScoresB, result); err != nil {
		h.l.Error("failed to create match", zap.Error(err))
		return echo.ErrInternalServerError
	}

	if result == match.Draw {
		allPlayers := append(req.TeamA, req.TeamB...)
		if err := h.statisticService.UpdateGameStatisticsByMemberIds(ctx, allPlayers, req.GameId, statistic.ResultDraw); err != nil {
			h.l.Error("failed to update statistics for draw", zap.Error(err))
			return echo.ErrInternalServerError
		}
	} else {
		if err := h.statisticService.UpdateGameStatisticsByMemberIds(ctx, winners, req.GameId, statistic.ResultWin); err != nil {
			h.l.Error("failed to update statistics for winners", zap.Error(err))
			return echo.ErrInternalServerError
		}

		if err := h.statisticService.UpdateGameStatisticsByMemberIds(ctx, losers, req.GameId, statistic.ResultLoss); err != nil {
			h.l.Error("failed to update statistics for losers", zap.Error(err))
			return echo.ErrInternalServerError
		}
	}

	if !req.Rated {
		return c.NoContent(http.StatusCreated)
	}

	isDraw := result == match.Draw
	if err := h.ratingService.UpdateRatings(ctx, isDraw, winners, losers); err != nil {
		h.l.Error("failed to update ratings", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}
