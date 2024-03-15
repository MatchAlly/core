package game

type GameType string

const (
	FreeForAllGameType GameType = "ffa"
	TeamGameType       GameType = "team"
)

type Game struct {
	Id     uint     `gorm:"primaryKey"`
	ClubId uint     `gorm:"not null"`
	Name   string   `gorm:"not null"`
	Type   GameType `gorm:"not null"`
}
