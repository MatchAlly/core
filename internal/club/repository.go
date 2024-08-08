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
	GetClubIDsWithUserID(ctx context.Context, userId uint) ([]uint, error)
	GetMembers(ctx context.Context, id uint) ([]Member, error)
	CreateClub(ctx context.Context, Club *Club) (clubId uint, err error)
	AddUserToClub(ctx context.Context, userId uint, clubId uint, role Role) error
	DeleteMember(ctx context.Context, memberId uint) error
	DeleteClub(ctx context.Context, id uint) error
	UpdateClub(ctx context.Context, id uint, name string) error
	UpdateMemberRole(ctx context.Context, memberId uint, role Role) error
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
	var c Club
	result := r.db.WithContext(ctx).
		First(&c, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return &c, nil
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

func (r *repository) GetClubIDsWithUserID(ctx context.Context, userId uint) ([]uint, error) {
	var clubIds []uint
	result := r.db.WithContext(ctx).
		Model(&Member{}).
		Where("user_id = ?", userId).
		Pluck("club_id", &clubIds)
	if result.Error != nil {
		return nil, result.Error
	}

	return clubIds, nil
}

func (r *repository) GetMembers(ctx context.Context, id uint) ([]Member, error) {
	var members []Member
	result := r.db.WithContext(ctx).
		Model(&Member{}).
		Where("club_id = ?", id).
		Find(&members)

	if result.Error != nil {
		return nil, result.Error
	}

	return members, nil
}

func (r *repository) CreateClub(ctx context.Context, c *Club) (uint, error) {
	result := r.db.WithContext(ctx).
		Create(&c)
	if result.Error != nil {
		pgErr, ok := result.Error.(*pq.Error)
		if ok && pgErr.Code == UniqueViolationCode {
			return 0, ErrDuplicateEntry
		}

		return 0, result.Error
	}

	return c.ID, nil
}

func (r *repository) AddUserToClub(ctx context.Context, userId uint, clubId uint, role Role) error {
	m := Member{
		UserId: userId,
		ClubId: clubId,
		Role:   role,
	}

	result := r.db.WithContext(ctx).
		Create(&m)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) DeleteMember(ctx context.Context, memberId uint) error {
	result := r.db.WithContext(ctx).
		Delete(&Member{}, memberId)
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

func (r *repository) UpdateMemberRole(ctx context.Context, memberId uint, role Role) error {
	result := r.db.WithContext(ctx).
		Model(&Member{}).
		Where("id = ?", memberId).
		Update("role", role)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
