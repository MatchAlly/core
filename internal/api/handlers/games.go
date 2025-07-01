package handlers

import (
	"context"
	"core/internal/game"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type getClubGamesRequest struct {
	ClubId uuid.UUID `path:"clubId"`
}

type getClubGamesResponse struct {
	Body struct {
		Games []getClubGamesResponseGame `json:"games"`
	}
}

type getClubGamesResponseGame struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (h *Handler) GetClubGames(ctx context.Context, req *getClubGamesRequest) (*getClubGamesResponse, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authorization.IsMember(ctx, userID, req.ClubId)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to get matches in this club")
	}

	games, err := h.club.GetGames(ctx, req.ClubId)
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
	ClubID uuid.UUID `path:"clubId"`
	Body   struct {
		Name string `json:"name" minLength:"1" maxLength:"50"`
	}
}

type postClubGameResponse struct {
	Body struct {
		GameID uuid.UUID `json:"gameId"`
		Name   string    `json:"name"`
	}
}

func (h *Handler) PostClubGame(ctx context.Context, req *postClubGameRequest) (*postClubGameResponse, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authorization.IsAdmin(ctx, userID, req.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to create games in this club")
	}

	gameID, err := h.game.CreateGame(ctx, req.ClubID, req.Body.Name)
	if err != nil {
		h.l.Error("failed to create game", "error", err)
		return nil, huma.Error500InternalServerError("failed to create game, try again later")
	}

	resp := &postClubGameResponse{}
	resp.Body.GameID = gameID
	resp.Body.Name = req.Body.Name

	return resp, nil
}

type getGameModesRequest struct {
	GameID uuid.UUID `path:"gameId"`
}

type getGameModesResponse struct {
	Body struct {
		Modes []getGameModesResponseMode `json:"modes"`
	}
}

type getGameModesResponseMode struct {
	ID   uuid.UUID `json:"id"`
	Mode string    `json:"mode"`
}

func (h *Handler) GetGameModes(ctx context.Context, req *getGameModesRequest) (*getGameModesResponse, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	game, err := h.game.GetGame(ctx, req.GameID)
	if err != nil {
		h.l.Error("failed to get game", "error", err)
		return nil, huma.Error500InternalServerError("failed to get game")
	}

	ok, err = h.authorization.IsMember(ctx, userID, game.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to view game modes")
	}

	modes, err := h.game.GetGameModes(ctx, req.GameID)
	if err != nil {
		h.l.Error("failed to get game modes", "error", err)
		return nil, huma.Error500InternalServerError("failed to get game modes")
	}

	mappedModes := make([]getGameModesResponseMode, len(modes))
	for i, m := range modes {
		mappedModes[i] = getGameModesResponseMode{
			ID:   m.ID,
			Mode: m.Mode.String(),
		}
	}

	resp := &getGameModesResponse{}
	resp.Body.Modes = mappedModes

	return resp, nil
}

type postGameModeRequest struct {
	GameID uuid.UUID `path:"gameId"`
	Body   struct {
		Mode string `json:"mode" enum:"FREE_FOR_ALL,TEAM,COOP"`
	}
}

func (h *Handler) PostGameMode(ctx context.Context, req *postGameModeRequest) (*struct{}, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	g, err := h.game.GetGame(ctx, req.GameID)
	if err != nil {
		h.l.Error("failed to get game", "error", err)
		return nil, huma.Error500InternalServerError("failed to get game")
	}

	ok, err = h.authorization.IsAdmin(ctx, userID, g.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to add game modes")
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

	if err := h.game.AddGameMode(ctx, req.GameID, mode); err != nil {
		h.l.Error("failed to add game mode", "error", err)
		return nil, huma.Error500InternalServerError("failed to add game mode")
	}

	return nil, nil
}

type deleteGameModeRequest struct {
	GameID uuid.UUID `path:"gameId"`
	Mode   string    `path:"mode" enum:"FREE_FOR_ALL,TEAM,COOP"`
}

func (h *Handler) DeleteGameMode(ctx context.Context, req *deleteGameModeRequest) (*struct{}, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	g, err := h.game.GetGame(ctx, req.GameID)
	if err != nil {
		h.l.Error("failed to get game", "error", err)
		return nil, huma.Error500InternalServerError("failed to get game")
	}

	ok, err = h.authorization.IsAdmin(ctx, userID, g.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to remove game modes")
	}

	var mode game.Mode
	switch req.Mode {
	case "FREE_FOR_ALL":
		mode = game.ModeFreeForAll
	case "TEAM":
		mode = game.ModeTeam
	case "COOP":
		mode = game.ModeCoop
	default:
		return nil, huma.Error400BadRequest("invalid game mode")
	}

	if err := h.game.RemoveGameMode(ctx, req.GameID, mode); err != nil {
		h.l.Error("failed to remove game mode", "error", err)
		return nil, huma.Error500InternalServerError("failed to remove game mode")
	}

	return nil, nil
}

type deleteGameRequest struct {
	GameID uuid.UUID `path:"gameId"`
}

func (h *Handler) DeleteGame(ctx context.Context, req *deleteGameRequest) (*struct{}, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	g, err := h.game.GetGame(ctx, req.GameID)
	if err != nil {
		h.l.Error("failed to get game", "error", err)
		return nil, huma.Error500InternalServerError("failed to get game")
	}

	ok, err = h.authorization.IsAdmin(ctx, userID, g.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to delete games in this club")
	}

	if err := h.game.DeleteGame(ctx, req.GameID); err != nil {
		h.l.Error("failed to delete game", "error", err)
		return nil, huma.Error500InternalServerError("failed to delete game")
	}

	return nil, nil
}

type putGameRequest struct {
	GameID uuid.UUID `path:"gameId"`
	Body   struct {
		Name string `json:"name" minLength:"1" maxLength:"50"`
	}
}

type putGameResponse struct {
	Body struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
}

func (h *Handler) PutGame(ctx context.Context, req *putGameRequest) (*putGameResponse, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	g, err := h.game.GetGame(ctx, req.GameID)
	if err != nil {
		h.l.Error("failed to get game", "error", err)
		return nil, huma.Error500InternalServerError("failed to get game")
	}

	ok, err = h.authorization.IsAdmin(ctx, userID, g.ClubID)
	if err != nil {
		h.l.Error("failed to check authorization", "error", err)
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to update games in this club")
	}

	g.Name = req.Body.Name
	if err := h.game.UpdateGame(ctx, g); err != nil {
		h.l.Error("failed to update game", "error", err)
		return nil, huma.Error500InternalServerError("failed to update game")
	}

	resp := &putGameResponse{}
	resp.Body.ID = g.ID
	resp.Body.Name = g.Name

	return resp, nil
}
