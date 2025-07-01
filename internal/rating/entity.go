package rating

import (
	"time"

	"github.com/google/uuid"
)

const (
	startMu    = 25.0
	startSigma = 3.0
)

type Rating struct {
	ID        uuid.UUID `db:"id"`
	MemberID  uuid.UUID `db:"member_id"`
	GameID    uuid.UUID `db:"game_id"`
	Mu        float64   `db:"mu"`
	Sigma     float64   `db:"sigma"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
