package club

import (
	"core/internal/member"

	"github.com/google/uuid"
)

type Initiator int

const (
	InitiatorNone Initiator = iota
	IniatorClub
	InitiatorUser
)

type Club struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt string    `db:"created_at"`
}

type Invite struct {
	ID        uuid.UUID   `db:"id"`
	ClubId    uuid.UUID   `db:"club_id"`
	UserId    uuid.UUID   `db:"user_id"`
	Initiator Initiator   `db:"initiator"`
	Role      member.Role `db:"role"`
}
