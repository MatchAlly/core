package club

import (
	"context"
	"core/internal/game"
	"core/internal/member"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrDuplicateEntry = fmt.Errorf("duplicate entry")
	ErrNotFound       = fmt.Errorf("not found")
)

type Repository interface {
	GetClub(ctx context.Context, id int) (*Club, error)
	GetClubs(ctx context.Context, ids []int) ([]Club, error)
	CreateClub(ctx context.Context, Club *Club) (clubId int, err error)
	DeleteClub(ctx context.Context, id int) error
	UpdateClub(ctx context.Context, id int, name string) error
	GetGames(ctx context.Context, clubID int) ([]game.Game, error)
	CreateMember(ctx context.Context, member *member.Member) error
	IsMember(ctx context.Context, userId, clubId int) (bool, error)
	CreateInvite(ctx context.Context, invite *Invite) error
	GetPendingInvites(ctx context.Context, clubId int) ([]Invite, error)
	GetUserInvites(ctx context.Context, userId int) ([]Invite, error)
	GetInvite(ctx context.Context, id int) (*Invite, error)
	DeleteInvite(ctx context.Context, id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetClub(ctx context.Context, id int) (*Club, error) {
	var c *Club

	err := r.db.GetContext(ctx, c, "SELECT * FROM clubs WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return c, nil
}

func (r *repository) GetClubs(ctx context.Context, ids []int) ([]Club, error) {
	var clubs []Club

	query, args, err := sqlx.In("SELECT * FROM clubs WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &clubs, query, args...)
	if err != nil {
		return nil, err
	}

	return clubs, nil
}

func (r *repository) CreateClub(ctx context.Context, c *Club) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO clubs (name) VALUES ($1) RETURNING id",
		c.Name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) DeleteClub(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM clubs WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateClub(ctx context.Context, id int, name string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE clubs SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetGames(ctx context.Context, clubID int) ([]game.Game, error) {
	var games []game.Game

	err := r.db.SelectContext(ctx, &games, "SELECT * FROM games WHERE club_id = $1", clubID)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (r *repository) CreateMember(ctx context.Context, member *member.Member) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO club_members (club_id, user_id, role) VALUES ($1, $2, $3)",
		member.ClubID, member.UserID, member.Role)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) IsMember(ctx context.Context, userId, clubId int) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM club_members WHERE user_id = $1 AND club_id = $2",
		userId, clubId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) CreateInvite(ctx context.Context, invite *Invite) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO club_invites (club_id, user_id, initiator) VALUES ($1, $2, $3)",
		invite.ClubId, invite.UserId, invite.Initiator)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetPendingInvites(ctx context.Context, clubId int) ([]Invite, error) {
	var invites []Invite
	err := r.db.SelectContext(ctx, &invites,
		"SELECT * FROM club_invites WHERE club_id = $1",
		clubId)
	if err != nil {
		return nil, err
	}
	return invites, nil
}

func (r *repository) GetUserInvites(ctx context.Context, userId int) ([]Invite, error) {
	var invites []Invite
	err := r.db.SelectContext(ctx, &invites,
		"SELECT * FROM club_invites WHERE user_id = $1",
		userId)
	if err != nil {
		return nil, err
	}
	return invites, nil
}

func (r *repository) GetInvite(ctx context.Context, id int) (*Invite, error) {
	var invite Invite
	err := r.db.GetContext(ctx, &invite,
		"SELECT * FROM club_invites WHERE id = $1",
		id)
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *repository) DeleteInvite(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM club_invites WHERE id = $1",
		id)
	if err != nil {
		return err
	}
	return nil
}
