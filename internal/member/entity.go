package member

type Role string

const (
	RoleAdmin   Role = "ADMIN"
	RoleManager Role = "MANAGER"
	RoleMember  Role = "MEMBER"
)

type Member struct {
	ID          int
	ClubID      int    `db:"club_id"`
	UserID      int    `db:"user_id"`
	DisplayName string `db:"display_name"`
	Role        Role
}
