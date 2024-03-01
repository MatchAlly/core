package user

import (
	"time"
)

type User struct {
	Id uint `gorm:"primaryKey"`

	Email string `gorm:"index"`
	Name  string `gorm:"index;not null"`
	Hash  string `gorm:"not null"`

	UpdatedAt time.Time
	CreatedAt time.Time
}
