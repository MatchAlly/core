package statistic

type Statistic struct {
	ID       uint
	MemberId uint `db:"member_id"`
	GameId   uint `db:"game_id"`

	Wins   int
	Draws  int
	Losses int
	Streak int
}
