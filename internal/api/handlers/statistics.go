package handlers

import (
	"context"
	"core/internal/statistic"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type getMemberStatisticsRequest struct {
	MemberID uuid.UUID  `path:"memberId"`
	GameID   *uuid.UUID `query:"gameId" required:"false"`
}

type getMemberStatisticsResponse struct {
	Body struct {
		Statistics []getMemberStatisticsResponseStatistic `json:"statistics"`
	}
}

type getMemberStatisticsResponseStatistic struct {
	GameID uuid.UUID `json:"gameId"`
	Wins   int       `json:"wins"`
	Losses int       `json:"losses"`
	Draws  int       `json:"draws"`
	Streak int       `json:"streak"`
}

func (h *Handler) GetMemberStatistics(ctx context.Context, req *getMemberStatisticsRequest) (*getMemberStatisticsResponse, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	// Get the member's club ID
	member, err := h.member.GetMember(ctx, req.MemberID)
	if err != nil {
		h.l.Error("failed to get member", "error", err)
		return nil, huma.Error500InternalServerError("failed to get member")
	}

	// Check if the user is authorized to view the member's statistics
	ok, err = h.authorization.IsMember(ctx, userID, member.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to view statistics in this club")
	}

	var stats []statistic.Statistic
	if req.GameID != nil {
		// Get statistics for a specific game
		stat, err := h.statistic.GetStatistics(ctx, req.MemberID, *req.GameID)
		if err != nil {
			h.l.Error("failed to get statistics", "error", err)
			return nil, huma.Error500InternalServerError("failed to get statistics")
		}
		stats = []statistic.Statistic{*stat}
	} else {
		// Get statistics for all games
		stats, err = h.statistic.GetStatisticsByGame(ctx, req.MemberID)
		if err != nil {
			h.l.Error("failed to get statistics", "error", err)
			return nil, huma.Error500InternalServerError("failed to get statistics")
		}
	}

	mappedStats := make([]getMemberStatisticsResponseStatistic, len(stats))
	for i, s := range stats {
		mappedStats[i] = getMemberStatisticsResponseStatistic{
			GameID: s.GameId,
			Wins:   s.Wins,
			Losses: s.Losses,
			Draws:  s.Draws,
			Streak: s.Streak,
		}
	}

	resp := &getMemberStatisticsResponse{}
	resp.Body.Statistics = mappedStats

	return resp, nil
}

type getGameRankingsRequest struct {
	GameID uuid.UUID `path:"gameId"`
}

type getGameRankingsResponse struct {
	Body struct {
		Rankings []getGameRankingsResponseRanking `json:"rankings"`
	}
}

type getGameRankingsResponseRanking struct {
	MemberID uuid.UUID `json:"memberId"`
	Wins     int       `json:"wins"`
	Losses   int       `json:"losses"`
	Draws    int       `json:"draws"`
	Streak   int       `json:"streak"`
}

func (h *Handler) GetGameRankings(ctx context.Context, req *getGameRankingsRequest) (*getGameRankingsResponse, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	// Get the game's club ID
	game, err := h.game.GetGame(ctx, req.GameID)
	if err != nil {
		h.l.Error("failed to get game", "error", err)
		return nil, huma.Error500InternalServerError("failed to get game")
	}

	// Check if the user is authorized to view the rankings
	ok, err = h.authorization.IsMember(ctx, userID, game.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to view rankings in this club")
	}

	// Get statistics for all members in the game
	stats, err := h.statistic.GetStatisticsByGame(ctx, req.GameID)
	if err != nil {
		h.l.Error("failed to get statistics", "error", err)
		return nil, huma.Error500InternalServerError("failed to get statistics")
	}

	mappedRankings := make([]getGameRankingsResponseRanking, len(stats))
	for i, s := range stats {
		mappedRankings[i] = getGameRankingsResponseRanking{
			MemberID: s.MemberId,
			Wins:     s.Wins,
			Losses:   s.Losses,
			Draws:    s.Draws,
			Streak:   s.Streak,
		}
	}

	resp := &getGameRankingsResponse{}
	resp.Body.Rankings = mappedRankings

	return resp, nil
}
