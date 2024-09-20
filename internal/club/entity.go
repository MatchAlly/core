package club

type Club struct {
	ID   uint
	Name string
}

type Invite struct {
	ID     uint
	ClubId uint `db:"club_id"`
	UserId uint `db:"user_id"`
}

type JoinRequest struct {
	ID     uint
	ClubId uint `db:"club_id"`
	UserId uint `db:"user_id"`
}
