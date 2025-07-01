package game

import "github.com/google/uuid"

type Mode int

const (
	ModeNone Mode = iota
	ModeFreeForAll
	ModeTeam
	ModeCoop
)

func (m Mode) String() string {
	switch m {
	case ModeFreeForAll:
		return "FREE_FOR_ALL"
	case ModeTeam:
		return "TEAM"
	case ModeCoop:
		return "COOP"
	default:
		return "NONE"
	}
}

type Game struct {
	ID     uuid.UUID `db:"id"`
	ClubID uuid.UUID `db:"club_id"`
	Name   string    `db:"name"`
}

type Gamemode struct {
	ID     uuid.UUID `db:"id"`
	GameID uuid.UUID `db:"game_id"`
	Mode   Mode      `db:"mode"`
}
