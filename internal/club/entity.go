package club

import (
	"core/internal/game"
	"core/internal/match"
	"core/internal/user"
	"time"
)

type Role string

const (
	AdminRole   Role = "admin"
	ManagerRole Role = "manager"
	MemberRole  Role = "member"
)

type Club struct {
	Id uint `gorm:"primaryKey"`

	Name string `gorm:"not null"`

	Users   []user.User   `gorm:"many2many:user_clubs;"`
	Games   []game.Game   `gorm:"constraint:OnDelete:CASCADE"`
	Matches []match.Match `gorm:"constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
}

type Member struct {
	Id uint `gorm:"primaryKey"`

	UserId uint `gorm:"not null"`
	ClubId uint `gorm:"not null"`
	Role   Role `gorm:"default:member"`

	CreatedAt time.Time
}
