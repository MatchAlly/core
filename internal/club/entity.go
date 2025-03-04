package club

import "core/internal/member"

type Initiator int

const (
	InitiatorNone Initiator = iota
	IniatorClub
	InitiatorUser
)

type Club struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
}

type Invite struct {
	ID        int         `db:"id"`
	ClubId    int         `db:"club_id"`
	UserId    int         `db:"user_id"`
	Initiator Initiator   `db:"initiator"`
	Role      member.Role `db:"role"`
}
