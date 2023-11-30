package match

import (
	"time"
)

type Result rune

const (
	TeamAWins Result = 'A'
	TeamBWins Result = 'B'
	Draw      Result = 'D'
)

type Match struct {
	Id uint `gorm:"primaryKey"`

	TeamA  []uint   `gorm:"serializer:json;not null"`
	TeamB  []uint   `gorm:"serializer:json;not null"`
	Sets   []string `gorm:"serializer:json;not null"`
	Result Result   `gorm:"not null"`

	CreatedAt time.Time
}
