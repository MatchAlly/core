package game

type GameType string

const (
	FreeForAllGameType GameType = "ffa"
	TeamGameType       GameType = "team"
	CoopGameType       GameType = "coop"
)

type Game struct {
	Id   uint     `gorm:"primaryKey"`
	Name string   `gorm:"not null"`
	Type GameType `gorm:"not null"`
}
