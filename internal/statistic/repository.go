package statistic

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetStatisticsByUserIds(ctx context.Context, userIds []uint) ([]*Statistic, error)
	GetStatisticsByUserId(ctx context.Context, userId uint) ([]Statistic, error)
	CreateStatistic(ctx context.Context, userId, gameId uint) error
	UpdateStatistics(ctx context.Context, stats []Statistic) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetStatisticsByUserIds(ctx context.Context, userIds []uint) ([]*Statistic, error) {
	var stats []*Statistic
	result := r.db.WithContext(ctx).
		Where("user_id IN ?", userIds).
		Find(&stats)
	if result.Error != nil {
		return nil, result.Error
	}

	return stats, nil
}

func (r *repository) GetStatisticsByUserId(ctx context.Context, userId uint) ([]Statistic, error) {
	var stats []Statistic
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Find(&stats)
	if result.Error != nil {
		return nil, result.Error
	}

	return stats, nil
}

func (r *repository) CreateStatistic(ctx context.Context, userId, gameId uint) error {
	stat := Statistic{
		UserId: userId,
		GameId: gameId,
		Wins:   0,
		Draws:  0,
		Losses: 0,
		Streak: 0,
	}

	result := r.db.WithContext(ctx).
		Create(&stat)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) UpdateStatistics(ctx context.Context, stats []Statistic) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			result := tx.WithContext(ctx).
				Model(&stat).
				Updates(stat)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}
