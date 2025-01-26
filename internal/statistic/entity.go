package statistic

import "time"

type Statistic struct {
	ID        int       `db:"id"`
	MemberId  int       `db:"member_id"`
	GameId    int       `db:"game_id"`
	Wins      int       `db:"wins"`
	Draws     int       `db:"draws"`
	Losses    int       `db:"losses"`
	Streak    int       `db:"streak"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
