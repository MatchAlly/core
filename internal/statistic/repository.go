package statistic

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetStatisticsByMemberId(ctx context.Context, memberId uint) ([]Statistic, error)
	GetGameStatisticsForMemberIds(ctx context.Context, memberIds []uint, gameId uint) ([]Statistic, error)
	CreateDefaultStatistic(ctx context.Context, memberId, gameId uint) error
	UpdateStatistics(ctx context.Context, statistics []Statistic) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetStatisticsByMemberId(ctx context.Context, memberId uint) ([]Statistic, error) {
	var stats []Statistic
	result := r.db.WithContext(ctx).
		Where("member_id = ?", memberId).
		Find(&stats)
	if result.Error != nil {
		return nil, result.Error
	}

	return stats, nil
}

func (r *repository) GetGameStatisticsForMemberIds(ctx context.Context, memberIds []uint, gameId uint) ([]Statistic, error) {
	var stats []Statistic
	result := r.db.WithContext(ctx).
		Where("member_id IN ?", memberIds).
		Where("game_id = ?", gameId).
		Find(&stats)
	if result.Error != nil {
		return nil, result.Error
	}

	return stats, nil
}

func (r *repository) CreateDefaultStatistic(ctx context.Context, memberId, gameId uint) error {
	stat := Statistic{
		MemberId: memberId,
		GameId:   gameId,
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
