package club

import (
	"core/internal/game"
	"core/internal/match"
	"core/internal/rating"
	"core/internal/statistic"
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

	Members []Member      `gorm:"constraint:OnDelete:CASCADE"`
	Games   []game.Game   `gorm:"constraint:OnDelete:CASCADE"`
	Matches []match.Match `gorm:"constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
}

type Member struct {
	Id uint `gorm:"primaryKey"`

	ClubId uint `gorm:"not null"`
	Role   Role `gorm:"default:member"`
	UserId uint `gorm:"not null"`

	Statistics []statistic.Statistic `gorm:"constraint:OnDelete:CASCADE"`
	Ratings    []rating.Rating       `gorm:"constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
}
