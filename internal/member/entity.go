package member

import "github.com/google/uuid"

type Role string

const (
	RoleNone     Role = "none"
	RoleObserver Role = "observer"
	RoleMember   Role = "member"
	RoleManager  Role = "manager"
	RoleAdmin    Role = "admin"
	RoleOwner    Role = "owner"
)

type Member struct {
	ID     uuid.UUID `db:"id"`
	ClubID uuid.UUID `db:"club_id"`
	UserID uuid.UUID `db:"user_id"`
	Role   Role      `db:"role"`
}
