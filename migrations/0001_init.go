package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migration00001Init contains the initial migration.
var Migration00001Init = &gormigrate.Migration{
	ID: "init_00001",
	Migrate: func(tx *gorm.DB) error {
		type Role string
		type GameType string
		type Result rune

		type Game struct {
			Id     uint     `gorm:"primaryKey"`
			ClubId uint     `gorm:"index;not null"`
			Name   string   `gorm:"not null"`
			Type   GameType `gorm:"not null"`
		}

		type Statistic struct {
			Id uint `gorm:"primaryKey"`

			MemberId uint `gorm:"not null"`
			GameId   uint `gorm:"not null"`

			Wins   int
			Draws  int
			Losses int
			Streak int

			CreatedAt time.Time
		}

		type Rating struct {
			Id uint `gorm:"primaryKey"`

			MemberId uint `gorm:"not null"`
			GameId   uint `gorm:"not null"`

			Value      float64 `gorm:"default:1000.0"`
			Deviation  float64
			Volatility float64 `gorm:"default:0.06"`

			CreatedAt time.Time
		}

		type Member struct {
			Id uint `gorm:"primaryKey"`

			ClubId uint `gorm:"not null"`
			Role   Role `gorm:"default:member"`
			UserId uint `gorm:"not null"`

			Statistics []Statistic `gorm:"constraint:OnDelete:CASCADE"`
			Ratings    []Rating    `gorm:"constraint:OnDelete:CASCADE"`

			CreatedAt time.Time
		}

		type Match struct {
			Id     uint `gorm:"primaryKey"`
			ClubId uint `gorm:"not null"`

			TeamA  []uint   `gorm:"serializer:json;not null"`
			TeamB  []uint   `gorm:"serializer:json;not null"`
			Sets   []string `gorm:"serializer:json;not null"`
			Result Result   `gorm:"not null"`

			CreatedAt time.Time
		}

		type Club struct {
			Id uint `gorm:"primaryKey"`

			Name string `gorm:"not null"`

			Members []Member `gorm:"constraint:OnDelete:CASCADE"`
			Games   []Game   `gorm:"constraint:OnDelete:CASCADE"`
			Matches []Match  `gorm:"constraint:OnDelete:CASCADE"`

			CreatedAt time.Time
		}

		return tx.AutoMigrate(
			&Game{},
			&Statistic{},
			&Rating{},
			&Member{},
			&Match{},
			&Club{},
		)
	},
}
