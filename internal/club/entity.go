package club

import (
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

	CreatedAt time.Time
}

type Invite struct {
	Id uint `gorm:"primaryKey"`

	ClubId uint `gorm:"primaryKey"`
	UserId uint `gorm:"primaryKey"`

	Accepted bool `gorm:"default:false"`
	Role     Role `gorm:"default:member"`

	CreatedAt time.Time
}

type ClubsGames struct {
	Id uint `gorm:"primaryKey"`

	ClubId uint `gorm:"primaryKey"`
	GameId uint `gorm:"primaryKey"`

	CreatedAt time.Time
}

type ClubsTournaments struct {
	Id uint `gorm:"primaryKey"`

	ClubId       uint `gorm:"primaryKey"`
	TournamentId uint `gorm:"primaryKey"`

	CreatedAt time.Time
}
