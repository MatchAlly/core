package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type getClubGamesRequest struct {
	ClubId int `path:"clubId" minimum:"1"`
}

type getClubGamesResponse struct {
	Games []getClubGamesResponseGame `json:"games"`
}

type getClubGamesResponseGame struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetClubGames(ctx context.Context, req *getClubGamesRequest) (*getClubGamesResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authZService.IsMember(ctx, userID, req.ClubId)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to get matches in this club")
	}

	games, err := h.clubService.GetGames(ctx, req.ClubId)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to get games, try again later")
	}

	mappedGames := make([]getClubGamesResponseGame, len(games))
	for i, g := range games {
		mappedGames[i] = getClubGamesResponseGame{
			ID:   g.ID,
			Name: g.Name,
		}
	}

	resp := &getClubGamesResponse{
		Games: mappedGames,
	}

	return resp, nil
}

type postClubGameRequest struct {
	ClubID int    `path:"clubId" minimum:"1"`
	Name   string `json:"name" minLength:"1" maxLength:"50"`
}

type postClubGameResponse struct {
	GameID int    `json:"gameId"`
	Name   string `json:"name"`
}

func (h *Handler) PostClubGame(ctx context.Context, req *postClubGameRequest) (*postClubGameResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authZService.IsAdmin(ctx, userID, req.ClubID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to create games in this club")
	}

	gameID, err := h.gameService.CreateGame(ctx, req.ClubID, req.Name)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to create game, try again later")
	}

	resp := &postClubGameResponse{
		GameID: gameID,
		Name:   req.Name,
	}

	return resp, nil
}
