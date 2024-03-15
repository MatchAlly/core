package invite

import (
	"time"
)

type Invite struct {
	Id uint `gorm:"primaryKey"`

	ClubId uint `gorm:"primaryKey"`
	UserId uint `gorm:"primaryKey"`

	CreatedAt time.Time
}
