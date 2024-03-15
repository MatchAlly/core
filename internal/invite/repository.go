package invite

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetInvitesByUserId(ctx context.Context, userId uint) ([]Invite, error)
	GetInvitesByClubId(ctx context.Context, clubId uint) ([]Invite, error)
	CreateInvite(ctx context.Context, userId, clubId uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetInvitesByUserId(ctx context.Context, userId uint) ([]Invite, error) {
	var invites []Invite
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Find(&invites)
	if result.Error != nil {
		return nil, result.Error
	}

	return invites, nil
}

func (r *repository) GetInvitesByClubId(ctx context.Context, clubId uint) ([]Invite, error) {
	var invites []Invite
	result := r.db.WithContext(ctx).
		Where("club_id = ?", clubId).
		Find(&invites)
	if result.Error != nil {
		return nil, result.Error
	}

	return invites, nil
}

func (r *repository) CreateInvite(ctx context.Context, userId, clubId uint) error {
	invite := Invite{
		UserId: userId,
		ClubId: clubId,
	}

	result := r.db.WithContext(ctx).
		Create(&invite)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
