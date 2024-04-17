package rating

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetRatingByMemberId(ctx context.Context, memberId uint) (*Rating, error)
	GetRatingsByMemberIds(ctx context.Context, memberIds []uint) ([]Rating, error)
	GetTopMembersByRating(ctx context.Context, topX int, memberIds []uint) (topXMemberIds []uint, ratings []int, err error)
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

func (r *repository) GetRatingByMemberId(ctx context.Context, memberId uint) (*Rating, error) {
	var rating Rating

	result := r.db.WithContext(ctx).
		Where("member_id = ?", memberId).
		First(&rating)
	if result.Error != nil {
		return nil, result.Error
	}

	return &rating, nil
}

func (r *repository) GetRatingsByMemberIds(ctx context.Context, memberIds []uint) ([]Rating, error) {
	var ratings []Rating
	result := r.db.WithContext(ctx).
		Where("member_id IN ?", memberIds).
		Find(&ratings)
	if result.Error != nil {
		return nil, result.Error
	}

	return ratings, nil
}

func (r *repository) GetTopMembersByRating(ctx context.Context, topX int, memberIds []uint) ([]uint, []int, error) {
	var topMemberIds []uint
	var ratings []int

	result := r.db.WithContext(ctx).
		Model(&Rating{}).
		Order("rating desc").
		Limit(topX).
		Pluck("member_id", &topMemberIds).
		Pluck("value", &ratings).
		Where("member_id IN ?", memberIds)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return topMemberIds, ratings, nil
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
