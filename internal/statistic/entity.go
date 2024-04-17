package statistic

import "time"

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
	Id uint `gorm:"primaryKey"`

	MemberId uint `gorm:"not null"`
	GameId   uint `gorm:"not null"`

	Wins   int
	Draws  int
	Losses int
	Streak int

	CreatedAt time.Time
}
