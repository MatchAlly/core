package game

type Mode string

const (
	ModeFreeForAll Mode = "FFA"
	ModeTeam       Mode = "TEAM"
	ModeCoop       Mode = "COOP"
)

type Game struct {
	ID     int
	ClubID int `db:"club_id"`
	Name   string
}

type Gamemode struct {
	ID     int
	GameID int `db:"game_id"`
	Mode   Mode
}
