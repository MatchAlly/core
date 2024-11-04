package club

type Club struct {
	ID   int
	Name string
}

type Invite struct {
	ID     int
	ClubId int `db:"club_id"`
	UserId int `db:"user_id"`
}

type JoinRequest struct {
	ID     int
	ClubId int `db:"club_id"`
	UserId int `db:"user_id"`
}
