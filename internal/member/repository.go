package member

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetMembersInClub(ctx context.Context, clubId uint) ([]Member, error)
	GetUserMemberships(ctx context.Context, userId uint) ([]Member, error)
	UpdateRole(ctx context.Context, memberId uint, role Role) error
	DeleteMembership(ctx context.Context, memberId uint) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetMembersInClub(ctx context.Context, clubId uint) ([]Member, error) {
	var members []Member
	err := r.db.SelectContext(ctx, &members, "SELECT * FROM members WHERE club_id = $1", clubId)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *repository) GetUserMemberships(ctx context.Context, userId uint) ([]Member, error) {
	var memberships []Member
	err := r.db.SelectContext(ctx, &memberships, "SELECT * FROM members WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}

	return memberships, nil
}

func (r *repository) UpdateRole(ctx context.Context, memberId uint, role Role) error {
	_, err := r.db.ExecContext(ctx, "UPDATE members SET role = $1 WHERE id = $2", role, memberId)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteMembership(ctx context.Context, memberId uint) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM members WHERE id = $1", memberId)
	if err != nil {
		return err
	}

	return nil
}
