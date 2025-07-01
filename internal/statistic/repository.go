package statistic

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetStatistics(ctx context.Context, memberID, gameID uuid.UUID) (*Statistic, error)
	GetStatisticsByGame(ctx context.Context, gameID uuid.UUID) ([]Statistic, error)
	CreateStatistics(ctx context.Context, stats *Statistic) (uuid.UUID, error)
	UpdateStatistics(ctx context.Context, stats *Statistic) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetStatistics(ctx context.Context, memberID, gameID uuid.UUID) (*Statistic, error) {
	var stats Statistic
	err := r.db.GetContext(ctx, &stats,
		"SELECT * FROM statistics WHERE member_id = $1 AND game_id = $2",
		memberID, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}
	return &stats, nil
}

func (r *repository) GetStatisticsByGame(ctx context.Context, gameID uuid.UUID) ([]Statistic, error) {
	var stats []Statistic
	err := r.db.SelectContext(ctx, &stats,
		"SELECT * FROM statistics WHERE game_id = $1 ORDER BY wins DESC, losses ASC",
		gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics by game: %w", err)
	}
	return stats, nil
}

func (r *repository) CreateStatistics(ctx context.Context, stats *Statistic) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO statistics (member_id, game_id, wins, losses, draws, streak) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		stats.MemberId, stats.GameId, stats.Wins, stats.Losses, stats.Draws, stats.Streak).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create statistics: %w", err)
	}
	return id, nil
}

func (r *repository) UpdateStatistics(ctx context.Context, stats *Statistic) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE statistics SET wins = $1, losses = $2, draws = $3, streak = $4, updated_at = CURRENT_TIMESTAMP WHERE member_id = $5 AND game_id = $6",
		stats.Wins, stats.Losses, stats.Draws, stats.Streak, stats.MemberId, stats.GameId)
	if err != nil {
		return fmt.Errorf("failed to update statistics: %w", err)
	}
	return nil
}
