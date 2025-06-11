package member

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetMembersInClub(ctx context.Context, clubId int) ([]Member, error)
	GetUserMemberships(ctx context.Context, userId int) ([]Member, error)
	GetMember(ctx context.Context, id int) (*Member, error)
	CreateMember(ctx context.Context, member *Member) error
	UpdateRole(ctx context.Context, memberId int, role Role) error
	DeleteMember(ctx context.Context, memberId int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetMembersInClub(ctx context.Context, clubId int) ([]Member, error) {
	var members []Member
	err := r.db.SelectContext(ctx, &members, "SELECT * FROM members WHERE club_id = $1", clubId)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *repository) GetUserMemberships(ctx context.Context, userId int) ([]Member, error) {
	var memberships []Member
	err := r.db.SelectContext(ctx, &memberships, "SELECT * FROM members WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}

	return memberships, nil
}

func (r *repository) GetMember(ctx context.Context, id int) (*Member, error) {
	var member Member
	err := r.db.GetContext(ctx, &member, "SELECT * FROM club_members WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *repository) CreateMember(ctx context.Context, member *Member) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO club_members (club_id, user_id, role) VALUES ($1, $2, $3)",
		member.ClubID, member.UserID, member.Role)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateRole(ctx context.Context, memberId int, role Role) error {
	_, err := r.db.ExecContext(ctx, "UPDATE members SET role = $1 WHERE id = $2", role, memberId)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteMember(ctx context.Context, memberId int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM members WHERE id = $1", memberId)
	if err != nil {
		return err
	}

	return nil
}
