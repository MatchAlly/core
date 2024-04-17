package game

import "gorm.io/gorm"

type GameType string

const (
	FreeForAllGameType GameType = "ffa"
	TeamGameType       GameType = "team"
)

type Game struct {
	gorm.Model

	ClubId uint     `gorm:"index;not null"`
	Name   string   `gorm:"not null"`
	Type   GameType `gorm:"not null"`
}
