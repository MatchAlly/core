package member

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
	ID          int    `db:"id"`
	ClubID      int    `db:"club_id"`
	UserID      int    `db:"user_id"`
	DisplayName string `db:"display_name"`
	Role        Role   `db:"role"`
}
