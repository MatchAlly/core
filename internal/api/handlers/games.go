package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type getClubGamesRequest struct {
	ClubId int `path:"clubId" minimum:"1"`
}

type getClubGamesResponse struct {
	Body struct {
		Games []getClubGamesResponseGame `json:"games"`
	}
}

type getClubGamesResponseGame struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetClubGames(ctx context.Context, req *getClubGamesRequest) (*getClubGamesResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authZService.IsMember(ctx, userID, req.ClubId)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to get matches in this club")
	}

	games, err := h.clubService.GetGames(ctx, req.ClubId)
	if err != nil {
		h.l.Error("failed to get games", "error", err)
		return nil, huma.Error500InternalServerError("failed to get games, try again later")
	}

	mappedGames := make([]getClubGamesResponseGame, len(games))
	for i, g := range games {
		mappedGames[i] = getClubGamesResponseGame{
			ID:   g.ID,
			Name: g.Name,
		}
	}

	resp := &getClubGamesResponse{}
	resp.Body.Games = mappedGames

	return resp, nil
}

type postClubGameRequest struct {
	ClubID int `path:"clubId" minimum:"1"`
	Body   struct {
		Name string `json:"name" minLength:"1" maxLength:"50"`
	}
}

type postClubGameResponse struct {
	Body struct {
		GameID int    `json:"gameId"`
		Name   string `json:"name"`
	}
}

func (h *Handler) PostClubGame(ctx context.Context, req *postClubGameRequest) (*postClubGameResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authZService.IsAdmin(ctx, userID, req.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to create games in this club")
	}

	gameID, err := h.gameService.CreateGame(ctx, req.ClubID, req.Body.Name)
	if err != nil {
		h.l.Error("failed to create game", "error", err)
		return nil, huma.Error500InternalServerError("failed to create game, try again later")
	}

	resp := &postClubGameResponse{}
	resp.Body.GameID = gameID
	resp.Body.Name = req.Body.Name

	return resp, nil
}
