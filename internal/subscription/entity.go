package subscription

import (
	"time"

	"github.com/google/uuid"
)

type Tier string

const (
	TierNone  Tier = "none"
	TierFree  Tier = "free"
	TierMinor Tier = "minor"
	TierMajor Tier = "major"
)

type Subscription struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Tier      Tier      `db:"tier"`
	CreatedAt time.Time `db:"created_at"`
}
