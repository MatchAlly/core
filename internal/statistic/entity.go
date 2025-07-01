package statistic

import (
	"time"

	"github.com/google/uuid"
)

type Statistic struct {
	ID        uuid.UUID `db:"id"`
	MemberId  uuid.UUID `db:"member_id"`
	GameId    uuid.UUID `db:"game_id"`
	Wins      int       `db:"wins"`
	Draws     int       `db:"draws"`
	Losses    int       `db:"losses"`
	Streak    int       `db:"streak"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
