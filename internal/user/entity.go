package user

import (
	"core/internal/club"
	"time"
)

type User struct {
	Id uint `gorm:"primaryKey"`

	Email string `gorm:"uniqueIndex"`
	Name  string `gorm:"index;not null"`
	Hash  string `gorm:"not null"`

	Memberships []club.Member `gorm:"constraint:OnDelete:CASCADE"`

	UpdatedAt time.Time
	CreatedAt time.Time
}
