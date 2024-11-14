package match

import (
	"core/internal/member"
	"time"
)

type Match struct {
	ID        int       `json:"id" db:"id"`
	ClubID    int       `json:"club_id" db:"club_id"`
	GameID    int       `json:"game_id" db:"game_id"`
	Sets      []string  `json:"sets" db:"sets"`
	Teams     []Team    `json:"teams,omitempty"` // Must be loaded by joins
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Team struct {
	ID      int             `json:"id" db:"id"`
	ClubID  int             `json:"club_id" db:"club_id"`
	Members []member.Member `json:"members,omitempty"` // Must be loaded by joins
}
