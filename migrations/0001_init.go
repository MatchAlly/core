package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
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
			gorm.Model

			ClubId uint     `gorm:"index;not null"`
			Name   string   `gorm:"not null"`
			Type   GameType `gorm:"not null"`
		}

		type Statistic struct {
			gorm.Model

			MemberId uint `gorm:"not null"`
			GameId   uint `gorm:"not null"`

			Wins   int
			Draws  int
			Losses int
			Streak int
		}

		type Rating struct {
			gorm.Model

			MemberId uint `gorm:"not null"`
			GameId   uint `gorm:"not null"`

			Value      float64 `gorm:"default:1000.0"`
			Deviation  float64
			Volatility float64 `gorm:"default:0.06"`
		}

		type Invite struct {
			gorm.Model

			ClubId uint `gorm:"primaryKey"`
			UserId uint `gorm:"primaryKey"`
		}

		type Match struct {
			gorm.Model

			ClubId uint `gorm:"not null"`

			TeamA  []uint   `gorm:"serializer:json;not null"`
			TeamB  []uint   `gorm:"serializer:json;not null"`
			Sets   []string `gorm:"serializer:json;not null"`
			Result Result   `gorm:"not null"`
		}

		if err := tx.AutoMigrate(
			&Game{},
			&Statistic{},
			&Rating{},
			&Invite{},
			&Match{},
		); err != nil {
			return errors.Wrap(err, "failed to migrate tables with no foreign keys")
		}

		type Member struct {
			gorm.Model

			ClubId uint `gorm:"not null"`
			Role   Role `gorm:"default:member"`
			UserId uint `gorm:"not null"`

			Statistics []Statistic `gorm:"constraint:OnDelete:CASCADE"`
			Ratings    []Rating    `gorm:"constraint:OnDelete:CASCADE"`
		}

		if err := tx.AutoMigrate(
			&Member{},
		); err != nil {
			return errors.Wrap(err, "failed to migrate members")
		}

		type Club struct {
			gorm.Model

			Name string `gorm:"not null"`

			Members []Member `gorm:"constraint:OnDelete:CASCADE"`
			Games   []Game   `gorm:"constraint:OnDelete:CASCADE"`
			Matches []Match  `gorm:"constraint:OnDelete:CASCADE"`
			Invites []Invite `gorm:"constraint:OnDelete:CASCADE"`
		}

		if err := tx.AutoMigrate(
			&Club{},
		); err != nil {
			return errors.Wrap(err, "failed to migrate clubs")
		}

		return nil
	},
}
