package migrations

import (
	"core/internal/club"
	"core/internal/match"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migration00001Init contains the initial migration.
var Migration00001Init = &gormigrate.Migration{
	ID: "init_00001",
	Migrate: func(tx *gorm.DB) error {
		type User struct {
			Id uint `gorm:"primaryKey"`

			Email   string `gorm:"index"`
			Name    string `gorm:"index;not null"`
			Hash    string
			Virtual bool `gorm:"default:false"`

			CreatedAt time.Time
		}

		type Statistic struct {
			Id uint `gorm:"primaryKey"`

			UserId uint `gorm:"not null"`
			GameId uint `gorm:"not null"`

			Wins   int
			Draws  int
			Losses int
			Streak int

			CreatedAt time.Time
		}

		type Rating struct {
			Id uint `gorm:"primaryKey"`

			UserId uint `gorm:"not null"`
			GameId uint `gorm:"not null"`

			Value      float64 `gorm:"default:1000.0"`
			Deviation  float64
			Volatility float64 `gorm:"default:0.06"`

			CreatedAt time.Time
		}

		type Game struct {
			Id   uint   `gorm:"primaryKey"`
			Name string `gorm:"not null"`
		}

		type Club struct {
			Id uint `gorm:"primaryKey"`

			Name string `gorm:"not null"`

			CreatedAt time.Time
		}

		type ClubsUsers struct {
			Id uint `gorm:"primaryKey"`

			ClubId uint `gorm:"primaryKey"`
			UserId uint `gorm:"primaryKey"`

			Accepted bool      `gorm:"default:false"`
			Role     club.Role `gorm:"default:member"`

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

		type Match struct {
			Id uint `gorm:"primaryKey"`

			TeamA  []uint       `gorm:"serializer:json;not null"`
			TeamB  []uint       `gorm:"serializer:json;not null"`
			Sets   []string     `gorm:"serializer:json;not null"`
			Result match.Result `gorm:"not null"`

			CreatedAt time.Time
		}

		return tx.AutoMigrate(
			&User{},
			&Statistic{},
			&Rating{},
			&Club{},
			&Match{},
			&Game{},
			&ClubsUsers{},
			&ClubsGames{},
			&ClubsTournaments{},
		)
	},
}
