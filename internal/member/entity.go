package member

type Role string

const (
	AdminRole   Role = "ADMIN"
	ManagerRole Role = "MANAGER"
	MemberRole  Role = "MEMBER"
)

type Member struct {
	ID          int
	ClubID      int    `db:"club_id"`
	UserID      int    `db:"user_id"`
	DisplayName string `db:"display_name"`
	Role        Role
}
