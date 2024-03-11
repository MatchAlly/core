package club

import (
	"context"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const UniqueViolationCode = pq.ErrorCode("23505")

var (
	ErrDuplicateEntry = errors.New("already exists")
	ErrNotFound       = errors.New("not found")
)

type Repository interface {
	GetClub(ctx context.Context, id uint) (*Club, error)
	GetClubs(ctx context.Context, ids []uint) ([]Club, error)
	GetUserIdsInClub(ctx context.Context, id uint) ([]uint, error)
	GetInvitesByUserId(ctx context.Context, userId uint) ([]Invite, error)
	InviteToClub(ctx context.Context, userIds []uint, clubId uint) error
	CreateClub(ctx context.Context, Club *Club) (clubId uint, err error)
	AddUserToClub(ctx context.Context, userId uint, clubId uint, role Role) error
	RemoveUserFromClub(ctx context.Context, userId uint, clubId uint) error
	DeleteClub(ctx context.Context, id uint) error
	UpdateClub(ctx context.Context, id uint, name string) error
	UpdateUserRole(ctx context.Context, userId uint, clubId uint, role Role) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetClub(ctx context.Context, id uint) (*Club, error) {
	var club Club
	result := r.db.WithContext(ctx).
		First(&club, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return &club, nil
}

func (r *repository) GetClubs(ctx context.Context, ids []uint) ([]Club, error) {
	var clubs []Club
	result := r.db.WithContext(ctx).
		Find(&clubs, ids)
	if result.Error != nil {
		return nil, result.Error
	}

	return clubs, nil
}

func (r *repository) GetUserIdsInClub(ctx context.Context, id uint) ([]uint, error) {
	var invites []Invite
	result := r.db.WithContext(ctx).
		Where("Club_id = ?", id).
		Find(&invites)
	if result.Error != nil {
		return nil, result.Error
	}

	var userIds []uint
	for _, invite := range invites {
		userIds = append(userIds, invite.UserId)
	}

	return userIds, nil
}

func (r *repository) GetInvitesByUserId(ctx context.Context, userId uint) ([]Invite, error) {
	var invites []Invite
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND accepted = ?", userId, false).
		Find(&invites)
	if result.Error != nil {
		return nil, result.Error
	}

	return invites, nil
}

func (r *repository) CreateClub(ctx context.Context, Club *Club) (uint, error) {
	result := r.db.WithContext(ctx).
		Create(&Club)
	if result.Error != nil {
		pgErr, ok := result.Error.(*pq.Error)
		if ok && pgErr.Code == UniqueViolationCode {
			return 0, ErrDuplicateEntry
		}

		return 0, result.Error
	}

	return Club.Id, nil
}

func (r *repository) AddUserToClub(ctx context.Context, userId uint, clubId uint, role Role) error {
	invite := &Invite{
		ClubId:   clubId,
		UserId:   userId,
		Accepted: true,
		Role:     role,
	}

	result := r.db.WithContext(ctx).
		Create(&invite)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) RemoveUserFromClub(ctx context.Context, userId uint, clubId uint) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND Club_id = ?", userId, clubId).
		Delete(&Invite{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) DeleteClub(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Delete(&Club{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) UpdateClub(ctx context.Context, id uint, name string) error {
	result := r.db.WithContext(ctx).
		Model(&Club{}).
		Where("id = ?", id).
		Update("name", name)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) UpdateUserRole(ctx context.Context, userId uint, clubId uint, role Role) error {
	result := r.db.WithContext(ctx).
		Model(&Invite{}).
		Where("user_id = ? AND club_id = ?", userId, clubId).
		Update("role", role)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) InviteToClub(ctx context.Context, userIds []uint, clubId uint) error {
	var invites []Invite
	for _, userId := range userIds {
		invites = append(invites, Invite{
			ClubId:   clubId,
			UserId:   userId,
			Accepted: false,
			Role:     MemberRole,
		})
	}

	result := r.db.WithContext(ctx).
		Create(&invites)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
