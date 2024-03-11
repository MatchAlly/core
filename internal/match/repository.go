package match

import (
	"context"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const UniqueViolationCode = pq.ErrorCode("23505")

var (
	ErrDuplicateEntry = errors.New("already exists")
)

type Repository interface {
	CreateMatch(ctx context.Context, match *Match) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) CreateMatch(ctx context.Context, match *Match) error {
	result := r.db.WithContext(ctx).
		Create(&match)
	if result.Error != nil {
		pgErr, ok := result.Error.(*pq.Error)
		if ok && pgErr.Code == UniqueViolationCode {
			return ErrDuplicateEntry
		}

		return result.Error
	}

	return nil
}
