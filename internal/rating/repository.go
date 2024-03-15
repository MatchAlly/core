package rating

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetRatingByUserId(ctx context.Context, userId uint) (*Rating, error)
	GetRatingsByUserIds(ctx context.Context, userIds []uint) ([]Rating, error)
	GetTopXAmongUserIdsByRating(ctx context.Context, topX int, userIds []uint) (topXUserIds []uint, ratings []int, err error)
	CreateRating(ctx context.Context, rating *Rating) error
	UpdateRating(ctx context.Context, ratings *Rating) error
	UpdateRatings(ctx context.Context, ratings []Rating) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetRatingByUserId(ctx context.Context, userId uint) (*Rating, error) {
	var rating Rating

	result := r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		First(&rating)
	if result.Error != nil {
		return nil, result.Error
	}

	return &rating, nil
}

func (r *repository) GetRatingsByUserIds(ctx context.Context, userIds []uint) ([]Rating, error) {
	var ratings []Rating
	result := r.db.WithContext(ctx).
		Where("user_id IN ?", userIds).
		Find(&ratings)
	if result.Error != nil {
		return nil, result.Error
	}

	return ratings, nil
}

func (r *repository) GetTopXAmongUserIdsByRating(ctx context.Context, topX int, userIds []uint) ([]uint, []int, error) {
	var topXUserIds []uint
	var ratings []int

	result := r.db.WithContext(ctx).
		Model(&Rating{}).
		Order("rating desc").
		Limit(topX).
		Pluck("user_id", &topXUserIds).
		Pluck("value", &ratings).
		Where("user_id IN ?", userIds)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return topXUserIds, ratings, nil
}

func (r *repository) CreateRating(ctx context.Context, rating *Rating) error {
	result := r.db.WithContext(ctx).
		Create(rating)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) UpdateRatings(ctx context.Context, ratings []Rating) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, rating := range ratings {
			result := tx.WithContext(ctx).
				Model(&rating).
				Updates(rating)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func (r *repository) UpdateRating(ctx context.Context, rating *Rating) error {
	result := r.db.WithContext(ctx).
		Model(&rating).
		Updates(rating)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
