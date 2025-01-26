package member

type Role int

const (
	RoleNone Role = iota
	RoleMember
	RoleManager
	RoleAdmin
)

func (r Role) String() string {
	return [...]string{"none", "member", "manager", "admin"}[r]
}

type Member struct {
	ID          int    `db:"id"`
	ClubID      int    `db:"club_id"`
	UserID      int    `db:"user_id"`
	DisplayName string `db:"display_name"`
	Role        Role   `db:"role"`
}
