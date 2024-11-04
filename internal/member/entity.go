package member

type Role string

const (
	AdminRole   Role = "admin"
	ManagerRole Role = "manager"
	MemberRole  Role = "member"
)

type Member struct {
	ID          int
	ClubID      int    `db:"club_id"`
	UserID      int    `db:"user_id"`
	DisplayName string `db:"display_name"`
	Role        Role
}
