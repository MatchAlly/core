package user

import (
	"context"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const ErrCodeMySQLDuplicateEntry uint16 = 1062

var (
	ErrDuplicateEntry = errors.New("already exists")
	ErrNotFound       = errors.New("not found")
)

type Repository interface {
	GetUser(ctx context.Context, id uint) (*User, error)
	GetUsers(ctx context.Context, ids []uint) ([]*User, error)
	GetUsersInClub(ctx context.Context, clubId uint) ([]User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUsersByEmails(ctx context.Context, emails []string) ([]*User, error)
	CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint) error
	UpdateUser(ctx context.Context, user *User) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetUser(ctx context.Context, id uint) (*User, error) {
	var user *User
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *RepositoryImpl) GetUsers(ctx context.Context, ids []uint) ([]*User, error) {
	var users []*User
	result := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (r *RepositoryImpl) GetUsersInClub(ctx context.Context, clubId uint) ([]User, error) {
	var users []User
	result := r.db.WithContext(ctx).
		Where("club_id = ?", clubId).
		Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (r *RepositoryImpl) CreateUser(ctx context.Context, user *User) error {
	result := r.db.WithContext(ctx).
		Create(&user)
	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return ErrDuplicateEntry
		}

		return result.Error
	}

	return nil
}

func (r *RepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return &user, nil
}

func (r *RepositoryImpl) GetUsersByEmails(ctx context.Context, emails []string) ([]*User, error) {
	var users []*User
	result := r.db.WithContext(ctx).
		Where("email IN ?", emails).
		Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (r *RepositoryImpl) DeleteUser(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&User{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *RepositoryImpl) UpdateUser(ctx context.Context, user *User) error {
	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", user.Id).
		Updates(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
