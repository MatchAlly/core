package rating

import "time"

const (
	maxDeviation = 2.015
	minDeviation = 0.175

	maxRating   = 3 * maxDeviation
	startRating = 0.0
	minRating   = -3 * maxDeviation

	maxVolatility   = 0.08
	startVolatility = 0.06
	minVolatility   = 0.04
	tau             = 0.5

	resultMultiplierWin  = 1.0
	resultMultiplierDraw = 0.5
	resultMultiplierLoss = 0.0
)

type Rating struct {
	ID int

	MemberID int `db:"member_id"`
	GameID   int `db:"game_id"`

	Value      float64
	Deviation  float64
	Volatility float64

	UpdatedAt time.Time
}
