package match

import (
	"core/internal/game"
	"core/internal/member"
	"time"
)

type Match struct {
	ID        int       `json:"id" db:"id"`
	ClubID    int       `json:"club_id" db:"club_id"`
	GameID    int       `json:"game_id" db:"game_id"`
	Gamemode  game.Mode `json:"gamemode" db:"gamemode"`
	Ranked    bool      `json:"ranked" db:"ranked"`
	Sets      []string  `json:"sets" db:"sets"`
	Teams     []Team    `json:"teams,omitempty"` // Must be loaded by joins
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Team struct {
	ID      int             `json:"id" db:"id"`
	ClubID  int             `json:"club_id" db:"club_id"`
	Members []member.Member `json:"members,omitempty"` // Must be loaded by joins
}
