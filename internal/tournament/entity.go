package tournament

import "time"

type Tournament struct {
	Id uint `gorm:"primaryKey"`

	Name     string            `gorm:"not null"`
	GameID   uint              `gorm:"not null"`
	NumTeams uint              `gorm:"not null"`
	Teams    []TournamentTeam  `gorm:"not null"`
	Matches  []TournamentMatch `gorm:"not null"`

	CreatedAt time.Time
}

type TournamentTeam struct {
	Id uint `gorm:"primaryKey"`

	TournamentID uint `gorm:"not null"`
	TeamID       uint `gorm:"not null"`

	InitialSeed    uint
	FinalPlacement uint

	CreatedAt time.Time
}

type TournamentMatch struct {
	Id uint `gorm:"primaryKey"`

	TournamentID uint `gorm:"not null"`
	MatchID      uint `gorm:"not null"`

	Team1ID uint `gorm:"not null"`
	Team2ID uint `gorm:"not null"`

	Team1Score uint
	Team2Score uint

	CreatedAt time.Time
}
