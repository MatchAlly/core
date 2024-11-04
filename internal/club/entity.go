package club

type Initiator string

const (
	ClubInitiator Initiator = "CLUB"
	UserInitiator Initiator = "USER"
)

type Club struct {
	ID        int
	Name      string
	CreatedAt string `db:"created_at"`
}

type Invite struct {
	ID        int
	ClubId    int `db:"club_id"`
	UserId    int `db:"user_id"`
	Initiator Initiator
}
