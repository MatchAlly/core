package subscription

import "time"

type Tier int

const (
	TierNone Tier = iota
	TierFree
	TierMinor
	TierMajor
)

type Subscription struct {
	ID                     int       `db:"id"`
	UserID                 int       `db:"user_id"`
	ManagedOrganizationIDs []int     `db:"managed_organization_ids"`
	TotalManagedUsers      int       `db:"total_managed_users"`
	Tier                   Tier      `db:"tier"`
	CreatedAt              time.Time `db:"created_at"`
	UpdatedAt              time.Time `db:"updated_at"`
}
