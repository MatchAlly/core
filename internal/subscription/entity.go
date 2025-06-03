package subscription

import "time"

type Tier string

const (
	TierNone  Tier = "none"
	TierFree  Tier = "free"
	TierMinor Tier = "minor"
	TierMajor Tier = "major"
)

type Subscription struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Tier      Tier      `db:"tier"`
	CreatedAt time.Time `db:"created_at"`
}
