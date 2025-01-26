package game

type Mode int

const (
	ModeNone Mode = iota
	ModeFreeForAll
	ModeTeam
	ModeCoop
)

type Game struct {
	ID     int    `db:"id"`
	ClubID int    `db:"club_id"`
	Name   string `db:"name"`
}

type Gamemode struct {
	ID     int  `db:"id"`
	GameID int  `db:"game_id"`
	Mode   Mode `db:"mode"`
}
