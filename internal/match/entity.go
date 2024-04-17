package match

import (
	"gorm.io/gorm"
)

type Result rune

const (
	TeamAWins Result = 'A'
	TeamBWins Result = 'B'
	Draw      Result = 'D'
)

type Match struct {
	gorm.Model

	ClubId uint `gorm:"not null"`

	TeamA  []uint   `gorm:"serializer:json;not null"`
	TeamB  []uint   `gorm:"serializer:json;not null"`
	Sets   []string `gorm:"serializer:json;not null"`
	Result Result   `gorm:"not null"`
}
