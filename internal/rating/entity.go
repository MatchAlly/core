package rating

import (
	"time"
)

const (
	startMu    = 25.0
	startSigma = 3.0
)

type Rating struct {
	ID        int
	MemberID  int `db:"member_id"`
	GameID    int `db:"game_id"`
	Mu        float64
	Sigma     float64
	UpdatedAt time.Time `db:"updated_at"`
}
