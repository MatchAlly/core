package statistic

import "time"

type Statistic struct {
	ID int

	MemberId int `db:"member_id"`
	GameId   int `db:"game_id"`

	Wins   int
	Draws  int
	Losses int
	Streak int

	UpdatedAt time.Time
}
