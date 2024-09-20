package game

type Game struct {
	ID     uint
	ClubID uint `db:"club_id"`
	Name   string
}
