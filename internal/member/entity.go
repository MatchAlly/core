package member

type Role string

const (
	AdminRole   Role = "admin"
	ManagerRole Role = "manager"
	MemberRole  Role = "member"
)

type Member struct {
	ID          uint
	ClubID      uint   `db:"club_id"`
	UserID      uint   `db:"user_id"`
	DisplayName string `db:"display_name"`
	Role        Role
}
