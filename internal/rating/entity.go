package rating

import (
	"time"
)

const (
	startMu    = 25.0
	startSigma = 3.0
)

type Rating struct {
	ID        int       `db:"id"`
	MemberID  int       `db:"member_id"`
	GameID    int       `db:"game_id"`
	Mu        float64   `db:"mu"`
	Sigma     float64   `db:"sigma"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
