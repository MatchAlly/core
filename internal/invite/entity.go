package invite

import (
	"gorm.io/gorm"
)

type Invite struct {
	gorm.Model

	ClubId uint `gorm:"primaryKey"`
	UserId uint `gorm:"primaryKey"`
}
