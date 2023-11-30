package statistic

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetStatisticsByUserIds(ctx context.Context, userIds []uint) ([]*Statistic, error)
	GetStatisticByUserId(ctx context.Context, userId uint) (*Statistic, error)
	GetTopXAmongUserIdsByWins(ctx context.Context, topX int, userIds []uint) (topXUserIds []uint, values []int, err error)
	GetTopXAmongUserIdsByStreak(ctx context.Context, topX int, userIds []uint) (topXUserIds []uint, values []int, err error)
	CreateStatistic(ctx context.Context, userId uint) error
	UpdateStatistics(ctx context.Context, stats []Statistic) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetStatisticsByUserIds(ctx context.Context, userIds []uint) ([]*Statistic, error) {
	var stats []*Statistic
	result := r.db.WithContext(ctx).
		Where("user_id IN ?", userIds).
		Find(&stats)
	if result.Error != nil {
		return nil, result.Error
	}

	return stats, nil
}

func (r *RepositoryImpl) GetStatisticByUserId(ctx context.Context, userId uint) (*Statistic, error) {
	var stats Statistic
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		First(&stats)
	if result.Error != nil {
		return nil, result.Error
	}

	return &stats, nil
}

func (r *RepositoryImpl) GetTopXAmongUserIdsByWins(ctx context.Context, topX int, userIds []uint) ([]uint, []int, error) {
	var topXUserIds []uint
	var wins []int

	result := r.db.
		WithContext(ctx).
		Model(&Statistic{}).
		Order("wins desc").
		Limit(topX).
		Pluck("user_id", &topXUserIds).
		Pluck("wins", &wins).
		Where("user_id IN ?", userIds)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return topXUserIds, wins, nil
}

func (r *RepositoryImpl) GetTopXAmongUserIdsByStreak(ctx context.Context, topX int, userIds []uint) ([]uint, []int, error) {
	var topXUserIds []uint
	var streaks []int

	result := r.db.WithContext(ctx).
		Model(&Statistic{}).
		Order("streak desc").
		Limit(topX).
		Where("user_id IN ?", userIds).
		Pluck("user_id", &topXUserIds).
		Pluck("streak", &streaks)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return topXUserIds, streaks, nil
}

func (r *RepositoryImpl) CreateStatistic(ctx context.Context, userId uint) error {
	stat := Statistic{
		UserId: userId,
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

func (r *RepositoryImpl) UpdateStatistics(ctx context.Context, stats []Statistic) error {
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
