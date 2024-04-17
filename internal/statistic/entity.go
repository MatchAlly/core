package statistic

import (
	"gorm.io/gorm"
)

type MatchResult int

const (
	ResultWin MatchResult = iota
	ResultLoss
	ResultDraw
)

type Measure string

const (
	MeasureWins   Measure = "wins"
	MeasureStreak Measure = "streak"
)

type Statistic struct {
	gorm.Model

	MemberId uint `gorm:"not null"`
	GameId   uint `gorm:"not null"`

	Wins   int
	Draws  int
	Losses int
	Streak int
}
