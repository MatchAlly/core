package game

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetGame(ctx context.Context, id uint) (*Game, error)
	GetGames(ctx context.Context, ids []uint) ([]*Game, error)
	CreateGame(ctx context.Context, game *Game) error
	UpdateGame(ctx context.Context, game *Game) error
	DeleteGame(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetGame(ctx context.Context, id uint) (*Game, error) {
	var game Game
	result := r.db.WithContext(ctx).First(&game, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &game, nil
}

func (r *repository) GetGames(ctx context.Context, ids []uint) ([]*Game, error) {
	var games []*Game
	result := r.db.WithContext(ctx).Find(&games, ids)
	if result.Error != nil {
		return nil, result.Error
	}

	return games, nil
}

func (r *repository) CreateGame(ctx context.Context, game *Game) error {
	result := r.db.WithContext(ctx).Create(game)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) UpdateGame(ctx context.Context, game *Game) error {
	result := r.db.WithContext(ctx).Save(game)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) DeleteGame(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Game{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
