package game

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	// Games
	GetGame(ctx context.Context, id uuid.UUID) (*Game, error)
	GetGames(ctx context.Context, ids []uuid.UUID) ([]Game, error)
	CreateGame(ctx context.Context, clubID uuid.UUID, name string) (uuid.UUID, error)
	UpdateGame(ctx context.Context, game *Game) error
	DeleteGame(ctx context.Context, id uuid.UUID) error

	// Game modes
	GetGameModes(ctx context.Context, gameID uuid.UUID) ([]Gamemode, error)
	AddGameMode(ctx context.Context, gameID uuid.UUID, mode Mode) error
	RemoveGameMode(ctx context.Context, gameID uuid.UUID, mode Mode) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) GetGame(ctx context.Context, id uuid.UUID) (*Game, error) {
	return s.repo.GetGame(ctx, id)
}

func (s *service) GetGames(ctx context.Context, ids []uuid.UUID) ([]Game, error) {
	return s.repo.GetGames(ctx, ids)
}

func (s *service) CreateGame(ctx context.Context, clubID uuid.UUID, name string) (uuid.UUID, error) {
	// Validate game name
	if len(name) < 1 || len(name) > 50 {
		return uuid.Nil, fmt.Errorf("game name must be between 1 and 50 characters")
	}

	// Check for duplicate name
	unique, err := s.repo.IsGameNameUnique(ctx, clubID, name, uuid.Nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to check game name uniqueness: %w", err)
	}
	if !unique {
		return uuid.Nil, fmt.Errorf("game name already exists in this club")
	}

	game := &Game{
		ClubID: clubID,
		Name:   name,
	}

	return s.repo.CreateGame(ctx, game)
}

func (s *service) UpdateGame(ctx context.Context, game *Game) error {
	// Validate game name
	if len(game.Name) < 1 || len(game.Name) > 50 {
		return fmt.Errorf("game name must be between 1 and 50 characters")
	}

	// Check for duplicate name
	unique, err := s.repo.IsGameNameUnique(ctx, game.ClubID, game.Name, game.ID)
	if err != nil {
		return fmt.Errorf("failed to check game name uniqueness: %w", err)
	}
	if !unique {
		return fmt.Errorf("game name already exists in this club")
	}

	return s.repo.UpdateGame(ctx, game)
}

func (s *service) DeleteGame(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteGame(ctx, id)
}

func (s *service) GetGameModes(ctx context.Context, gameID uuid.UUID) ([]Gamemode, error) {
	return s.repo.GetGameModes(ctx, gameID)
}

func (s *service) AddGameMode(ctx context.Context, gameID uuid.UUID, mode Mode) error {
	if mode == ModeNone {
		return fmt.Errorf("invalid game mode")
	}
	return s.repo.AddGameMode(ctx, gameID, mode)
}

func (s *service) RemoveGameMode(ctx context.Context, gameID uuid.UUID, mode Mode) error {
	if mode == ModeNone {
		return fmt.Errorf("invalid game mode")
	}
	return s.repo.RemoveGameMode(ctx, gameID, mode)
}
