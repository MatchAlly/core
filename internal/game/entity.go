package game

type Game struct {
	ID     int
	ClubID int `db:"club_id"`
	Name   string
}
