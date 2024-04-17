package user

import (
	"core/internal/club"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email string `gorm:"uniqueIndex"`
	Name  string `gorm:"index;not null"`
	Hash  string `gorm:"not null"`

	Memberships []club.Member `gorm:"constraint:OnDelete:CASCADE"`
}
